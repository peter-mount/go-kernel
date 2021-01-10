package util

import "sync"

// SyncMap is a synchronized map with accessors similar to the Java Map interface
type syncMap struct {
	mutex sync.Mutex
	m     map[string]interface{}
}

// NewSyncMap creates a new Synchronous Map
func NewSyncMap() Map {
	m := &syncMap{}
	m.Clear()
	return m
}

func (m *syncMap) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m = make(map[string]interface{})
}

func (m *syncMap) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.m)
}

func (m *syncMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *syncMap) Get(k string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.m[k]
}

func (m *syncMap) Get2(k string) (interface{}, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	v, e := m.m[k]
	return v, e
}

func (m *syncMap) Put(k string, v interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	old := m.m[k]
	m.m[k] = v
	return old
}

func (m *syncMap) PutIfAbsent(k string, v interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists {
		return old
	}

	m.m[k] = v
	return nil
}
func (m *syncMap) Remove(k string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return old
	}

	return nil
}

func (m *syncMap) RemoveIfEquals(k string, v interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists && old == v {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *syncMap) ReplaceIfEquals(k string, oldValue interface{}, newValue interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	old, exists := m.m[k]
	if exists && old == oldValue {
		m.m[k] = newValue
		return true
	}

	return false
}

func (m *syncMap) Replace(k string, v interface{}) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.m[k]
	if exists {
		m.m[k] = v
		return true
	}

	return false
}

func (m *syncMap) ContainsKey(k string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, exists := m.m[k]
	return exists
}

func (m *syncMap) ComputeIfAbsent(k string, f func(string) interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v, exists := m.m[k]
	if !exists {
		v = f(k)
		m.m[k] = v
	}

	return v
}

func (m *syncMap) ComputeIfPresent(k string, f func(string, interface{}) interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v, exists := m.m[k]
	if exists {
		v = f(k, v)
		m.m[k] = v
	}

	return v
}

func (m *syncMap) Compute(k string, f func(string, interface{}) interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	v := f(k, m.m[k])
	m.m[k] = v

	return v
}

func (m *syncMap) Merge(k string, v interface{}, f func(interface{}, interface{}) interface{}) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ev, exists := m.m[k]

	if !exists || ev == nil {
		ev = v
	} else {
		ev = f(ev, v)
	}

	m.m[k] = ev
	return ev
}

func (m *syncMap) ExecIfPresent(k string, f func(interface{})) {
	v, e := m.Get2(k)
	if e {
		f(v)
	}
}

func (m *syncMap) ForEach(f func(string, interface{})) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.m {
		f(k, v)
	}
}

func (m *syncMap) Keys() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var a []string
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *syncMap) ForEachAsync(f func(string, interface{})) {
	for _, k := range m.Keys() {
		if v, e := m.Get2(k); e {
			f(k, v)
		}
	}
}
