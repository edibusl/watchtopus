package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"watchtopus/agent/collectors"
	"watchtopus/orm"
)

func main() {
	initConfigs()
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
		fmt.Println(string(strJson))

		// Send metrics JSON array to the server
		baseUrl := viper.GetStringSlice("servers")[0]
		resp, err := http.Post(baseUrl+"/report", "application/json", bytes.NewBuffer(strJson))
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			fmt.Printf("An error occured: %s\n", err)
		} else {
			fmt.Println("Sent report successfully")
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}

func initConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("servers", []string{"http://127.0.0.1:9200"})

	// Set filename
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Set multiple paths to search for the conf file, including the running dir
	viper.AddConfigPath(".")
	dir, err := os.Getwd()
	if err == nil {
		viper.AddConfigPath(dir)
		viper.AddConfigPath(dir + "/agent")
	}

	//Read config
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		fmt.Printf("Error reading configs file: %s. Using default keys. \n", err)
	}
}
