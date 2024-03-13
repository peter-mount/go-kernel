package util

// hashSet is an un-synchronized set with accessors similar to the Java Set interface
type hashSet[T comparable] struct {
	m map[T]interface{}
}

// NewHashSet creates a new Set. This set is not synchronised
func NewHashSet[T comparable]() Set[T] {
	s := &hashSet[T]{}
	s.Clear()
	return s
}

func (m *hashSet[T]) Size() int {
	return len(m.m)
}

func (m *hashSet[T]) IsEmpty() bool {
	return m.Size() == 0
}

func (m *hashSet[T]) Add(v T) bool {
	_, e := m.m[v]
	if !e {
		m.m[v] = nil
	}
	return !e
}

func (m *hashSet[T]) AddAll(v ...T) bool {
	var r bool
	for _, e := range v {
		r = m.Add(e) || r
	}
	return r
}

func (m *hashSet[T]) Clear() {
	m.m = make(map[T]interface{})
}

func (m *hashSet[T]) Remove(k T) bool {
	_, exists := m.m[k]
	if exists {
		delete(m.m, k)
		return true
	}

	return false
}

func (m *hashSet[T]) Contains(k T) bool {
	_, exists := m.m[k]
	return exists
}

func (m *hashSet[T]) Slice() []T {
	var a []T
	for k, _ := range m.m {
		a = append(a, k)
	}
	return a
}

func (m *hashSet[T]) ForEach(f func(T)) {
	for k, _ := range m.m {
		f(k)
	}
}

func (m *hashSet[T]) ForEachFailFast(f func(T) error) error {
	for k, _ := range m.m {
		err := f(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *hashSet[T]) ForEachAsync(f func(T)) {
	for _, k := range m.Slice() {
		f(k)
	}
}

func (m *hashSet[T]) Iterator() Iterator[T] {
	a := m.Slice()
	return NewIterator[T](a...)
}

func (m *hashSet[T]) ReverseIterator() Iterator[T] {
	a := m.Slice()
	a = reverseSlice(a)
	return NewIterator[T](a...)
}
