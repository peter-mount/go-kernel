package util

import "sync"

type PriorityQueue[T any] struct {
	mutex   sync.Mutex
	entries []PriorityEntry[T]
}

type PriorityEntry[T any] struct {
	Priority int
	Element  T
}

func (p *PriorityQueue[T]) Add(e T) {
	p.AddPriority(0, e)
}

func (p *PriorityQueue[T]) AddPriority(priority int, e T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ent := PriorityEntry[T]{Priority: priority, Element: e}

	for i, existing := range p.entries {
		if existing.Priority > priority {
			p.entries = append(p.entries, ent)
			copy(p.entries[i+1:], p.entries[i:])
			p.entries[i] = ent
			return
		}
	}

	p.entries = append(p.entries, ent)
}

func (p *PriorityQueue[T]) IsEmpty() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.entries) == 0
}

// clone returns a clone of the entry list. Used to take a snapshot, see ForEach
func (p *PriorityQueue[T]) clone() []PriorityEntry[T] {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	e := make([]PriorityEntry[T], len(p.entries))
	copy(e, p.entries)

	return e
}

// ForEach will call a function for each entry in the PriorityQueue.
// The Queue is not modified during this call.
func (p *PriorityQueue[T]) ForEach(f func(T) error) error {
	// Run on a clone() so the underlying queue can change whilst we are running
	for _, e := range p.clone() {
		err := f(e.Element)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drain will call a function for each entry in the PriorityQueue,
// removing the entry from the head of the queue.
func (p *PriorityQueue[T]) Drain(f func(T) error) error {
	for {
		e, exists := p.Pop()
		if !exists {
			return nil
		}

		err := f(e)
		if err != nil {
			return err
		}
	}
}

// Pop removes the first entry from the queue.
func (p *PriorityQueue[T]) Pop() (T, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Don't use p.IsEmpty() to save a round of locking
	if len(p.entries) == 0 {
		// Before generics, we could return nil, but this is the only way with generics
		var emptyVal T
		return emptyVal, false
	}

	e := p.entries[0]
	p.entries = p.entries[1:]
	return e.Element, true
}
