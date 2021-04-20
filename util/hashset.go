package util

// hashSet is an un-synchronized set with accessors similar to the Java Set interface
type hashSet struct {
	m map[interface{}]interface{}
}

// NewHashSet creates a new Set. This set is not synchronised
func NewHashSet() Set {
	s := &hashSet{}
	s.Clear()
	return s
}

func (m *hashSet) Size() int {
	return len(m.m)
}

func (m *hashSet) IsEmpty() bool {
	return m.Size() == 0
}

func (m *hashSet) Add(v interface{}) bool {
	_, e := m.m[v]
	if !e {
		m.m[v] = nil
	}
	return !e
}

func (m *hashSet) AddAll(v ...interface{}) bool {
	var r bool
	for _, e := range v {
		r = m.Add(e) || r
	}
	return r
}

func (m *hashSet) Clear() {
	m.m = make(map[interface{}]interface{})
}

func (m *hashSet) Remove(k interface{}) bool {
	_, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *hashSet) Contains(k interface{}) bool {
	_, exists := m.m[k]
	return exists
}

func (m *hashSet) Slice() []interface{} {
	var a []interface{}
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *hashSet) ForEach(f func(interface{})) {
	for k, _ := range m.m {
		f(k)
	}
}

func (m *hashSet) ForEachFailFast(f func(interface{}) error) error {
	for k, _ := range m.m {
		err := f(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *hashSet) ForEachAsync(f func(interface{})) {
	for _, k := range m.Slice() {
		f(k)
	}
}

func (m *hashSet) Iterator() Iterator {
	a := m.Slice()
	return NewIterator(a...)
}

func (m *hashSet) ReverseIterator() Iterator {
	a := m.Slice()
	a = reverseSlice(a)
	return NewIterator(a...)
}
