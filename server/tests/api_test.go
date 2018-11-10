package tests

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"testing"
	"time"
	"watchtopus/infra"
	"watchtopus/orm"
	"watchtopus/server/services"
	"watchtopus/server/utils"
)

func initDefaultConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("elastics.host", "http://127.0.0.1:9200")
	viper.SetDefault("listener.port", 3001)
}

func TestMain(m *testing.M) {
	// Setup
	infra.InitLogger()
	initDefaultConfigs()
	infra.InitConfigs("server")
	services.HostsInit()
	utils.InitElasticsearch()
	go services.StartApiServer()

	code := m.Run()

	os.Exit(code)
}

func TestReport(t *testing.T) {
	var allMetrics = []orm.MetricFloat{
		{
			HostId:      "abc123",
			HostIp:      "127.0.0.1",
			Timestamp:   time.Now(),
			Key:         "cpu.user.cpu1",
			Category:    "cpu",
			SubCategory: "user",
			Component:   "cpu1",
			Val:         50.2,
		},
	}

	strJson, _ := json.Marshal(allMetrics)

	// Post the report to the server
	resp, err := http.Post("http://localhost:3001/report", "application/json", bytes.NewBuffer(strJson))
	if err != nil || resp.StatusCode != 200 {
		t.Error("Calling POST /report failed")
	}

	// Get last report of this host
	resp, err = http.Get("http://localhost:3001/report/abc123/last/cpu.user.cpu1")
	if err != nil || resp.StatusCode != 200 {
		t.Error("Calling GET /report failed")
	}
	lastReport := infra.ParseResponseBody(resp)

	// Assert to verify that the correct values were written for the last report
	var key, hostId, val string
	json.Unmarshal(*lastReport["key"], &key)
	json.Unmarshal(*lastReport["hostId"], &hostId)
	json.Unmarshal(*lastReport["val"], &val)

	if key != allMetrics[0].Key {
		t.Error("Bad key")
	}
	if hostId != allMetrics[0].HostId {
		t.Error("Bad hostId")
	}
	if val != infra.FloatToString(allMetrics[0].Val) {
		t.Error("Bad val")
	}
}
