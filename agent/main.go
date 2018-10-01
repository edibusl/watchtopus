package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"watchtopus/agent/collectors"
	"watchtopus/orm"
)

func main() {
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
		resp, err := http.Post("http://localhost:3000/report", "application/json", bytes.NewBuffer(strJson))
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
