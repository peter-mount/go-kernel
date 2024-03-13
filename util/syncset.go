package util

import "sync"

// syncSet is a synchronized set with accessors similar to the Java Set interface
type syncSet[T comparable] struct {
	mutex sync.Mutex
	m     map[T]interface{}
}

// NewSyncSet creates a new Synchronous Set
func NewSyncSet[T comparable]() Set[T] {
	s := &syncSet[T]{}
	s.Clear()
	return s
}

func (m *syncSet[T]) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.m)
}

func (m *syncSet[T]) IsEmpty() bool {
	return m.Size() == 0
}

func (m *syncSet[T]) Add(v T) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, e := m.m[v]
	if !e {
		m.m[v] = nil
	}
	return !e
}

func (m *syncSet[T]) AddAll(v ...T) bool {
	var r bool
	for _, e := range v {
		r = m.Add(e) || r
	}
	return r
}

func (m *syncSet[T]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.m = make(map[T]interface{})
}

func (m *syncSet[T]) Remove(k T) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *syncSet[T]) Contains(k T) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, exists := m.m[k]
	return exists
}

func (m *syncSet[T]) Slice() []T {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var a []T
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *syncSet[T]) ForEach(f func(T)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, _ := range m.m {
		f(k)
	}
}

func (m *syncSet[T]) ForEachFailFast(f func(T) error) error {
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

func (m *syncSet[T]) ForEachAsync(f func(T)) {
	for _, k := range m.Slice() {
		f(k)
	}
}

func (m *syncSet[T]) Iterator() Iterator[T] {
	a := m.Slice()
	return NewIterator[T](a...)
}

func (m *syncSet[T]) ReverseIterator() Iterator[T] {
	a := m.Slice()
	a = reverseSlice(a)
	return NewIterator(a...)
}
