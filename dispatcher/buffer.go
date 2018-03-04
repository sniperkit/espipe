package dispatcher

import (
	"github.com/khezen/espipe/document"
)

// Buffer -
type Buffer interface {
	Append(msg document.Document) error
	Flush() error
	Flusher() func()
}

const bufferLimit = 1000
