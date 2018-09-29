package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"watchtopus-agent/collectors"
)

func main() {
	for {
		var cpuMetrics []collectors.MetricFloat
		cpuMetrics = collectors.CollectCpu()

		strJson, _ := json.Marshal(cpuMetrics)
		fmt.Println(string(strJson))

		resp, err := http.Post("http://localhost:3000/report", "application/json", bytes.NewBuffer(strJson))
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			fmt.Printf("An error occured: %s\n", err)
		}else{
			fmt.Println("Sent report successfully")
		}

		if resp != nil{
			resp.Body.Close()
		}
	}
}

