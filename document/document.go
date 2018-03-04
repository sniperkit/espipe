package document

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/khezen/espipe/template"
	"github.com/khezen/espipe/uuid"
)

const (
	anonymous = "anonymous"
)

// ID is uuid for documents
type ID string

// Type is type for documents
type Type string

// Document has a JSON body which must be indexed in given Template as given Type.
type Document struct {
	Template  *template.Template
	Type      Type
	ID        ID
	Timestamp time.Time
	Body      []byte
}

// NewDocument creates a document from es index, document type and its body
func NewDocument(index *template.Template, docType Type, body []byte) (*Document, error) {
	id := uuid.New()
	var bodyMap map[string]interface{}
	err := json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}
	t := time.Now().UTC()
	bodyMap["post_date"] = t.Format(time.RFC3339)
	body, err = json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	return &Document{
		index,
		docType,
		ID(id),
		t,
		body,
	}, nil
}

// SizeKB returns document size in KB
func (d *Document) SizeKB() float64 {
	return float64(len(d.Body)) / 1000
}

// Request returns the JSON request to be append to the bulk
func (d *Document) Request() ([]byte, error) {

	request := make(map[string]interface{})

	//{ "index" : { "_index" : "logs-2017.05.28", "_type" : "log", "_id" : "1" } }
	docDescription := make(map[string]interface{})
	docDescription["_index"] = RenderIndex(d.Template.Name, d.Timestamp)
	docDescription["_type"] = d.Type
	docDescription["_id"] = d.ID

	request["index"] = docDescription
	body, err := json.Marshal(request)
	body = append(body, '\n')
	body = append(body, d.Body...)
	body = append(body, '\n')
	if err != nil {
		return nil, err
	}
	return body, nil
}

// RenderIndex - logs: logs-2017.05.26
func RenderIndex(templateName template.Name, t time.Time) string {
	indexBuf := bytes.NewBufferString(string(templateName))
	indexBuf.WriteString("-")
	indexBuf.WriteString(t.Format("2006.01.02"))
	return indexBuf.String()
}
