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
	config  configuration.Configuration
	mutex   sync.RWMutex
}

// NewDispatcher creates a new Dispatcher object
func NewDispatcher(config *configuration.Configuration) (*Dispatcher, error) {
	buffers := make(map[template.Name]Buffer)
	client := elastic.NewClient(config)
	return &Dispatcher{
		client,
		buffers,
		*config,
		sync.RWMutex{},
	}, nil
}

// Dispatch takes incoming message into Elasticsearch
func (d *Dispatcher) Dispatch(document *document.Document) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.ensureBuffer(document)
	go d.buffers[document.Template.Name].Append(*document)
}

func (d *Dispatcher) ensureBuffer(document *document.Document) {
	if _, ok := d.buffers[document.Template.Name]; !ok {
		var buffer Buffer
		if d.config.Redis.Enabled {
			buffer = RedisBuffer(document.Template, &d.config)
		} else {
			buffer = DefaultBuffer(document.Template, d.Client)
		}
		if d.config.Redis.AutoFlush {
			go buffer.Flusher()()
		}
		d.buffers[document.Template.Name] = buffer
	}
}
