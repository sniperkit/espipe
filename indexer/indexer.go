package indexer

import (
	"github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/dispatcher"
	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/errors"
	"github.com/khezen/espipe/template"
)

// Indexer indexes document in bulk request to elasticsearch
type Indexer struct {
	config             configuration.Configuration
	Dispatcher         *dispatcher.Dispatcher
	availableTemplates map[template.Name]template.Template
	availableResources map[template.Name]map[document.Type]bool
}

// New - Create new service for serving web REST requests
func New(config configuration.Configuration) (*Indexer, error) {
	d, err := dispatcher.NewDispatcher(&config)
	if err != nil {
		return nil, err
	}
	availableTemplates := make(map[template.Name]template.Template)
	availableResources := make(map[template.Name]map[document.Type]bool)
	for _, template := range config.Templates {
		availableTemplates[template.Name] = template
		availableResources[template.Name] = make(map[document.Type]bool)
		types, err := template.GetTypes()
		if err != nil {
			return nil, err
		}
		for _, t := range types {
			availableResources[template.Name][document.Type(t)] = true
		}
	}
	return &Indexer{
		config,
		d,
		availableTemplates,
		availableResources,
	}, nil
}

// Index document
func (i *Indexer) Index(docTemplate template.Name, docType document.Type, docBytes []byte) error {
	template, ok := i.availableTemplates[docTemplate]
	if !ok {
		return errors.ErrPathNotFound
	}
	if _, ok := i.availableResources[template.Name][docType]; !ok {
		return errors.ErrPathNotFound
	}
	// NO ERRORS -> DISPATCH
	document, err := document.NewDocument(&template, docType, docBytes)
	if err != nil {
		return err
	}
	i.Dispatcher.Dispatch(document)
	return nil
}
