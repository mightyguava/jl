package jl

import (
	"encoding/json"
)

// FieldFinder locates a field in the Entry and returns it.
type FieldFinder func(entry *Entry) interface{}

// ByNames locates fields by their top-level key name in the JSON log entry, and returns the field as a json.RawMessage.
func ByNames(names ...string) FieldFinder {
	return func(entry *Entry) interface{} {
		for _, name := range names {
			if v, ok := entry.Partials[name]; ok {
				return v
			}
		}
		return nil
	}
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

// JavaExceptionFinder finds a Java exception containing a stracetrace and returns it as a JavaExceptions.
func JavaExceptionFinder(entry *Entry) interface{} {
	var java struct {
		Exceptions []*JavaException `json:"exceptions"`
	}
	if err := json.Unmarshal(entry.Raw, &java); err == nil && len(java.Exceptions) > 0 {
		return JavaExceptions(java.Exceptions)
	}
	return nil
}
