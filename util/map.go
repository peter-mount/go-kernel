package util

// Map is a go equivalent to the Java Map interface
type Map interface {
	Collection
	// ContainsKey returns true if the map contains a key
	ContainsKey(k string) bool
	// Get returns a value from the map or nil if absent
	Get(k string) interface{}
	// Get returns a value from the map or nil if absent
	Get2(k string) (interface{}, bool)
	// Put sets a value into the map returning the existing entry or nil of absent
	Put(k string, v interface{}) interface{}
	// PutIfAbsent sets a value in the map if the key does not already exist.
	PutIfAbsent(k string, v interface{}) interface{}
	// Remove removes a key, returning the value present
	Remove(k string) interface{}
	// RemoveIfEquals removes a key if it's the same value.
	// Similar to java.util.Map.Remove(k,v)
	RemoveIfEquals(k string, v interface{}) bool
	// ReplaceIfEquals replaces a entry for a key if it's current value is the one specified
	// Similar to java.util.Map.replace(k,old,new)
	ReplaceIfEquals(k string, oldValue interface{}, newValue interface{}) bool
	// Replace replaces a entry for a key if it currently has a value.
	// Similar to java.util.Map.replace(k,v)
	Replace(k string, v interface{}) bool
	// ComputeIfAbsent a value if the specified key is not already associated with a value (or is mapped to nil).
	// The returned value will be the value in the map (or the new one inserted)
	ComputeIfAbsent(k string, f func(string) interface{}) interface{}
	// ComputeIfPresent will call a function if a value exists in the map.
	// The returned value from the function will be the new value for the key in the map
	// & will be the returned value from this call.
	ComputeIfPresent(k string, f func(string, interface{}) interface{}) interface{}
	// Compute will call a function with the current value in the map (will be nil if absent)
	// The returned value from the function will be the new value for the key in the map
	// & will be the returned value from this call.
	Compute(k string, f func(string, interface{}) interface{}) interface{}
	// Merge will ether use the supplied value as the new value if an entry does not exist (or is nil)
	// otherwise will use the supplied function to return the value to be set.
	// For example, in Java: map.merge(key, msg, String::concat) would set to msg if key was absent,
	// or the existing value with msg concatenated.
	Merge(k string, v interface{}, f func(interface{}, interface{}) interface{}) interface{}
	// ForEach will call a function for each key,value pair in the map.
	// WARNING: For synchronized maps the function will be called from within the lock so the supplied
	// function cannot modify the map. See ForEachAsync
	ForEach(f func(string, interface{}))
	// ForEachAsync will call a function for each key, value pair in the map.
	// Unlike ForEach this makes a copy of the key's internally then traverses it, calling the function
	// outside of any lock as long as the value exists at the time the function would be called.
	// This allows for the map to be modified whilst the ForEachAsync is iterating.
	// WARNING: This is expensive as it makes a copy of the keys into a slice
	ForEachAsync(f func(string, interface{}))
	// Keys returns a slice of all keys currently in the map. The entries in this slice is only guaranteed to be correct
	// for the moment the function was called. The map can be changed afterwards
	Keys() []string
	// ExecIfPresent will execute a function if a value exists in the map
	ExecIfPresent(k string, f func(interface{}))
}
