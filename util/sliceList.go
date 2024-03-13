package util

import (
	"errors"
)

type sliceList[T comparable] struct {
	data []T
}

func NewList[T comparable](v ...T) List[T] {
	l := &sliceList[T]{}
	if len(v) > 0 {
		l.AddAll(v...)
	}
	return l
}

func (s *sliceList[T]) Clear() {
	s.data = nil
}

func (s *sliceList[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *sliceList[T]) Size() int {
	return len(s.data)
}

func (s *sliceList[T]) ForEach(f func(T)) {
	for _, v := range s.data {
		f(v)
	}
}

func (s *sliceList[T]) ForEachFailFast(f func(T) error) error {
	for _, v := range s.data {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sliceList[T]) ForEachAsync(f func(T)) {
	s.ForEach(f)
}

func (s *sliceList[T]) Add(v T) bool {
	s.data = append(s.data, v)
	return true
}

func (s *sliceList[T]) AddAll(v ...T) {
	s.data = append(s.data, v...)
}

func (s *sliceList[T]) AddIndex(i int, v T) {
	l := len(s.data)

	if i < 0 || i >= l {
		panic(errors.New("list index out of bounds"))
	}

	if i == 0 {
		s.data = append([]T{v}, s.data...)
	} else if i == l {
		s.data = append(s.data, v)
	} else {
		s.data = append(s.data[:i+1], s.data[i:]...)
		s.data[i] = v
	}
}

func (s *sliceList[T]) Contains(v T) bool {
	for _, e := range s.data {
		if e == v {
			return true
		}
	}
	return false
}

func (s *sliceList[T]) Get(i int) T {
	if i < 0 || i >= len(s.data) {
		panic(errors.New("index out of bounds"))
	}
	return s.data[i]
}

func (s *sliceList[T]) IndexOf(v T) int {
	for i, e := range s.data {
		if e == v {
			return i
		}
	}
	return -1
}

func (s *sliceList[T]) FindIndexOf(f Predicate[T]) int {
	for i, e := range s.data {
		if f(e) {
			return i
		}
	}
	return -1
}

func (s *sliceList[T]) Remove(v T) bool {
	for i, e := range s.data {
		if e == v {
			return s.RemoveIndex(i)
		}
	}
	return false
}

func (s *sliceList[T]) RemoveIndex(i int) bool {
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

func copySlice[T any](s []T) []T {
	var a []T
	if s != nil {
		a = make([]T, len(s))
		copy(a, s)
	}
	return a
}

func reverseSlice[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s *sliceList[T]) Iterator() Iterator[T] {
	a := copySlice(s.data)
	return NewIterator[T](a...)
}

func (s *sliceList[T]) ReverseIterator() Iterator[T] {
	a := copySlice(s.data)
	a = reverseSlice(a)
	return NewIterator[T](a...)
}
