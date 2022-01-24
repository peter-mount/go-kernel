package test

import (
	"github.com/peter-mount/go-kernel"
	"testing"
)

type testDepService1 struct {
	Id int
}

type testDepService2 struct {
	Id int
	s1 *testDepService1 `kernel:"inject"`
}

type testDepService3 struct {
	Id int
	s1 *testDepService1 `kernel:"inject"`
}

// TestDependency_Existing deploys two services referencing a third.
//
// If the third is referenced from one of them, then it's state should be seen from
// the other one.
//
// This ensures that the common service is the same one and not a unique instance - i.e.
// the kernel is picking up the new object
func TestDependency_Existing(t *testing.T) {
	s2 := &testDepService2{Id: 42}
	s3 := &testDepService3{Id: 31415}

	err := kernel.Launch(s2, s3)
	if err != nil {
		t.Fatal(err)
	}

	val := 9876
	s2.s1.Id = val
	if s3.s1.Id != s2.s1.Id {
		t.Errorf("not same instance, %d!=%d expected %d", s2.s1.Id, s3.s1.Id, val)
	}
}
