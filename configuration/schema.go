package configuration

import (
	"github.com/khezen/espipe/template"
)

// Configuration contains all configuration for the logger
type Configuration struct {
	Redis         Redis         `json:"redis"`
	Elasticsearch Elasticsearch `json:"elasticsearch"`
}

type Redis struct {
	Enabled   bool   `json:"enabled"`
	Address   string `json:"address"`
	Password  string `json:"password"`
	Partition int    `json:"partition"`
}

type Elasticsearch struct {
	Address   string              `json:"addr"`
	Templates []template.Template `json:"templates"`
	AWSAuth   *AWSAuth            `json:"AWSAuth,omitempty"`
	BasicAuth *BasicAuth          `json:"basicAuth,omitempty"`
}

// BasicAuth - username & password for HTTP Basic Auth
type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AWSAuth provide credential for AWS services signing
type AWSAuth struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}
