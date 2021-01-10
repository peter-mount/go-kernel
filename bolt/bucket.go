package bolt

import (
	bbolt "github.com/etcd-io/bbolt"
)

type Bucket struct {
	tx     *Tx
	bucket *bbolt.Bucket
}

// Tx returns the tx of the bucket.
func (b *Bucket) Tx() *Tx {
	return b.tx
}

// Cursor creates a cursor associated with the bucket.
// The cursor is only valid as long as the transaction is open.
// Do not use a cursor after the transaction is closed.
func (b *Bucket) Cursor() *Cursor {
	// Allocate and return a cursor.
	return &Cursor{bucket: b, cursor: b.bucket.Cursor()}
}

// Bucket retrieves a nested bucket by name.
// Returns nil if the bucket does not exist.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) Bucket(name string) *Bucket {
	nb := b.bucket.Bucket([]byte(name))
	if nb == nil {
		return nil
	}
	return &Bucket{tx: b.tx, bucket: nb}
}

// CreateBucket creates a new bucket at the given key and returns the new bucket.
// Returns an error if the key already exists, if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) CreateBucket(name string) (*Bucket, error) {
	nb, err := b.bucket.CreateBucket([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Bucket{tx: b.tx, bucket: nb}, nil
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist and returns a reference to it.
// Returns an error if the bucket name is blank, or if the bucket name is too long.
// The bucket instance is only valid for the lifetime of the transaction.
func (b *Bucket) CreateBucketIfNotExists(name string) (*Bucket, error) {
	nb, err := b.bucket.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Bucket{tx: b.tx, bucket: nb}, nil
}

// DeleteBucket deletes a bucket at the given key.
// Returns an error if the bucket does not exists, or if the key represents a non-bucket value.
func (b *Bucket) DeleteBucket(name string) error {
	return b.bucket.DeleteBucket([]byte(name))
}

// Get retrieves the value for a key in the bucket.
// Returns a nil value if the key does not exist or if the key is a nested bucket.
// The returned value is only valid for the life of the transaction.
func (b *Bucket) Get(key string) []byte {
	return b.bucket.Get([]byte(key))
}

// Put sets the value for a key in the bucket.
// If the key exist then its previous value will be overwritten.
// Supplied value must remain valid for the life of the transaction.
// Returns an error if the bucket was created from a read-only transaction, if the key is blank, if the key is too large, or if the value is too large.
func (b *Bucket) Put(key string, value []byte) error {
	return b.bucket.Put([]byte(key), value)
}

// Delete removes a key from the bucket.
// If the key does not exist then nothing is done and a nil error is returned.
// Returns an error if the bucket was created from a read-only transaction.
func (b *Bucket) Delete(key string) error {
	return b.bucket.Delete([]byte(key))
}

// Sequence returns the current integer for the bucket without incrementing it.
func (b *Bucket) Sequence() uint64 {
	return b.bucket.Sequence()
}

// SetSequence updates the sequence number for the bucket.
func (b *Bucket) SetSequence(v uint64) error {
	return b.bucket.SetSequence(v)
}

// NextSequence returns an autoincrementing integer for the bucket.
func (b *Bucket) NextSequence() (uint64, error) {
	seq, err := b.bucket.NextSequence()
	return seq, err
}

// ForEach executes a function for each key/value pair in a bucket.
// If the provided function returns an error then the iteration is stopped and
// the error is returned to the caller. The provided function must not modify
// the bucket; this will result in undefined behavior.
func (b *Bucket) ForEach(fn func(k string, v []byte) error) error {
	return b.bucket.ForEach(func(k, v []byte) error {
		return fn(string(k[:]), v)
	})
}
