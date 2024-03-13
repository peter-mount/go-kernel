package util

import "sync"

// SyncMap is a synchronized map with accessors similar to the Java Map interface
type syncMap[V comparable] struct {
	mutex sync.Mutex
	m     map[string]V
}

// NewSyncMap creates a new Synchronous Map
func NewSyncMap[V comparable]() Map[V] {
	m := &syncMap[V]{}
	m.Clear()
	return m
}

func (m *syncMap[V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m = make(map[string]V)
}

func (m *syncMap[V]) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.m)
}

func (m *syncMap[V]) IsEmpty() bool {
	return m.Size() == 0
}

func (m *syncMap[V]) Get(k string) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.m[k]
}

func (m *syncMap[V]) Get2(k string) (V, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, e := m.m[k]
	return v, e
}

func (m *syncMap[V]) Put(k string, v V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	old := m.m[k]
	m.m[k] = v
	return old
}

func (m *syncMap[V]) PutIfAbsent(k string, v V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists {
		return old
	}

	m.m[k] = v

	var nilResponse V
	return nilResponse
}
func (m *syncMap[V]) Remove(k string) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return old
	}

	var nilResponse V
	return nilResponse
}

func (m *syncMap[V]) RemoveIfEquals(k string, v V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists && old == v {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *syncMap[V]) ReplaceIfEquals(k string, oldValue V, newValue V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists && old == oldValue {
		m.m[k] = newValue
		return true
	}

	return false
}

func (m *syncMap[V]) Replace(k string, v V) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.m[k]
	if exists {
		m.m[k] = v
		return true
	}

	return false
}

func (m *syncMap[V]) ContainsKey(k string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, exists := m.m[k]
	return exists
}

func (m *syncMap[V]) ComputeIfAbsent(k string, f func(string) V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v, exists := m.m[k]
	if !exists {
		v = f(k)
		m.m[k] = v
	}

	return v
}

func (m *syncMap[V]) ComputeIfPresent(k string, f func(string, V) V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v, exists := m.m[k]
	if exists {
		v = f(k, v)
		m.m[k] = v
	}

	return v
}

func (m *syncMap[V]) Compute(k string, f func(string, V) V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v := f(k, m.m[k])
	m.m[k] = v

	return v
}

func (m *syncMap[V]) Merge(k string, v V, f func(V, V) V) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ev, exists := m.m[k]

	if !exists {
		ev = v
	} else {
		ev = f(ev, v)
	}

	m.m[k] = ev
	return ev
}

func (m *syncMap[V]) ExecIfPresent(k string, f func(V)) {
	v, e := m.Get2(k)
	if e {
		f(v)
	}
}

func (m *syncMap[V]) ForEach(f func(string, V)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.m {
		f(k, v)
	}
}

func (m *syncMap[V]) Keys() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var a []string
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *syncMap[V]) ForEachAsync(f func(string, V)) {
	for _, k := range m.Keys() {
		if v, e := m.Get2(k); e {
			f(k, v)
		}
	}
}
