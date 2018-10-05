package utils

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var _context context.Context
var _esClient *elastic.Client
var logger = logging.MustGetLogger("watchtopus")

func InitElasticsearch() {
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

	// Ping the ElasticSearch server to get e.g. the version number
	info, _, err := _esClient.Ping(viper.GetString("elastics.host")).Do(_context)
	if err != nil {
		// Handle error
		panic(err)
	}
	logger.Noticef("Elasticsearch ready. Version %s.", info.Version.Number)

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

func GetESClient() *elastic.Client {
	return _esClient
}
