package configuration

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var singleton *Configuration
var path = "/etc/espipe/config.json"
var rwMutex = sync.RWMutex{}

func setSingleton(config *Configuration) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	singleton = config
}

func getSingleton() *Configuration {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return singleton
}

// Get the config
func Get() (config *Configuration, err error) {
	config = getSingleton()
	if config == nil {
		config, err = loadConfig(path)
		if err != nil {
			config, err = loadConfig("config.json")
			if err != nil {
				return nil, err
			}
		}
	}
	return config, nil
}

func loadConfig(configFile string) (*Configuration, error) {
	// Load config file
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Configuration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	singleton = &config
	return &config, nil
}
