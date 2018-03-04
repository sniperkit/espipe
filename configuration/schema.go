package configuration

import (
	"github.com/khezen/espipe/template"
)

// Configuration contains all configuration for the logger
type Configuration struct {
	Elasticsearch string              `json:"elasticsearch"`
	Templates     []template.Template `json:"templates"`
	AWSAuth       *AWSAuth            `json:"AWSAuth,omitempty"`
	BasicAuth     *BasicAuth          `json:"basicAuth,omitempty"`
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
