package util

type List interface {
	Collection
	Iterable
	// Add an entry in the list
	Add(interface{}) bool
	// Add an entry at the specified position in the list
	AddIndex(int, interface{})
	// Contains returns true if the set contains the value
	Contains(interface{}) bool
	// Get returns the element at a specific index
	Get(int) interface{}
	// IndexOf returns the index of the first occurrence of the specified element or -1 if not present.
	IndexOf(interface{}) int
	// Remove removes the supplied entry
	Remove(interface{}) bool
	// RemoveIndex removes the entry at the specific index
	RemoveIndex(i int) bool
	// FindIndexOf returns the index of the first occurrence that matches the provided predicate
	FindIndexOf(Predicate) int
}
