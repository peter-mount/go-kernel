package bolt

import (
	"encoding/json"
)

// GetJSON retrieves the value for a key in the bucket.
// Returns a nil value if the key does not exist or if the key is a nested bucket.
// The returned value is only valid for the life of the transaction.
func (b *Bucket) GetJSON(key string, val interface{}) bool {
	v := b.Get(key)
	if v == nil {
		return false
	}

	err := json.Unmarshal(v, val)
	return err == nil
}

// PutJSON sets the value for a key in the bucket.
// If the key exist then its previous value will be overwritten.
// Supplied value must remain valid for the life of the transaction.
// Returns an error if the bucket was created from a read-only transaction, if the key is blank, if the key is too large, or if the value is too large.
func (b *Bucket) PutJSON(key string, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return b.bucket.Put([]byte(key), v)
}
