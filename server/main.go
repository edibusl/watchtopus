package main

import (
	"context"
	"github.com/go-martini/martini"
	"github.com/olivere/elastic"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strconv"
	"watchtopus/infra"
)

var _context context.Context
var _esClient *elastic.Client
var logger = logging.MustGetLogger("watchtopus")

func main() {
	infra.InitLogger()
	initDefaultConfigs()
	infra.InitConfigs("server")

	initElastic()
	startApiServer()
}

func initDefaultConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("elastics.host", "http://127.0.0.1:9200")
	viper.SetDefault("listener.port", 3000)
}

func startApiServer() {
	m := martini.Classic()
	registerEndpoints(m)

	// Disable martini logger
	m.Logger(log.New(ioutil.Discard, "", 0))

	// Start listening
	m.RunOnAddr(":" + strconv.Itoa(viper.GetInt("listener.port")))
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
	logger.Noticef("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := _esClient.IndexExists("metrics").Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		logger.Fatalf("Error! Index does not exist.")
	}
}
