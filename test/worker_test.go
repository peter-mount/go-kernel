package test

import (
	"context"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"log"
	"sync/atomic"
	"testing"
)

type workertestservice struct {
	queue    task.Queue `kernel:"worker"`
	res      int
	expected int
	ctr      uint64
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

func TestWorker_Thread(t *testing.T) {
	err := kernel.Launch(&testThread{t: t})
	if err != nil {
		t.Fatal(err)
	}
}

type testThread struct {
	daemon *kernel.Daemon `kernel:"inject"`
	worker task.Queue     `kernel:"worker"`
	t      *testing.T
	count  uint64
	total  uint64
}

func (t *testThread) PostInit() error {
	// Mark ourselves as a daemon
	t.daemon.SetDaemon()
	return nil
}

func (t *testThread) Start() error {
	//go t.test()

	for i := 0; i < 10; i++ {
		t.add(i)
	}
	return nil
}

func (t *testThread) add(i int) {
	nt := atomic.AddUint64(&t.total, 1)
	log.Printf("Add %d %d", nt, i)
	t.worker.AddTask(task.Of(t.runtest).WithValue("i", i))
}

func (t *testThread) runtest(ctx context.Context) error {
	log.Println("run")
	nc := atomic.AddUint64(&t.count, 1)
	i := ctx.Value("i").(int)
	log.Printf("Run %d/%d i=%d", nc, t.total, i)

	if i < 7 {
		t.add(i + 10)
	}
	if i == 8 {
		t.worker.AddPriorityTask(100, func(ctx context.Context) error {
			t.daemon.ClearDaemon()
			return nil
		})
	}
	return nil
}

func (t *testThread) test() {
}
