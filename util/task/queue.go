package task

import (
	"context"
	"errors"
	"github.com/peter-mount/go-kernel/util"
)

type Queue interface {
	AddTask(t Task) Queue
	AddPriorityTask(priority int, task Task) Queue
}

const (
	ctxKey = "task.Queue"
)

type defaultQueue struct {
	tasks util.PriorityQueue
}

func NewQueue() Queue {
	return &defaultQueue{}
}

// AddTask appends a Task to be performed once all Handler's have run.
func (q *defaultQueue) AddTask(t Task) Queue {
	q.tasks.Add(t)
	return q
}

// AddPriorityTask appends a Task to be performed once all Handler's have run.
func (q *defaultQueue) AddPriorityTask(priority int, task Task) Queue {
	q.tasks.AddPriority(priority, task)
	return q
}

// GetQueue returns the Queue contained in this Context
func GetQueue(ctx context.Context) Queue {
	if tc, ok := ctx.Value(ctxKey).(Queue); ok {
		return tc
	}

	return nil
}

// Run runs all tasks in the Queue until either the queue is empty or a task returns an error
func Run(queue Queue, ctx context.Context) error {
	if q, ok := queue.(*defaultQueue); ok {
		// Ensure we have a reference to the Queue in the context
		if ctx.Value(ctxKey) == nil {
			ctx = context.WithValue(ctx, ctxKey, queue)
		}

		// Run each task in sequence until either an error or the queue is empty
		return q.tasks.Drain(func(i interface{}) error {
			return i.(Task).Do(ctx)
		})
	}

	return errors.New("unsupported Queue")
}
