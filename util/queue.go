package util

// A queue is a slice which can have entries appended/removed from it.
// The type of queue, e.g. FIFO or LIFO is dependent on the implementation.
type Queue interface {
	Collection
	Iterable
	// Offer an entry to the queue. Returns true if the queue accepted the entry
	Offer(interface{}) bool
	// Poll removes the first entry in the queue or returns nil if empty
	Poll() interface{}
	// Peek returns the first entry like Poll() but does not remove the entry.
	Peek() interface{}
}
