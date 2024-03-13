package util

type List[T any] interface {
	Collection
	Iterable[T]
	// Add an entry to the list, usually on the end however this is up to the implementation.
	Add(T) bool
	// AddIndex adds an entry at the specified position in the list
	AddIndex(int, T)
	// Contains returns true if the set contains the value
	Contains(T) bool
	// Get returns the element at a specific index
	Get(int) T
	// IndexOf returns the index of the first occurrence of the specified element or -1 if not present.
	IndexOf(T) int
	// Remove removes the supplied entry
	Remove(T) bool
	// RemoveIndex removes the entry at the specific index
	RemoveIndex(i int) bool
	// FindIndexOf returns the index of the first occurrence that matches the provided predicate
	FindIndexOf(Predicate[T]) int
	// AddAll adds all entries to the list
	AddAll(...T)
}
