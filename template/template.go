package template

import "errors"

// Template descrbies an elasticsearch Template
type Template struct {
	Name         Name        `json:"name"`
	BufferSizeKB float64     `json:"bufferSizeKB"`
	TimerMS      float64     `json:"timerMS"`
	Body         interface{} `json:"body"`
}

// Name is the name of an Template
type Name string

// GetTypes return declared for the given Template
func (Template *Template) GetTypes() ([]string, error) {
	body := Template.Body.(map[string]interface{})
	typesMap, ok := body["mappings"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Expected map[string]interface{}. Got something else.")
	}
	types := make([]string, 0, len(typesMap))
	for t := range typesMap {
		types = append(types, t)
	}
	return types, nil
}
