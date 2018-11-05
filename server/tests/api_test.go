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
		orm.MetricFloat{
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

	resp, err := http.Post("http://localhost:3001/report", "application/json", bytes.NewBuffer(strJson))
	if err != nil || resp.StatusCode != 200 {
		t.Error("Calling /report failed")
	}

	// TODO - Verify that the data was written to elasticsearch
}
