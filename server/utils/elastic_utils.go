package utils

import (
	"context"
	"encoding/json"
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

func GetBodyKey(docSource *json.RawMessage, key string) string {
	var value string

	// json.RawMessage is just a byte[], so we need to marshal it to a string map
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(*docSource, &objmap)

	// Each key in the map is again a json.RawMessage, so we need to check if the key exists
	// and then marshal again the value
	if err == nil {
		if _, ok := objmap[key]; ok {
			json.Unmarshal(*objmap[key], &value)
		}
	}

	return value
}
