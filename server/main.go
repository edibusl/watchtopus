package main

import (
	"context"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/olivere/elastic"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"time"
	"watchtopus/orm"
)

var _context context.Context
var _esClient *elastic.Client
var log = logging.MustGetLogger("server")

func main() {
	initLogger()
	initConfigs()
	initElastic()
	startApiServer()
}

func startApiServer() {
	m := martini.Classic()
	m.Post("/report", func(res http.ResponseWriter, req *http.Request) {
		//Decode the JSON data
		decoder := json.NewDecoder(req.Body)
		var data []orm.MetricFloat
		err := decoder.Decode(&data)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Received data: %s", data[0].Key)

		//Save this metric as a document in the "metrics" index in ES
		for _, metric := range data {
			metric.Timestamp = time.Now()
			res, err := _esClient.Index().Index("metrics").Type("_doc").BodyJson(metric).Do(req.Context())
			if err != nil {
				// Handle error
				log.Error(err)
			}
			log.Infof("Indexed %s to index %s\n", res.Id, res.Index)
		}

		res.WriteHeader(200)
	})
	m.RunOnAddr(":" + strconv.Itoa(viper.GetInt("listener.port")))
}

func initLogger() {
	// Create a new logging backend
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	// Create a format for the logger
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	// Create a chain of logger backends. We'll use the last one which includes all the options.
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	// Set the global logging backend
	logging.SetBackend(backendLeveled)
}

func initConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("elastics.host", "http://127.0.0.1:9200")
	viper.SetDefault("listener.port", 3000)

	// Set filename
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Set multiple paths to search for the conf file, including the running dir
	viper.AddConfigPath(".")
	dir, err := os.Getwd()
	if err == nil {
		viper.AddConfigPath(dir)
		viper.AddConfigPath(dir + "/server")
	}

	//Read config
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Warningf("Error reading configs file: %s. Using default keys. \n", err)
	}
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
	info, code, err := _esClient.Ping(viper.GetString("elastics.host")).Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	log.Noticef("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := _esClient.IndexExists("metrics").Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		log.Fatalf("Error! Index does not exist.")
	}
}
