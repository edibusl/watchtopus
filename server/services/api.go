package services

import (
	"context"
	"encoding/json"
	_ "fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"watchtopus/infra"
	"watchtopus/orm"
	"watchtopus/server/utils"
)

func StartApiServer() {
	m := martini.Classic()
	registerEndpoints(m)

	// Disable martini logger
	m.Logger(log.New(ioutil.Discard, "", 0))

	//Enable CORS - Preflight requests
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Start listening
	port := viper.GetInt("listener.port")
	logger.Noticef("Listening to API requests on port %d", port)
	m.RunOnAddr(":" + strconv.Itoa(port))
}

func registerEndpoints(martiniServer *martini.ClassicMartini) {
	martiniServer.Post("/report", postReport)
	martiniServer.Get("/report/:host/last/:key", getLastReport)
	martiniServer.Get("/hosts/list", getHostsList)
	martiniServer.Get("/hosts/:host", getHostConfigs)
	martiniServer.Post("/hosts/:host", setHostConfigs)
}

func postReport(res http.ResponseWriter, req *http.Request) {
	//Decode the JSON utils
	decoder := json.NewDecoder(req.Body)
	var data []orm.MetricFloat
	err := decoder.Decode(&data)
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf("Received utils: %s", data[0].Key)

	//Parse host that sent this report
	hostIp := infra.ParseHost(req.RemoteAddr)

	//Save this metric as a document in the "metrics" index in ES
	var hostId string
	for _, metric := range data {
		// Set some more values for this doc
		metric.Timestamp = time.Now()
		metric.HostIp = hostIp

		// Save this metrics as a doc in ES "metrics" index
		res, err := utils.GetESClient().Index().Index("metrics").Type("_doc").BodyJson(metric).Do(req.Context())
		if err != nil {
			logger.Error(err)
		}
		logger.Debugf("Indexed %s to index %s\n", res.Id, res.Index)

		hostId = metric.HostId
	}

	//Save host in hosts list (if there was at least once metric)
	if hostId != "" {
		hostsSave(hostId, req.Context())
	}

	res.WriteHeader(http.StatusOK)
}

func getHostsList(res http.ResponseWriter, req *http.Request) {
	//Save host in hosts list
	hosts := hostsGetList(req.Context())
	hostsJson, _ := json.Marshal(hosts)

	res.WriteHeader(http.StatusOK)
	res.Write(hostsJson)
}

func getHostConfigs(params martini.Params, res http.ResponseWriter) {
	configs, err := hostsGetHostConfigs(params["host"], context.Background())

	if err == nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
	res.Write([]byte(configs))
}

func setHostConfigs(req *http.Request, params martini.Params) (int, string) {
	bodyBuffer, _ := ioutil.ReadAll(req.Body)
	sBody := string(bodyBuffer)

	//Save host in hosts list
	err := hostsSetHostConfigs(params["host"], sBody, context.Background())
	if err == nil {
		return http.StatusOK, "{}"
	} else {
		return http.StatusBadRequest, "{}"
	}
}

func getLastReport(params martini.Params, res http.ResponseWriter) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("hostId", params["host"]))
	query.Must(elastic.NewTermQuery("key", params["key"]))

	// Search with a term query
	searchResult, err := utils.GetESClient().Search().
		Index("metrics").
		Type("_doc").
		Query(query). // specify the query
		Pretty(true).
		Do(context.Background())

	if err != nil {
		// Handle error
		logger.Errorf("Error while fetching hosts list from ES %s", err.Error())
	}

	lastResult := make(map[string]string, 0)
	if searchResult.Hits.TotalHits > 0 {
		hit := searchResult.Hits.Hits[0]
		lastResult = map[string]string{
			"key":    utils.GetBodyKey(hit.Source, "key"),
			"hostId": utils.GetBodyKey(hit.Source, "hostId"),
			"val":    infra.FloatToString(utils.GetBodyKeyFloat(hit.Source, "val")),
		}
	} else {
		// No hits
		logger.Info("No reports found\n")
	}

	if err == nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
	lastResultJson, _ := json.Marshal(lastResult)
	res.Write(lastResultJson)
}

func addDefaultHeaders(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
}
