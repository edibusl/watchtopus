package main

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"watchtopus/infra"
	"watchtopus/server/services"
	"watchtopus/server/utils"
)

var logger = logging.MustGetLogger("watchtopus")

func main() {
	infra.InitLogger()
	initDefaultConfigs()
	infra.InitConfigs("server")
	services.HostsInit()

	utils.InitElasticsearch()
	services.StartApiServer()
}

func initDefaultConfigs() {
	// Set defaults in case that the conf file is not found
	viper.SetDefault("elastics.host", "http://127.0.0.1:9200")
	viper.SetDefault("listener.port", 3000)
}
