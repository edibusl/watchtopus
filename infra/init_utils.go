package infra

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"os"
)

var logger = logging.MustGetLogger("watchtopus")

func InitLogger() {
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

	// Set custom log level that was passed as en environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := logging.LogLevel(logLevel)
		if err == nil {
			backendLeveled.SetLevel(level, "")
		}
	}

	// Set the global logging backend
	logging.SetBackend(backendLeveled)
}

func InitConfigs(componentName string) {
	// Set filename
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Set multiple paths to search for the conf file, including the running dir
	viper.AddConfigPath(".")
	dir, err := os.Getwd()
	if err == nil {
		viper.AddConfigPath(dir)
		viper.AddConfigPath(dir + "/" + componentName + "/conf")
	}

	// Read config
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		logger.Warningf("Error reading configs file: %s. Using default keys. \n", err)
	}

	// Check if hostId exists in configs
	hostId := viper.GetString("hostId")
	if hostId == "" {
		// Generate a random UUID, and take the first 7 chars as the hostId
		uuid := fmt.Sprintf("%s", uuid.Must(uuid.NewV4()))[1:8]
		viper.Set("hostId", uuid)

		// Save back the config file
		viper.WriteConfig()
	}
}
