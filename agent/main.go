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
		cpuMetrics := collectors.CollectCpu()
		memMetrics := collectors.CollectMem()

		var allMetrics []orm.MetricFloat
		allMetrics = append(cpuMetrics, memMetrics...)

		strJson, _ := json.Marshal(allMetrics)
		fmt.Println(string(strJson))

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
