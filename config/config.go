//Package config wraps Moon configuration
package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//Config parametrizes Moon's configuration
type Config struct {
	Boards         []BoardConfig  `toml:"boards"`
	PostgresConfig PostgresConfig `toml:"postgres"`
	LnxConfig      LnxConfig      `toml:"lnx"`
}

//BoardConfig parametrizes Moon's configuration
//for indexing a board in Lnx
type BoardConfig struct {
	Name          string `toml:"name"`
	ForceRecreate bool   `toml:"force_recreate"`
}

//PostgresConfig parametrizes configuration
//for the db connection
type PostgresConfig struct {
	ConnectionString string `toml:"connection_string"`
}

//LnxConfig parametrizes configuration for
//Lnx searching and indexing
type LnxConfig struct {
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	BatchSize      int    `toml:"batch_size"`
	NapTime        string `toml:"nap_time"`
	ReaderThreads  int    `toml:"reader_threads"`
	MaxConcurrency int    `toml:"max_concurrency"`
	WriterBuffer   int    `toml:"writer_buffer"`
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
func LoadConfig() Config {
	configFile := os.Getenv("MOON_CONFIG")

	if configFile == "" {
		configFile = "./config.toml"
	}

	f, err := os.Open(configFile)

	if err != nil {
		log.Fatalf("Error loading configuration file: %s", err)
	}

	defer f.Close()

	var conf Config

	if _, err := toml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalln(err)
	}

	return conf
}
