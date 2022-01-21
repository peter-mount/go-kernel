package util

import "testing"

func TestPriorityQueue_Add(t *testing.T) {
	var queue PriorityQueue

	if !queue.IsEmpty() {
		t.Errorf("PriorityQueue not empty on creation")
	}

	queue.Add(0)

	if queue.IsEmpty() {
		t.Errorf("PriorityQueue empty after adding")
	}
}

func TestPriorityQueue_AddPriority(t *testing.T) {
	var queue PriorityQueue

	if !queue.IsEmpty() {
		t.Errorf("PriorityQueue not empty on creation")
	}

	// Add entries 1..100 4 times, so we should have a set of entries in sorted order
	for j := 0; j < 4; j++ {
		for i := 1; i < 100; i++ {
			queue.AddPriority(i, i)
		}
	}

	// This should get added to the front
	queue.Add(0)

	if queue.IsEmpty() {
		t.Errorf("PriorityQueue empty after adding")
	}

	v := 0

	for !queue.IsEmpty() {
		val, ok := queue.Pop()
		if !ok {
			t.Errorf("PriorityQueue empty whilst expecting an entry")
		}

		if e, ok := val.(int); ok {
			if e < v || e > (v+1) {
				t.Errorf("PriorityQueue out of sequence, expected %d or %d, got %d",
					v, v+1, e)
			}
			v = e
		} else {
			t.Errorf("PriorityQueue non-int returned, expected %d or %d, got %v",
				v, v+1, e)
		}
	}

}
