package jl

import (
	"encoding/json"
	"sort"
)

var SpecialFields = []string{
	"timestamp",
	"time",
	"level",
	"thread",
	"message",
	"msg",
	"logger",
	"exceptions",
}

func newEntry(m *Line, specialFields []string) *Entry {
	var specialKeys = stringSet(specialFields)
	var special, sorted []*field
	for _, k := range SpecialFields {
		if v, ok := m.Partials[k]; ok {
			special = append(special, newField(k, v))
		}
	}
	var sortedKeys = sortKeys(m.Partials)
	for _, k := range sortedKeys {
		if _, ok := specialKeys[k]; ok {
			continue
		}
		v := m.Partials[k]
		sorted = append(sorted, newField(k, v))
	}
	return &Entry{
		rawMessage:    m.Raw,
		partials: m.Partials,
		sortedFields:  sorted,
		specialFields: special,
	}
}

type field struct {
	Key   string
	Value interface{}
}

type Entry struct {
	rawMessage    []byte
	partials map[string]json.RawMessage
	fieldMap      map[string]interface{}
	sortedFields  []*field
	specialFields []*field
}

func (e *Entry) Level() Level {
	if levelField, ok := e.fieldMap["level"]; ok {
		if levelStr, ok := levelField.(string); ok {
			if l, err := ParseLevel(levelStr); err == nil {
				return l
			}
		}
	}
	return InfoLevel
}

func (e *Entry) Time() (string, bool) {
	if timestamp, ok := e.fieldMap["timestamp"]; ok {
		str, ok := timestamp.(string)
		return str, ok
	}
	if timestamp, ok := e.fieldMap["time"]; ok {
		str, ok := timestamp.(string)
		return str, ok
	}
	return "", false
}

func (e *Entry) Message() (string, bool) {
	if message, ok := e.fieldMap["message"]; ok {
		str, ok := message.(string)
		return str, ok
	}
	if message, ok := e.fieldMap["msg"]; ok {
		str, ok := message.(string)
		return str, ok
	}
	return "", false
}

func newField(k string, v interface{}) *field {
	return &field{Key: k, Value: v}
}

func stringSet(l []string) map[string]interface{} {
	m := make(map[string]interface{})
	for _, v := range l {
		m[v] = nil
	}
	return m
}

func sortKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.StringSlice(keys).Sort()
	return keys
}
