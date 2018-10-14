package collectors

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
	"watchtopus/orm"
)

type PingResponse struct {
	host    string
	success bool
}

var pingConfigs []string

func InitPing(configs *json.RawMessage) {
	err := json.Unmarshal(*configs, &pingConfigs)
	if err != nil {
		logger.Errorf("Error while parsing config pingsList that was received from server. %s", err.Error())
	} else {
		logger.Notice("Updated agent configs from server successfully.")
	}
}

func CollectPing(ch chan []orm.MetricFloat) {
	metrics := make([]orm.MetricFloat, 0)

	if pingConfigs == nil || len(pingConfigs) == 0 {
		ch <- metrics
		return
	}

	// Make pings to all hosts in parallel
	chPing := make(chan PingResponse)
	for _, pingHost := range pingConfigs {
		go makePing(pingHost, chPing)
	}

	// Wait for results
	var results []PingResponse
	for i := 0; i < len(pingConfigs); i++ {
		result := <-chPing
		results = append(results, result)
	}

	// Create metrics out of results
	for _, pingResult := range results {
		val := 0.0
		if pingResult.success {
			val = 1.0
		}

		metrics = append(metrics, orm.MetricFloat{
			HostId:      viper.GetString("hostId"),
			Key:         "custom.ping",
			Val:         val,
			Category:    "custom",
			SubCategory: "ping",
			Component:   pingResult.host})
	}

	ch <- metrics
}

func makePing(hostDnsOrIp string, ch chan PingResponse) {
	// Check if this is an IP address
	var ip string
	addr := net.ParseIP(hostDnsOrIp)
	if addr != nil {
		ip = hostDnsOrIp
	} else {
		// It's not an IP address, probably a DNS. try to lookup the IP
		hostName, err := net.LookupHost(hostDnsOrIp)
		if len(hostName) == 0 || err != nil {
			logger.Errorf("Bad hostname in ping list configuration '%s': %s", hostDnsOrIp, err.Error())
			ch <- PingResponse{host: hostDnsOrIp, success: false}
		}

		ip = hostName[0]
	}

	// IP address resolved. Prepare for pinging
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		logger.Errorf("Error when pinging '%s': %s", hostDnsOrIp, err.Error())
		ch <- PingResponse{host: hostDnsOrIp, success: false}
	}
	p.AddIPAddr(ra)

	// Listen in a callback to an ICMP response
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ch <- PingResponse{host: hostDnsOrIp, success: true}
	}

	// Listen in a callback to a timeout
	p.OnIdle = func() {
		ch <- PingResponse{host: hostDnsOrIp, success: false}
	}

	// Listen to incoming ICMP replies using udp sockets.
	// This does not require root permissions, but does require running the following command on the aganet host:
	// sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
	p.Network("udp")

	// Timeout of 1 second
	p.MaxRTT = 1 * time.Second

	// Start
	err = p.Run()
	if err != nil {
		logger.Errorf("Error while pinging '%s': %s.", hostDnsOrIp, err.Error())
		ch <- PingResponse{host: hostDnsOrIp, success: false}
	}
}
