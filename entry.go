package jl

import "sort"

var SpecialFields = []string{
	"timestamp",
	"level",
	"thread",
	"message",
	"logger",
}

func newEntry(m *Line, specialFields []string) *Entry {
	var specialKeys = stringSet(specialFields)
	var special, sorted []*field
	for _, k := range SpecialFields {
		if v, ok := m.JSON[k]; ok {
			special = append(special, newField(k, v))
		}
	}
	var sortedKeys = sortKeys(m.JSON)
	for _, k := range sortedKeys {
		if _, ok := specialKeys[k]; ok {
			continue
		}
		v := m.JSON[k]
		sorted = append(sorted, newField(k, v))
	}
	return &Entry{
		fieldMap:      m.JSON,
		sortedFields:  sorted,
		specialFields: special,
	}
}

type field struct {
	Key   string
	Value interface{}
}

type Entry struct {
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

func sortKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.StringSlice(keys).Sort()
	return keys
}