package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Config parametrizes Kaguya's configuration.
type Config struct {
	Boards         []BoardConfig
	PostgresConfig PostgresConfig
	LnxConfig      LnxConfig
	BatchSize      int
}

type BoardConfig struct {
	Name          string
	ForceRecreate bool
}

type PostgresConfig struct {
	ConnectionString string
}

type LnxConfig struct {
	Host          string
	Port          int
	Configuration IndexConfiguration
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
//Errors might be returned due to IO or invalid JSON.
func LoadConfig() (Config, error) {
	configFile := os.Getenv("MOON_CONFIG")

	if configFile == "" {
		configFile = "./config.json"
	}

	blob, err := ioutil.ReadFile(configFile)

	if err != nil {
		return Config{}, fmt.Errorf("Error loading configuration file: %s", err)
	}

	var conf Config

	err = json.Unmarshal(blob, &conf)

	if err != nil {
		return Config{}, fmt.Errorf(
			"Error unmarshalling configuration file contents to JSON:\n File contents: %s\n Error message: %s",
			blob,
			err,
		)
	}

	return conf, nil
}
