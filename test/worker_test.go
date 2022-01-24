package test

import (
	"context"
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/util/task"
	"testing"
)

type workertestservice struct {
	queue    task.Queue `kernel:"worker"`
	res      int
	expected int
}

func (w *workertestservice) Name() string {
	return "workertestservice"
}

func (w *workertestservice) Start() error {
	w.add(10)
	w.add(50)
	w.add(12)
	return nil
}

func (w *workertestservice) add(priority int) {
	w.queue.AddPriorityTask(priority, func(_ context.Context) error {
		w.res += priority
		return nil
	})
	w.expected += priority
}

// TestService_Inject tests the kernel injects the common worker queue
func TestWorker_Inject(t *testing.T) {

	s := &workertestservice{}

	err := kernel.Launch(s)
	if err != nil {
		t.Fatal(err)
	}

	if s.queue == nil {
		t.Fatal("No injected service")
	}

	if s.res != s.expected {
		t.Errorf("Invalid result, expected %d got %d", s.expected, s.res)
	}
}
