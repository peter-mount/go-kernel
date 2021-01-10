package bolt

import (
	bbolt "github.com/etcd-io/bbolt"
)

type Tx struct {
	tx *bbolt.Tx
}

// Update executes a function within the context of a read-write managed transaction.
// If no error is returned from the function then the transaction is committed.
// If an error is returned then the entire transaction is rolled back.
// Any error that is returned from the function or returned from the commit is
// returned from the Update() method.
func (s *BoltService) Update(fn func(*Tx) error) error {
	return s.db.Update(func(btx *bbolt.Tx) error {
		return fn(&Tx{tx: btx})
	})
}

// View executes a function within the context of a managed read-only transaction.
// Any error that is returned from the function is returned from the View() method.
func (s *BoltService) View(fn func(*Tx) error) error {
	return s.db.View(func(btx *bbolt.Tx) error {
		return fn(&Tx{tx: btx})
	})
}

// Batch calls fn as part of a batch. It behaves similar to Update,
// except:
//
// 1. concurrent Batch calls can be combined into a single Bolt
// transaction.
//
// 2. the function passed to Batch may be called multiple times,
// regardless of whether it returns error or not.
//
// This means that Batch function side effects must be idempotent and
// take permanent effect only after a successful return is seen in
// caller.
//
// The maximum batch size and delay can be adjusted with DB.MaxBatchSize
// and DB.MaxBatchDelay, respectively.
//
// Batch is only useful when there are multiple goroutines calling it.
func (s *BoltService) Batch(fn func(*Tx) error) error {
	return s.db.Batch(func(btx *bbolt.Tx) error {
		return fn(&Tx{tx: btx})
	})
}

// Bucket retrieves a bucket by name.
// Returns nil if the bucket does not exist.
// The bucket instance is only valid for the lifetime of the transaction.
func (t *Tx) Bucket(name string) *Bucket {
	b := t.tx.Bucket([]byte(name))
	if b == nil {
		return nil
	}
	return &Bucket{tx: t, bucket: b}
}

// CreateBucket creates a new bucket.
// Returns an error if the bucket already exists, if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (t *Tx) CreateBucket(name string) (*Bucket, error) {
	b, err := t.tx.CreateBucket([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Bucket{tx: t, bucket: b}, nil
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist.
// Returns an error if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (t *Tx) CreateBucketIfNotExists(name string) (*Bucket, error) {
	b, err := t.tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Bucket{tx: t, bucket: b}, nil
}

// DeleteBucket deletes a bucket.
// Returns an error if the bucket cannot be found or if the key represents a non-bucket value.
func (t *Tx) DeleteBucket(name string) error {
	return t.tx.DeleteBucket([]byte(name))
}

// OnCommit adds a handler function to be executed after the transaction successfully commits.
func (t *Tx) OnCommit(fn func()) {
	t.tx.OnCommit(fn)
}

// ForEach iterates over all bucket names
func (t *Tx) ForEach(fn func(k string, v *Bucket) error) error {
	return t.tx.ForEach(func(k []byte, b *bbolt.Bucket) error {
		return fn(string(k[:]), &Bucket{tx: t, bucket: b})
	})
}
