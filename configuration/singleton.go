package configuration

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var singleton *Configuration
var singRWMut = sync.RWMutex{}
var path = "/etc/espipe/config.json"
var pathRWMut = sync.RWMutex{}

func setSingleton(config *Configuration) {
	singRWMut.Lock()
	defer singRWMut.Unlock()
	singleton = config
}

func getGet() *Configuration {
	singRWMut.RLock()
	defer singRWMut.RUnlock()
	return singleton
}

// Set the config
func Set(configPath string) {
	pathRWMut.Lock()
	defer pathRWMut.Unlock()
	setSingleton(nil)
	path = configPath
}

// Get the config
func Get() (config *Configuration, err error) {
	config = getGet()
	if config == nil {
		pathRWMut.RLock()
		defer pathRWMut.RUnlock()
		config, err = loadConfig(path)
		if err != nil {
			return nil, err
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
