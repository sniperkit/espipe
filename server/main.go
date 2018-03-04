package server

import (
	"github.com/khezen/espipe/configuration"
)

var (
	configFile = "/etc/espipe/config.json"
)

func main() {
	quit := make(chan error)
	var err error
	config, err := configuration.Get()
	if config == nil {
		configuration.Set("config.json")
		config, err = configuration.Get()
	}
	config.Redis.AutoFlush = true
	server, err := New(*config, quit)
	if err != nil {
		panic(err)
	}
	go server.ListenAndServe()
	err = <-quit
	panic(err)
}
