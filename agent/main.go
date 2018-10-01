package main

import (
	"bytes"
	"encoding/json"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"net/http"
	"watchtopus/agent/collectors"
	"watchtopus/infra"
	"watchtopus/orm"
)

var logger = logging.MustGetLogger("watchtopus")

func main() {
	infra.InitLogger()
	initDefaultConfigs()
	infra.InitConfigs("agent")

	collect()
}

func collect() {
	for {
		// Collect metrics in parallel threads (using goroutines)
		// Use a channel to blocking wait for all threads to return the metrics data
		ch := make(chan []orm.MetricFloat)
		go collectors.CollectCpu(ch)
		go collectors.CollectMem(ch)
		metrics1, metrics2 := <-ch, <-ch

		// Combine metrics to a single array
		var allMetrics []orm.MetricFloat
		allMetrics = append(metrics1, metrics2...)

		// Encode metrics array to JSON string
		strJson, _ := json.Marshal(allMetrics)
		logger.Debug(string(strJson))

		// Send metrics JSON array to the server
		baseUrl := viper.GetStringSlice("servers")[0]
		resp, err := http.Post(baseUrl+"/report", "application/json", bytes.NewBuffer(strJson))
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			logger.Errorf("An error occured: %s\n", err)
		} else {
			logger.Debug("Sent report successfully")
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

func initDefaultConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("servers", []string{"http://127.0.0.1:9200"})
}
