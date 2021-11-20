package util

type PriorityQueue []PriorityEntry

type PriorityEntry struct {
	Priority int
	Element  interface{}
}

func (p *PriorityQueue) Add(e interface{}) {
	p.AddPriority(0, e)
}

func (p *PriorityQueue) AddPriority(priority int, e interface{}) {
	ent := PriorityEntry{Priority: priority, Element: e}

	for i, existing := range *p {
		if existing.Priority > priority {
			*p = append(*p, ent)
			copy((*p)[i+1:], (*p)[i:])
			(*p)[i] = ent
			return
		}
	}

	*p = append(*p, ent)
}

func (p *PriorityQueue) IsEmpty() bool {
	return len(*p) == 0
}

// ForEach will call a function for each entry in the PriorityQueue.
// The Queue is not modified during this call.
func (p *PriorityQueue) ForEach(f func(interface{}) error) error {
	for _, e := range *p {
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
	for !p.IsEmpty() {
		e, exists := p.Pop()
		if exists {
			err := f(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Pop removes the first entry from the queue.
func (p *PriorityQueue) Pop() (interface{}, bool) {
	if p.IsEmpty() {
		return nil, false
	}

	e := (*p)[0]
	*p = (*p)[1:]
	return e.Element, true
}
