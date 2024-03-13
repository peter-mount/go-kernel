package kernel

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"time"
)

type Worker struct {
	daemon *Daemon `kernel:"inject"`
	tasks  util.PriorityQueue[task.Task]
}

// AddTask adds a task with priority 0
func (w *Worker) AddTask(task task.Task) task.Queue {
	w.tasks.Add(task)
	return w
}

// AddPriorityTask adds a task with a specific priority.
// Tasks with a higher priority value will run AFTER those with a lower value.
func (w *Worker) AddPriorityTask(priority int, task task.Task) task.Queue {
	w.tasks.AddPriority(priority, task)
	return w
}

func (w *Worker) Start() error {
	// If in webserver mode then run tasks in the background
	if w.daemon.IsWebserver() {
		go func() {
			for {
				_ = w.runDaemon()
			}
		}()
	}
	return nil
}

// Run kernel stage. This just calls RunTasks()
func (w *Worker) Run() error {
	if !w.daemon.IsWebserver() {
		return w.runDaemon()
	}
	return nil
}

func (w *Worker) runDaemon() error {
	run := true
	for run {
		if err := w.run(context.Background()); err != nil {
			return err
		}
		run = w.daemon.IsDaemon()
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

func (w *Worker) run(ctx context.Context) error {
	// Ensure we have a reference to the Queue in the context
	ctx = context.WithValue(ctx, ctxKey, w)

	// Run each task in sequence until either an error or the queue is empty
	return w.tasks.Drain(func(t task.Task) error {
		return t.Do(ctx)
	})
}

const (
	ctxKey = "task.Queue"
)
