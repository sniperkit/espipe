package dispatcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/elastic"
	"github.com/khezen/espipe/template"
)

// buffer is related to a template
// It sends messages in bulk to elasticsearch.
type buffer struct {
	Template  *template.Template
	client    *elastic.Client
	Kill      chan error
	documents []document.Document
	mutex     sync.RWMutex
}

// DefaultBuffer creates a new buffer
func DefaultBuffer(template *template.Template, client *elastic.Client) Buffer {
	buffer := &buffer{
		template,
		client,
		make(chan error),
		make([]document.Document, 0),
		sync.RWMutex{},
	}
	return buffer
}

// Append to buffer
func (b *buffer) Append(msg document.Document) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.documents = append(b.documents, msg)
	return nil
}

// Flush the buffer
func (b *buffer) Flush() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if len(b.documents) == 0 {
		return nil
	}
	bulk := make([]byte, 0, 5000000)
	for _, doc := range b.documents {
		req, err := doc.Request()
		if err != nil {
			return err
		}
		bulk = append(bulk, req...)
	}
	err := b.client.Bulk(bulk)
	if err != nil {
		return err
	}
	b.documents = make([]document.Document, 0, bufferLimit)
	return nil
}

// Flusher flushes every {configuration/config.go::Template.TimerMS}
func (b *buffer) Flusher() func() {
	ticker := time.NewTicker(time.Duration(b.Template.FlushPeriodMS) * time.Millisecond)
	return func() {
		for {
			select {
			case <-b.Kill:
				return
			case <-ticker.C:
				go func() {
					err := b.Flush()
					if err != nil {
						fmt.Println(err)
					}
				}()
				break
			}
		}
	}
}
