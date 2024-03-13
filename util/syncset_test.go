package util

import "testing"

func TestSyncSet_Add(t *testing.T) {
	s := NewSyncSet[string]()

	val := "testValue"

	if !s.Add(val) {
		t.Errorf("Failed to add value")
	}

	if s.Add(val) {
		t.Errorf("Added value twice")
	}

	if s.IsEmpty() {
		t.Errorf("Set empty")
	}

	if s.Size() != 1 {
		t.Errorf("Set not containing correct number of values")
	}
}

func TestSyncSet_Remove(t *testing.T) {
	s := NewSyncSet[string]()

	val := "testValue"

	if !s.Add(val) {
		t.Errorf("Failed to add value")
	}

	if s.IsEmpty() {
		t.Errorf("Set empty")
	}

	if s.Size() != 1 {
		t.Errorf("Set not containing correct number of values")
	}

	if !s.Remove(val) {
		t.Errorf("Failed to remnove value")
	}

	if !s.IsEmpty() {
		t.Errorf("Set not empty")
	}

	if s.Size() != 0 {
		t.Errorf("Set not containing correct number of values")
	}
}
