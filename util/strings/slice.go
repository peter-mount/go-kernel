package strings

import (
	"io"
	"sort"
	"strings"
)

type StringSlice []string

// Of returns a StringSlice from a slice of strings
func Of(s ...string) StringSlice {
	return s
}

// ForEach invokes a StringHandler for each entry in a StringSlice. First one to return an error terminates the loop.
func (s StringSlice) ForEach(f StringHandler) error {
	for _, b := range s {
		err := f(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// Write will write a StringSlice to a writer with each entry being a single line.
func (s StringSlice) Write(w io.Writer) error {
	_, err := w.Write([]byte(strings.Join(s, "\n")))
	return err
}

// IsEmpty returns true if a StringSlice has no entries
func (s StringSlice) IsEmpty() bool {
	return len(s) == 0
}

// Join returns the content of a StringSlice with the specified separator between each entry
func (s StringSlice) Join(sep string) string {
	return strings.Join(s, sep)
}

// Join2 is similar to Join except it adds the prefix & suffix to the final result
func (s StringSlice) Join2(prefix, suffix, sep string) string {
	return prefix + s.Join(sep) + suffix
}

// Sort sorts the string slice using a case-insensitive comparator.
func (s StringSlice) Sort() StringSlice {
	sort.SliceStable(s, func(i, j int) bool {
		return strings.ToLower(s[i]) < strings.ToLower(s[j])
	})
	return s
}

type StringSliceHandler func(StringSlice) (StringSlice, error)

func (a StringSliceHandler) Do(s StringSlice) (StringSlice, error) {
	if a != nil {
		return a(s)
	}
	return nil, nil
}

func (a StringSliceHandler) Then(b StringSliceHandler) StringSliceHandler {
	return func(s StringSlice) (StringSlice, error) {
		s, err := a(s)
		if err != nil {
			return nil, err
		}
		return b(s)
	}
}
