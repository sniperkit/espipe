package server

import (
	"github.com/khezen/espipe/configuration"
)

var (
	configFile = "/etc/espipe/config.json"
)

func main() {
	quit := make(chan error)

	config, err := configuration.LoadConfig(configFile)
	if err != nil {
		config, err = configuration.LoadConfig("config.json")
		if err != nil {
			panic(err)
		}
	}

	server, err := New(config, quit)
	if err != nil {
		panic(err)
	}

	go server.ListenAndServe()
	err = <-quit
	panic(err)
}
