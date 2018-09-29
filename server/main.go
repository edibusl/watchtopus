package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/olivere/elastic"
	"log"
	"net/http"
	"time"
)

type MetricFloat struct {
	Key string			`json:"key"`
	Val float64			`json:"val"`
	Category string		`json:"category"` //cpu, memory, network
	SubCategory string	`json:"subcategory"` //user, system, idle
	Component string	`json:"component"` //Optional - cpu0, cpu1, cpu2, etc..
	Timestamp time.Time	`json:"timestamp"` //Optional - cpu0, cpu1, cpu2, etc..
}




var _context context.Context
var _esClient *elastic.Client

func main() {
	initElastic()

	m := martini.Classic()
	m.Post("/report", func(res http.ResponseWriter, req *http.Request) {
		//Decode the JSON data
		decoder := json.NewDecoder(req.Body)
		var data []MetricFloat
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		log.Printf("Received data: %s", data[0].Key)

		//Save this metric as a document in the "metrics" index in ES
		for _, metric := range data{
			metric.Timestamp = time.Now()
			res, err := _esClient.Index().Index("metrics").Type("_doc").BodyJson(metric).Do(req.Context())
			if err != nil {
				// Handle error
				panic(err)
			}
			fmt.Printf("Indexed %s to index %s\n", res.Id, res.Index)
		}



		res.WriteHeader(200)
	})
	m.RunOnAddr(":3000")
}


func initElastic() {
	_context = context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	var err error
	_esClient, err = elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := _esClient.Ping("http://127.0.0.1:9200").Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)


	// Use the IndexExists service to check if a specified index exists.
	exists, err := _esClient.IndexExists("metrics").Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		fmt.Println("Error! Index does not exist.")
	}
}
