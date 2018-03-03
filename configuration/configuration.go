package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/khezen/espipe/template"
)

// Configuration contains all configuration for the logger
type Configuration struct {
	Elasticsearch string              `json:"elasticsearch"`
	Templates     []template.Template `json:"templates"`
	AWSAuth       *AWSAuth            `json:"AWSAuth,omitempty"`
	BasicAuth     *Crendentials       `json:"basicAuth,omitempty"`
}

// Crendentials - username & password for HTTP Basic Auth
type Crendentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AWSAuth provide credential for AWS services signing
type AWSAuth struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}

// LoadConfig reads the configuration from the config JSON file
func LoadConfig(configFile string) (Configuration, error) {

	// Load config file
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}
