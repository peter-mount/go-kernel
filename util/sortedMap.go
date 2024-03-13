package util

import "github.com/peter-mount/go-kernel/v2/util/strings"

type SortedMap[V any] map[string]V

func NewSortedMap[V any]() *SortedMap[V] {
	m := make(SortedMap[V])
	return &m
}

// AddAll will add all entries in the source map to this instance
func (m *SortedMap[V]) AddAll(source map[string]V) *SortedMap[V] {
	for k, v := range source {
		(*m)[k] = v
	}
	return m
}

// DecodeMap will add all entries in the source map to this instance.
func (m *SortedMap[V]) DecodeMap(source map[interface{}]interface{}) *SortedMap[V] {
	_ = m.decodeMap(source)
	return m
}

// Decode will add all entries in the source to this instance if it's a map. If not it does nothing.
func (m *SortedMap[V]) Decode(source interface{}) *SortedMap[V] {
	_ = IfMap(source, m.decodeMap)
	return m
}

func (m *SortedMap[V]) decodeMap(source map[interface{}]interface{}) error {
	for k, v := range source {
		ks := DecodeString(k, "")
		if ks != "" {
			(*m)[ks] = v.(V)
		}
	}
	return nil
}

func (m *SortedMap[V]) Keys() strings.StringSlice {
	var a strings.StringSlice
	for k, _ := range *m {
		a = append(a, k)
	}
	return a
}

func (m *SortedMap[V]) ForEach(f func(string, V) error) error {
	return m.Keys().
		Sort().
		ForEach(func(k string) error {
			return f(k, (*m)[k])
		})
}
