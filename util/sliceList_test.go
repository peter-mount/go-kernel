package util

import "testing"

func TestSlice_IsEmpty(t *testing.T) {
	list := NewList[int]()

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

func TestSlice_Remove(t *testing.T) {
	list := NewList[int]()

	v1 := 1
	v2 := 2
	v3 := 3

	list.Add(v1)
	list.Add(v2)
	list.Add(v3)

	if list.Size() != 3 {
		t.Errorf("List not correct size, expected 3 got %d", list.Size())
	}

	list.Remove(v2)

	if list.Size() != 2 {
		t.Errorf("List not correct size, expected 2 got %d", list.Size())
	}

	if list.Get(0) != v1 {
		t.Errorf("Element 0 returned %v not %v", list.Get(0), v1)
	}

	if list.Get(1) != v3 {
		t.Errorf("Element 1 returned %v not %v", list.Get(1), v3)
	}
}
