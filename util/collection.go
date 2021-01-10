package util

type Collection interface {
	// Clear removes all entries from the set
	Clear()
	// IsEmpty returns true if the set is empty
	IsEmpty() bool
	// Size returns the size of the set
	Size() int
}
