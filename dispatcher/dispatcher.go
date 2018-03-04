package dispatcher

import (
	"sync"

	"github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/elastic"
	"github.com/khezen/espipe/template"
)

// Dispatcher dispatch logs message to Elasticsearch
type Dispatcher struct {
	Client  *elastic.Client
	buffers map[template.Name]Buffer
	relieve chan Buffer
	mutex   sync.RWMutex
}

// NewDispatcher creates a new Dispatcher object
func NewDispatcher(config *configuration.Configuration) (*Dispatcher, error) {
	buffers := make(map[template.Name]Buffer)
	client := elastic.NewClient(config)
	return &Dispatcher{
		client,
		buffers,
		make(chan Buffer),
		sync.RWMutex{},
	}, nil
}

// Dispatch takes incoming message into Elasticsearch
func (d *Dispatcher) Dispatch(document *document.Document) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.ensureBuffer(document)
	d.buffers[document.Template.Name].Append(*document)
}

func (d *Dispatcher) ensureBuffer(document *document.Document) {
	if _, ok := d.buffers[document.Template.Name]; !ok {
		buffer := DefaultBuffer(document.Template, d.Client)
		go buffer.Flusher()()
		d.buffers[document.Template.Name] = buffer
	}
}
