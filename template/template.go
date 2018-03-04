package template

import "errors"

// Template descrbies an elasticsearch Template
type Template struct {
	Name          Name        `json:"name"`
	FlushPeriodMS float64     `json:"flushPeriodMS"`
	Body          interface{} `json:"body"`
}

// Name is the name of an Template
type Name string

// GetTypes return declared for the given Template
func (Template *Template) GetTypes() ([]Name, error) {
	body := Template.Body.(map[string]interface{})
	typesMap, ok := body["mappings"].(map[string]interface{})
	if !ok {
		return nil, errors.New("ErrUnexpectedCast - correct is map[string]interface{}")
	}
	types := make([]Name, 0, len(typesMap))
	for t := range typesMap {
		types = append(types, Name(t))
	}
	return types, nil
}
