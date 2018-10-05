package services

import (
	"context"
	"encoding/json"
	_ "fmt"
	"github.com/go-martini/martini"
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

	// Start listening
	port := viper.GetInt("listener.port")
	logger.Noticef("Listening to API requests on port %d", port)
	m.RunOnAddr(":" + strconv.Itoa(port))
}

func registerEndpoints(martiniServer *martini.ClassicMartini) {
	martiniServer.Post("/report", report)
	martiniServer.Get("/hosts/list", getHostsList)
	martiniServer.Get("/hosts/:host", getHostConfigs)
	martiniServer.Post("/hosts/:host", setHostConfigs)
}

func report(res http.ResponseWriter, req *http.Request) {
	//Decode the JSON utils
	decoder := json.NewDecoder(req.Body)
	var data []orm.MetricFloat
	err := decoder.Decode(&data)
	if err != nil {
		logger.Error(err)
	}
	logger.Debugf("Received utils: %s", data[0].Key)

	//Parse host that sent this report
	hostName := infra.ParseHost(req.RemoteAddr)

	//Save this metric as a document in the "metrics" index in ES
	for _, metric := range data {
		// Set some more values for this doc
		metric.Timestamp = time.Now()
		metric.Host = hostName

		// Save this metrics as a doc in ES "metrics" index
		res, err := utils.GetESClient().Index().Index("metrics").Type("_doc").BodyJson(metric).Do(req.Context())
		if err != nil {
			logger.Error(err)
		}
		logger.Debugf("Indexed %s to index %s\n", res.Id, res.Index)
	}

	//Save host in hosts list
	hostsSave(hostName, req.Context())

	res.WriteHeader(http.StatusOK)
}

func getHostsList(res http.ResponseWriter, req *http.Request) {
	//Save host in hosts list
	hosts := hostsGetList(req.Context())
	hostsJson, _ := json.Marshal(hosts)

	res.Write(hostsJson)
	res.WriteHeader(http.StatusOK)
}

func getHostConfigs(params martini.Params) (int, string) {
	configs := hostsGetHostConfigs(params["host"], context.Background())
	return http.StatusOK, configs
}

func setHostConfigs(req *http.Request, params martini.Params) (int, string) {
	bodyBuffer, _ := ioutil.ReadAll(req.Body)
	sBody := string(bodyBuffer)

	//Save host in hosts list
	err := hostsSetHostConfigs(params["host"], sBody, context.Background())
	if err != nil {
		return http.StatusOK, "{}"
	} else {
		return http.StatusBadRequest, "{}"
	}
}
