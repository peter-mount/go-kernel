package util

// Set is a go equivalent to the Java Set interface
type Set interface {
	Collection
	Iterable
	// Add adds an entry to the set
	Add(interface{}) bool
	// Contains returns true if the set contains the value
	Contains(interface{}) bool
	// Remove removes the supplied entry
	Remove(interface{}) bool
	// Slice returns a slice containing a snapshot of all entries in the set
	Slice() []interface{}
}
