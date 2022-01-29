package util

import "github.com/peter-mount/go-kernel/util/strings"

type SortedMap map[string]interface{}

func NewSortedMap() *SortedMap {
	m := make(SortedMap)
	return &m
}

// AddAll will add all entries in the source map to this instance
func (m *SortedMap) AddAll(source map[string]interface{}) *SortedMap {
	for k, v := range source {
		(*m)[k] = v
	}
	return m
}

// DecodeMap will add all entries in the source map to this instance.
func (m *SortedMap) DecodeMap(source map[interface{}]interface{}) *SortedMap {
	_ = m.decodeMap(source)
	return m
}

// Decode will add all entries in the source to this instance if it's a map. If not it does nothing.
func (m *SortedMap) Decode(source interface{}) *SortedMap {
	_ = IfMap(source, m.decodeMap)
	return m
}

func (m *SortedMap) decodeMap(source map[interface{}]interface{}) error {
	for k, v := range source {
		ks := DecodeString(k, "")
		if ks != "" {
			(*m)[ks] = v
		}
	}
	return nil
}

func (m *SortedMap) Keys() strings.StringSlice {
	var a strings.StringSlice
	for k, _ := range *m {
		a = append(a, k)
	}
	return a
}

func (m *SortedMap) ForEach(f func(string, interface{}) error) error {
	return m.Keys().
		Sort().
		ForEach(func(k string) error {
			return f(k, (*m)[k])
		})
}
