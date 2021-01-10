package util

import "testing"

func TestSlice_IsEmpty(t *testing.T) {
	list := NewList()

	if !list.IsEmpty() {
		t.Errorf("List not empty")
	}

	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.IsEmpty() {
		t.Errorf("List empty")
	}

	if list.Size() != 3 {
		t.Errorf("List size incorrect")
	}
}
