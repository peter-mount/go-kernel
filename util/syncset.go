package util

import "sync"

// syncSet is a synchronized set with accessors similar to the Java Set interface
type syncSet struct {
	mutex sync.Mutex
	m     map[interface{}]interface{}
}

// NewSyncSet creates a new Synchronous Set
func NewSyncSet() Set {
	s := &syncSet{}
	s.Clear()
	return s
}

func (m *syncSet) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.m)
}

func (m *syncSet) IsEmpty() bool {
	return m.Size() == 0
}

func (m *syncSet) Add(v interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, e := m.m[v]
	if !e {
		m.m[v] = nil
	}
	return !e
}

func (m *syncSet) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.m = make(map[interface{}]interface{})
}

func (m *syncSet) Remove(k interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *syncSet) Contains(k interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, exists := m.m[k]
	return exists
}

func (m *syncSet) Slice() []interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var a []interface{}
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *syncSet) ForEach(f func(interface{})) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, _ := range m.m {
		f(k)
	}
}

func (m *syncSet) ForEachFailFast(f func(interface{}) error) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, _ := range m.m {
		err := f(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *syncSet) ForEachAsync(f func(interface{})) {
	for _, k := range m.Slice() {
		f(k)
	}
}

func (m *syncSet) Iterator() Iterator {
	a := m.Slice()
	return NewIterator(a...)
}

func (m *syncSet) ReverseIterator() Iterator {
	a := m.Slice()
	a = reverseSlice(a)
	return NewIterator(a...)
}
