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
	config     configuration.Configuration
	templates  map[template.Name]template.Template
	types      map[template.Name]map[document.Type]bool
	dispatcher *dispatcher.Dispatcher
}

// New - Create new service for serving web REST requests
func New(config configuration.Configuration) (*Indexer, error) {
	d, err := dispatcher.NewDispatcher(&config)
	if err != nil {
		return nil, err
	}
	templates := make(map[template.Name]template.Template)
	types := make(map[template.Name]map[document.Type]bool)
	for _, template := range config.Elasticsearch.Templates {
		templates[template.Name] = template
		types[template.Name] = make(map[document.Type]bool)
		templateTypes, err := template.GetTypes()
		if err != nil {
			return nil, err
		}
		for _, t := range templateTypes {
			types[template.Name][document.Type(t)] = true
		}
	}
	return &Indexer{
		config,
		templates,
		types,
		d,
	}, nil
}

// Index document
func (i *Indexer) Index(docTemplate template.Name, docType document.Type, docBytes []byte) error {
	template, ok := i.templates[docTemplate]
	if !ok {
		return errors.ErrPathNotFound
	}
	if _, ok := i.types[template.Name][docType]; !ok {
		return errors.ErrPathNotFound
	}
	// NO ERRORS -> DISPATCH
	document, err := document.NewDocument(&template, docType, docBytes)
	if err != nil {
		return err
	}
	i.dispatcher.Dispatch(document)
	return nil
}
