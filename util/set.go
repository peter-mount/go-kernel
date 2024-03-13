package util

// Set is a go equivalent to the Java Set interface
type Set[T any] interface {
	Collection
	Iterable[T]
	// Add adds an entry to the set
	Add(T) bool
	// AddAll adds all supplied values to the set
	AddAll(v ...T) bool
	// Contains returns true if the set contains the value
	Contains(T) bool
	// Remove removes the supplied entry
	Remove(T) bool
	// Slice returns a slice containing a snapshot of all entries in the set
	Slice() []T
}
