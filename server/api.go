package main

import (
	"encoding/json"
	_ "fmt"
	"github.com/go-martini/martini"
	"net/http"
	"time"
	"watchtopus/infra"
	"watchtopus/orm"
)

func registerEndpoints(martiniServer *martini.ClassicMartini) {
	martiniServer.Post("/report", report)
}

func report(res http.ResponseWriter, req *http.Request) {
	//Decode the JSON data
	decoder := json.NewDecoder(req.Body)
	var data []orm.MetricFloat
	err := decoder.Decode(&data)
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf("Received data: %s", data[0].Key)

	//Save this metric as a document in the "metrics" index in ES
	for _, metric := range data {
		// Set some more values for this doc
		metric.Timestamp = time.Now()
		metric.Host = infra.ParseHost(req.RemoteAddr)

		// Save this metrics as a doc in ES "metrics" index
		res, err := _esClient.Index().Index("metrics").Type("_doc").BodyJson(metric).Do(req.Context())
		if err != nil {
			logger.Error(err)
		}
		logger.Debugf("Indexed %s to index %s\n", res.Id, res.Index)
	}

	res.WriteHeader(200)
}
