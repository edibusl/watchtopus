package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"watchtopus/agent/collectors"
	"watchtopus/infra"
	"watchtopus/orm"
)

var logger = logging.MustGetLogger("watchtopus")
var hostConfigs map[string]*json.RawMessage

func main() {
	infra.InitLogger()
	initDefaultConfigs()
	infra.InitConfigs("agent")
	fetchHostConfigsFromServer()
	go fetchHostConfigsPeriodically()

	logger.Notice("Agent started successfully. Collecting metrics.")

	collect()
}

func collect() {
	for {
		// Collect metrics in parallel threads (using goroutines)
		// Use a channel to blocking wait for all threads to return the metrics utils
		ch := make(chan []orm.MetricFloat)
		go collectors.CollectCpu(ch)
		go collectors.CollectMem(ch)
		go collectors.CollectPing(ch)
		metrics1, metrics2, metrics3 := <-ch, <-ch, <-ch

		// Combine metrics to a single array
		var allMetrics []orm.MetricFloat
		allMetrics = append(metrics1, metrics2...)
		allMetrics = append(allMetrics, metrics3...)

		// Encode metrics array to JSON string
		strJson, _ := json.Marshal(allMetrics)
		logger.Debug(string(strJson))

		// Send metrics JSON array to the server
		resp, err := http.Post(getBaseUrl()+"/report", "application/json", bytes.NewBuffer(strJson))
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			logger.Errorf("An error occurred: %s\n", err.Error())
		} else {
			logger.Debug("Sent report successfully")
		}

		if resp != nil {
			resp.Body.Close()
		}

		// No need to sleep since the CPU collector has a sleep of 2 seconds
		// (in order to make the diff between samples in different times)
	}
}

func fetchHostConfigsPeriodically() {
	for {
		// Update every minute
		time.Sleep(10 * time.Second)

		fetchHostConfigsFromServer()
	}
}

func fetchHostConfigsFromServer() {
	url := fmt.Sprintf("%s/hosts/%s", getBaseUrl(), viper.GetString("hostId"))
	resp, err := http.Get(url)
	if err != nil || (resp != nil && resp.StatusCode != 200) {
		if err != nil {
			logger.Errorf("An error occurred: %s\n", err.Error())
		} else {
			logger.Warningf("An error occurred. Status: %s\n", resp.Status)
		}

		return
	}

	// Parse host configs
	hostConfigs = infra.ParseResponseBody(resp)

	// Set the ping metrics configs
	if _, ok := hostConfigs["pingHosts"]; ok {
		collectors.InitPing(hostConfigs["pingHosts"])
	}
}

func getBaseUrl() string {
	return viper.GetStringSlice("servers")[0]
}

func initDefaultConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("servers", []string{"http://127.0.0.1:9200"})
}
