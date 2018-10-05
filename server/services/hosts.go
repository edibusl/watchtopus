package services

import (
	"context"
	"encoding/json"
	_ "fmt"
	"github.com/olivere/elastic"
	"github.com/op/go-logging"
	"watchtopus/infra"
	"watchtopus/server/utils"
)

var availableHosts []string
var logger = logging.MustGetLogger("watchtopus")

func HostsInit() {
	availableHosts = make([]string, 0)
}

func hostsSave(host string, ctx context.Context) {
	// Skip hosts that we already know that have a key in ES
	if infra.FindInArray(availableHosts, host) {
		return
	}

	// Check whether a key for this host already exists
	res, err := utils.GetESClient().Get().Index("hosts").Type("_doc").Id(host).Do(ctx)
	if err != nil || res == nil || !res.Found {
		// Key doesn't exist, create a key
		_, err := utils.GetESClient().Index().Index("hosts").Type("_doc").Id(host).BodyString(`{}`).Do(ctx)
		if err != nil {
			logger.Error("Error getting host's key from Elasticsearch")
		} else {
			// Success - Add the host to availableHosts to cache this host avoid
			// retrying adding this key on subsequent requests
			availableHosts = append(availableHosts, host)
		}
	} else {
		// Key exists - add to cache
		availableHosts = append(availableHosts, host)
	}
}

func hostsGetList(ctx context.Context) []string {
	hosts := make([]string, 0)

	// Search with a term query
	searchResult, err := utils.GetESClient().Search().
		Index("hosts"). // search in index "twitter"
		Type("_doc").
		Query(elastic.NewMatchAllQuery()). // specify the query
		Pretty(true).
		StoredFields("").
		Do(ctx)

	if err != nil {
		// Handle error
		logger.Errorf("Error while fetching hosts list from ES %s", err)
	}

	if searchResult.Hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			hosts = append(hosts, hit.Id)
		}
	} else {
		// No hits
		logger.Info("Found no hosts\n")
	}

	return hosts
}

func hostsGetHostConfigs(host string, ctx context.Context) string {
	res, err := utils.GetESClient().Get().
		Index("hosts").
		Type("_doc").
		Id(host).
		Do(ctx)

	if err != nil {
		// Handle error
		logger.Errorf("Error while getting host configs of host %s. Error: %s", host, err)
		panic(err)
	}

	// Convert document body to string
	j, err := json.Marshal(&res.Source)
	if err != nil {
		panic(err)
	}

	return string(j)
}

func hostsSetHostConfigs(host string, configs string, ctx context.Context) error {
	_, err := utils.GetESClient().Index().
		Index("hosts").
		Type("_doc").
		Id(host).
		BodyString(configs).
		Do(ctx)

	if err != nil {
		// Handle error
		logger.Errorf("Error while saving host configs of host %s. Error: %s", host, err)

		return err
	}

	return nil
}
