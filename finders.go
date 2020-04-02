package jl

import (
	"encoding/json"
	"strings"
)

// FieldFinder locates a field in the Entry and returns it.
type FieldFinder func(entry *Entry) interface{}

// ByNames locates fields by their top-level key name in the JSON log entry, and returns the field as a json.RawMessage.
func ByNames(names ...string) FieldFinder {
	return func(entry *Entry) interface{} {
		for _, name := range names {
			if v, ok := getDeep(entry, name); ok {
				return v
			}
		}
		return nil
	}
}

func getDeep(entry *Entry, name string) (interface{}, bool) {
	parts := strings.SplitN(name, ".", 2)
	key := parts[0]
	v, ok := entry.Partials[key]
	if !ok {
		return nil, false
	}
	if len(parts) == 1 {
		return v, true
	}
	var partials map[string]json.RawMessage
	err := json.Unmarshal(v, &partials)
	if err != nil {
		return nil, false
	}
	return getDeep(&Entry{
		Partials: partials,
		Raw:      v,
	}, parts[1])
}

// LogrusErrorFinder finds logrus error in the JSON log and returns it as a LogrusError.
func LogrusErrorFinder(entry *Entry) interface{} {
	var errStr, stack string
	if errV, ok := entry.Partials["error"]; !ok {
		return nil
	} else if err := json.Unmarshal(errV, &errStr); err != nil {
		return nil
	}
	if stackV, ok := entry.Partials["stack"]; !ok {
		return nil
	} else if err := json.Unmarshal(stackV, &stack); err != nil {
		return nil
	}
	return LogrusError{errStr, stack}
}
