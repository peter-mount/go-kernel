package task

import (
	"context"
	"fmt"
)

// WithValue adds a named value to the context
func (a Task) WithValue(key string, value interface{}) Task {
	if value == nil {
		panic(key)
	}

	return func(ctx context.Context) error {
		return a(context.WithValue(ctx, key, value))
	}
}

// WithContext copies the specified keys from a source context.
// It's the equivalent of WithValue(key,srcCtx.Value(key))
func (a Task) WithContext(srcCtx context.Context, keys ...string) Task {
	t := a
	for _, key := range keys {
		t = t.WithValue(key, srcCtx.Value(key))
	}
	return t
}

// RequireValue ensures that a key is defined within the current Context.
func (a Task) RequireValue(key string) Task {
	return func(ctx context.Context) error {
		if ctx.Value(key) == nil {
			return fmt.Errorf("required context Value %q missing", key)
		}
		return a.Do(ctx)
	}
}

// ValueProvider provides a value at runtime, used with UsingValue
type ValueProvider func(context.Context) (interface{}, error)

// UsingValue adds a named value to the context before passing it to the parent task
func (a Task) UsingValue(key string, f ValueProvider) Task {
	return func(ctx context.Context) error {
		value, err := f(ctx)
		if err == nil {
			err = a.Do(context.WithValue(ctx, key, value))
		}
		return err
	}
}

// ContextProvider provides a context based on another, used by Using
type ContextProvider func(context.Context) (context.Context, error)

// Using allows a ContextProvider to be used in a task, an alternate method to UsingValue
func (a Task) Using(f ContextProvider) Task {
	return func(ctx context.Context) error {
		ctx2, err := f(ctx)
		if err == nil {
			if ctx2 != nil {
				err = a.Do(ctx2)
			} else {
				err = a.Do(ctx)
			}
		}
		return err
	}
}
