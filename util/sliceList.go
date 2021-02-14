package util

import (
	"errors"
)

type sliceList struct {
	data []interface{}
}

func NewList() List {
	return &sliceList{}
}

func (s *sliceList) Clear() {
	s.data = nil
}

func (s *sliceList) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *sliceList) Size() int {
	return len(s.data)
}

func (s *sliceList) ForEach(f func(interface{})) {
	for _, v := range s.data {
		f(v)
	}
}

func (s *sliceList) ForEachFailFast(f func(interface{}) error) error {
	for _, v := range s.data {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sliceList) ForEachAsync(f func(interface{})) {
	s.ForEach(f)
}

func (s *sliceList) Add(v interface{}) bool {
	s.data = append(s.data, v)
	return true
}

func (s *sliceList) AddIndex(i int, v interface{}) {
	l := len(s.data)

	if i < 0 || i >= l {
		panic(errors.New("list index out of bounds"))
	}

	if i == 0 {
		s.data = append([]interface{}{v}, s.data...)
	} else if i == l {
		s.data = append(s.data, v)
	} else {
		s.data = append(s.data[:i+1], s.data[i:]...)
		s.data[i] = v
	}
}

func (s *sliceList) Contains(v interface{}) bool {
	for _, e := range s.data {
		if e == v {
			return true
		}
	}
	return false
}

func (s *sliceList) Get(i int) interface{} {
	if i < 0 || i >= len(s.data) {
		panic(errors.New("index out of bounds"))
	}
	return s.data[i]
}

func (s *sliceList) IndexOf(v interface{}) int {
	for i, e := range s.data {
		if e == v {
			return i
		}
	}
	return -1
}

func (s *sliceList) FindIndexOf(f Predicate) int {
	for i, e := range s.data {
		if f(e) {
			return i
		}
	}
	return -1
}

func (s *sliceList) Remove(v interface{}) bool {
	for i, e := range s.data {
		if e == v {
			return s.RemoveIndex(i)
		}
	}
	return false
}

func (s *sliceList) RemoveIndex(i int) bool {
	if i < 0 || i >= len(s.data) {
		return false
	}

	if i == 0 {
		if len(s.data) == 1 {
			s.data = nil
		} else {
			s.data = s.data[1:]
		}
	} else if i == (len(s.data) - 1) {
		s.data = s.data[:i]
	} else {
		s.data = append(s.data[:i], s.data[i+1:]...)
	}

	return true
}

func copySlice(s []interface{}) []interface{} {
	var a []interface{}
	if s != nil {
		a = make([]interface{}, len(s))
		copy(a, s)
	}
	return a
}

func reverseSlice(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s *sliceList) Iterator() Iterator {
	a := copySlice(s.data)
	return NewIterator(a...)
}

func (s *sliceList) ReverseIterator() Iterator {
	a := copySlice(s.data)
	a = reverseSlice(a)
	return NewIterator(a...)
}
