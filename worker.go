package kernel

import (
	"context"
	"github.com/peter-mount/go-kernel/util/task"
)

type Worker struct {
	tasks task.Queue
}

func (w *Worker) Name() string {
	return "worker"
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

// RunTasks runs any queued tasks, returning the first error or nil if all have run.
// this is provided to allow tasks to be run in the background if the main thread has been claimed
// e.g. the webserver is running.
func (w *Worker) RunTasks() error {
	return task.Run(w.tasks, context.Background())
}

// Run kernel stage. This just calls RunTasks()
func (w *Worker) Run() error {
	return w.RunTasks()
}
