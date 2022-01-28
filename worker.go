package kernel

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
	"log"
	"time"
)

type Worker struct {
	daemon *Daemon `kernel:"inject"`
	tasks  task.Queue
}

func (w *Worker) Start() error {
	w.tasks = task.NewQueue()

	if w.daemon.IsDaemon() {
		go w.runDaemon()
	}
	return nil
}

// AddTask adds a task with priority 0
func (w *Worker) AddTask(task task.Task) task.Queue {
	return w.tasks.AddTask(task)
}

// AddPriorityTask adds a task with a specific priority.
// Tasks with a higher priority value will run AFTER those with a lower value.
func (w *Worker) AddPriorityTask(priority int, task task.Task) task.Queue {
	return w.tasks.AddPriorityTask(priority, task)
}

// RunTasks runs any queued tasks, returning the first Error or nil if all have run.
// this is provided to allow tasks to be run in the background if the main thread has been claimed
// e.g. the webserver is running.
func (w *Worker) runDaemon() {
	for {
		err := w.Run()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// Run kernel stage. This just calls RunTasks()
func (w *Worker) Run() error {
	return task.Run(w.tasks, context.Background())
}
