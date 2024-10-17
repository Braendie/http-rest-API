package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/http-rest-API/internal/app/apiserver"
)

var (
	configPath string
)

// init initializes the configuration for the application.
// It defines a command-line flag "config-path" that specifies the path to the configuration file.
// If the flag is not provided, it defaults to "config/apiserver.toml".
func init() {
	flag.StringVar(&configPath, "config-path", "config/apiserver.toml", "path to config file")
}

// main is the entry point of the application.
// It parses command-line flags, loads the configuration file, and starts the server.
func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
