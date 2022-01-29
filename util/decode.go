package util

import (
	"math/big"
	"strconv"
)

// ForEachInterface will invoke a function for every entry in v if v is a slice
func ForEachInterface(v interface{}, f func(interface{}) error) error {
	if a, ok := v.([]interface{}); ok {
		for _, e := range a {
			if err := f(e); err != nil {
				return err
			}
		}
	}
	return nil
}

// IfMap will invoke a function if v is a map
func IfMap(v interface{}, f func(map[interface{}]interface{}) error) error {
	if m, ok := v.(map[interface{}]interface{}); ok {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

// IfMapEntry invokes a function if a map contains an entry
func IfMapEntry(m map[interface{}]interface{}, n interface{}, f func(interface{}) error) error {
	if v, ok := m[n]; ok {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func DecodeString(v interface{}, def string) string {
	var r string

	if v != nil {
		if s, ok := v.(string); ok {
			r = s
		} else if i, ok := v.(int); ok {
			r = strconv.Itoa(i)
		}
	}

	if r == "" {
		r = def
	}

	return r
}

func DecodeInt(v interface{}, def int) (int, bool) {
	r := def

	if v != nil {
		if s, ok := v.(string); ok {
			n := new(big.Int)
			n.SetString(s, 0)
			r = int(n.Int64())
			/*
			   i, err := strconv.Atoi(s)
			   if err != nil {
			     return 0, false
			   }
			   r = i
			*/
		} else if i, ok := v.(int); ok {
			r = i
		}
	}

	return r, true
}

func DecodeBool(v interface{}) (bool, bool) {

	if v != nil {
		if s, ok := v.(string); ok {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return false, false
			}
			return b, true
		} else if i, ok := v.(int); ok {
			return i != 0, true
		}

		b, ok := v.(bool)
		return b, ok
	}

	return false, true
}

func IfMapEntryString(m map[interface{}]interface{}, n string) string {
	var s string
	_ = IfMapEntry(m, n, func(i interface{}) error {
		s = DecodeString(i, "")
		return nil
	})
	return s
}

func IfMapEntryBool(m map[interface{}]interface{}, n string) bool {
	var s bool
	_ = IfMapEntry(m, n, func(i interface{}) error {
		s, _ = DecodeBool(i)
		return nil
	})
	return s
}

func IfMapEntryInt(m map[interface{}]interface{}, n string) int {
	var s int
	_ = IfMapEntry(m, n, func(i interface{}) error {
		s, _ = DecodeInt(i, 0)
		return nil
	})
	return s
}
