package task

import (
	"context"
)

// Task is a task that the Generator must run once all other Handler's have been run.
// They are usually tasks created by those Handlers.
type Task func(ctx context.Context) error

// Of creates a new Task forming a chain of the provided tasks
func Of(tasks ...Task) Task {
	var task Task
	for _, b := range tasks {
		task = task.Then(b)
	}
	return task
}

// Then joins two tasks together
func (a Task) Then(b Task) Task {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(ctx context.Context) error {
		err := a(ctx)
		if err == nil {
			err = b(ctx)
		}
		return err
	}
}

func (a Task) Do(ctx context.Context) error {
	if a != nil {
		return a(ctx)
	}
	return nil
}

// RunOnce will invoke a task exactly once.
// It uses a pointer to a boolean to store this state.
// It's useful for simple tasks but should be treated as Deprecated.
// Currently, here as Book still uses it as it only works for one Book not multiple books.
func (a Task) RunOnce(f *bool, t Task) Task {
	return a.Then(func(ctx context.Context) error {
		if !*f {
			*f = true
			return t(ctx)
		}
		return nil
	})
}

// Queue will defer the queued task onto the underlying Queue.
// If one or more tasks are provided then they will be queued if the flow reaches this location.
// If none, then the current task will be queued when run.
func (a Task) Queue(tasks ...Task) Task {
	return a.QueueWithPriority(0, tasks...)
}

// QueueWithPriority will defer the queued task onto the underlying Queue with a priority
// If one or more tasks are provided then they will be queued if the flow reaches this location.
// If none, then the current task will be queued when run.
func (a Task) QueueWithPriority(priority int, tasks ...Task) Task {
	if len(tasks) == 0 {
		return func(ctx context.Context) error {
			GetQueue(ctx).AddPriorityTask(priority, a)
			return nil
		}
	}

	r := a
	for _, task := range tasks {
		r = r.Then(task.Queue())
	}
	return r
}

// Guard wraps a task so that any error or panic returned by that task is ignored.
// It is used when you don't want the task to stop all other processing.
func (a Task) Guard() Task {
	return func(ctx context.Context) error {
		defer func() {
			_ = recover()
		}()

		_ = a.Do(ctx)
		return nil
	}
}
