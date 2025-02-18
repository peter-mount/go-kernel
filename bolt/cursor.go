package bolt

import (
	bbolt "go.etcd.io/bbolt"
)

type Cursor struct {
	bucket *Bucket
	cursor *bbolt.Cursor
}

// Bucket returns the bucket that this cursor was created from.
func (c *Cursor) Bucket() *Bucket {
	return c.bucket
}

// First moves the cursor to the first item in the bucket and returns its key and value.
// If the bucket is empty then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) First() (key string, value []byte) {
	k, v := c.cursor.First()
	return string(k[:]), v
}

// Last moves the cursor to the last item in the bucket and returns its key and value.
// If the bucket is empty then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Last() (key string, value []byte) {
	k, v := c.cursor.Last()
	return string(k[:]), v
}

// Next moves the cursor to the next item in the bucket and returns its key and value.
// If the cursor is at the end of the bucket then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Next() (key string, value []byte) {
	k, v := c.cursor.Next()
	return string(k[:]), v
}

// Prev moves the cursor to the previous item in the bucket and returns its key and value.
// If the cursor is at the beginning of the bucket then a nil key and value are returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Prev() (key string, value []byte) {
	k, v := c.cursor.Next()
	return string(k[:]), v
}

// Seek moves the cursor to a given key and returns it.
// If the key does not exist then the next key is used. If no keys
// follow, a nil key is returned.
// The returned key and value are only valid for the life of the transaction.
func (c *Cursor) Seek(seek string) (key string, value []byte) {
	k, v := c.cursor.Seek([]byte(seek))
	return string(k[:]), v
}

// Delete removes the current key/value under the cursor from the bucket.
// Delete fails if current key/value is a bucket or if the transaction is not writable.
func (c *Cursor) Delete() error {
	return c.cursor.Delete()
}
