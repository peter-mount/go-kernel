package util

import (
	"errors"
	"slices"
)

type Iterable[T any] interface {
	// ForEach calls a function for each entry in the collection.
	// WARNING: For synchronized sets the function is called from inside the lock so the set cannot be modified for the
	// duration of this call. See ForEachAsync
	ForEach(f func(T))
	// ForEachAsync calls a function for each entry in the collection.
	// WARNING: This can be expensive as it calls Slice() first to get a copy of the set before iterating over the
	// entries. However, it does mean the set can be modified during the call.
	ForEachAsync(f func(T))
	// ForEachFailFast calls a function for each entry in the collection.
	// WARNING: For synchronized sets the function is called from inside the lock so the set cannot be modified for the
	// duration of this call. See ForEachAsync
	ForEachFailFast(f func(T) error) error
	Iterator() Iterator[T]
	ReverseIterator() Iterator[T]
}

type Iterator[T any] interface {
	Iterable[T]
	HasNext() bool
	Next() T
}

type sliceIterator[T any] struct {
	slice []T
	pos   int
}

func NewIterator[T any](v ...T) Iterator[T] {
	return &sliceIterator[T]{slice: v}
}

func (i *sliceIterator[T]) HasNext() bool {
	return i.pos < len(i.slice)
}

func (i *sliceIterator[T]) Next() T {
	if i.HasNext() {
		v := i.slice[i.pos]
		i.pos++
		return v
	}
	panic(errors.New("iterator out of bounds"))
}

func (i *sliceIterator[T]) ForEach(f func(T)) {
	for _, v := range i.slice {
		f(v)
	}
}

func (i *sliceIterator[T]) ForEachAsync(f func(T)) {
	i.ForEach(f)
}

func (i *sliceIterator[T]) ForEachFailFast(f func(T) error) error {
	for _, v := range i.slice {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *sliceIterator[T]) Iterator() Iterator[T] {
	s := slices.Clone(i.slice)
	return NewIterator[T](s...)
}

func (i *sliceIterator[T]) ReverseIterator() Iterator[T] {
	s := slices.Clone(i.slice)
	slices.Reverse(s)
	return NewIterator[T](s...)
}
