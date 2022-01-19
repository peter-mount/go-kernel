package util

import "sync"

type PriorityQueue struct {
	mutex   sync.Mutex
	entries []PriorityEntry
}

type PriorityEntry struct {
	Priority int
	Element  interface{}
}

func (p *PriorityQueue) Add(e interface{}) {
	p.AddPriority(0, e)
}

func (p *PriorityQueue) AddPriority(priority int, e interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ent := PriorityEntry{Priority: priority, Element: e}

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

func (p *PriorityQueue) IsEmpty() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.entries) == 0
}

// clone returns a clone of the entry list. Used to take a snapshot, see ForEach
func (p *PriorityQueue) clone() []PriorityEntry {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	e := make([]PriorityEntry, len(p.entries))
	copy(e, p.entries)

	return e
}

// ForEach will call a function for each entry in the PriorityQueue.
// The Queue is not modified during this call.
func (p *PriorityQueue) ForEach(f func(interface{}) error) error {
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
func (p *PriorityQueue) Drain(f func(interface{}) error) error {
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
func (p *PriorityQueue) Pop() (interface{}, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Don't use p.IsEmpty() to save a round of locking
	if len(p.entries) == 0 {
		return nil, false
	}

	e := p.entries[0]
	p.entries = p.entries[1:]
	return e.Element, true
}
