package util

import "errors"

type Iterable interface {
	// ForEach calls a function for each entry in the collection.
	// WARNING: For synchronized sets the function is called from inside the lock so the set cannot be modified for the
	// duration of this call. See ForEachAsync
	ForEach(f func(interface{}))
	// ForEachAsync calls a function for each entry in the collection.
	// WARNING: This can be expensive as it calls Slice() first to get a copy of the set before iterating over the
	// entries. However, it does mean the set can be modified during the call.
	ForEachAsync(f func(interface{}))
	// ForEachFailFast calls a function for each entry in the collection.
	// WARNING: For synchronized sets the function is called from inside the lock so the set cannot be modified for the
	// duration of this call. See ForEachAsync
	ForEachFailFast(f func(interface{}) error) error
	Iterator() Iterator
	ReverseIterator() Iterator
}

type Iterator interface {
	Iterable
	HasNext() bool
	Next() interface{}
}

type sliceIterator struct {
	slice []interface{}
	pos   int
}

func NewIterator(v ...interface{}) Iterator {
	return &sliceIterator{slice: v}
}

func (i *sliceIterator) HasNext() bool {
	return i.pos < len(i.slice)
}

func (i *sliceIterator) Next() interface{} {
	if i.HasNext() {
		v := i.slice[i.pos]
		i.pos++
		return v
	}
	panic(errors.New("iterator out of bounds"))
}

func (i *sliceIterator) ForEach(f func(interface{})) {
	for _, v := range i.slice {
		f(v)
	}
}

func (i *sliceIterator) ForEachAsync(f func(interface{})) {
	i.ForEach(f)
}

func (i *sliceIterator) ForEachFailFast(f func(interface{}) error) error {
	for _, v := range i.slice {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *sliceIterator) Iterator() Iterator {
	// Just return ourselves
	return i
}

func (i *sliceIterator) ReverseIterator() Iterator {
	// Just return ourselves
	return i
}
