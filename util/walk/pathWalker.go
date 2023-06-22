package walk

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// PathWalker performs a function against a file being walked
type PathWalker func(path string, info os.FileInfo) error

// PathPredicate performs a test against a file being walked
type PathPredicate func(path string, info os.FileInfo) bool

// Not negates a PathPredicate
func (a PathPredicate) Not() PathPredicate {
	return func(path string, info os.FileInfo) bool {
		return !a(path, info)
	}
}

// NewPathWalker creates a new PathWalker
func NewPathWalker() PathWalker {
	return nil
}

// Do will call a PathWalker. If the walker is null then null is returned.
func (a PathWalker) Do(path string, info os.FileInfo) error {
	if a != nil {
		return a(path, info)
	}
	return nil
}

// Then performs the current PathWalker then another PathWalker
func (a PathWalker) Then(b PathWalker) PathWalker {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(path string, info os.FileInfo) error {
		err := a(path, info)
		if err != nil {
			return err
		}
		return b(path, info)
	}
}

var (
	// Used to stop the chain when predicate is false
	predicateFail = errors.New("predicate fail")
)

// If will test a path with a PathPredicate and allow the path to be processed only if the PathPredicate passes
func (a PathWalker) If(p PathPredicate) PathWalker {
	if a == nil {
		return nil
	}
	return func(path string, info os.FileInfo) error {
		if p(path, info) {
			return a(path, info)
		}
		return predicateFail
	}
}

// IfNot is the same as If except only allows processing to continue if the PathPredicate returns false
func (a PathWalker) IfNot(p PathPredicate) PathWalker {
	return a.If(p.Not())
}

// IsDir allows processing only if the current path is a Directory
func (a PathWalker) IsDir() PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(_ string, info os.FileInfo) bool {
		return info.IsDir()
	})
}

// IsFile allows processing only if the current path is a File
func (a PathWalker) IsFile() PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(_ string, info os.FileInfo) bool {
		return !info.IsDir()
	})
}

// PathContains allows processing if the path contains the provided string
func (a PathWalker) PathContains(s string) PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(path string, _ os.FileInfo) bool {
		return strings.Contains(path, s)
	})
}

// PathNotContain allows processing if the path does not contain the provided string
func (a PathWalker) PathNotContain(s string) PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(path string, _ os.FileInfo) bool {
		return !strings.Contains(path, s)
	})
}

// PathHasSuffix allows processing if the path has the provided suffix
func (a PathWalker) PathHasSuffix(s string) PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(path string, _ os.FileInfo) bool {
		return strings.HasSuffix(path, s)
	})
}

// PathHasNotSuffix allows processing if the path has not got the provided suffix
func (a PathWalker) PathHasNotSuffix(s string) PathWalker {
	if a == nil {
		return nil
	}
	return a.If(func(path string, _ os.FileInfo) bool {
		return !strings.HasSuffix(path, s)
	})
}

// FollowSymlinks will cause the walker to follow a symlink, unlike filepath.Walk() with refuses to do do.
func (a PathWalker) FollowSymlinks() PathWalker {
	return func(path string, info os.FileInfo) error {
		if (info.Mode() & os.ModeSymlink) == os.ModeSymlink {
			link, err := os.Readlink(path)
			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}
			link = filepath.Join(filepath.Dir(path), link)
			if err != nil {
				return err
			}
			return a.Walk(link)
		}

		return a.Do(path, info)
	}
}

// Walk performs the actual walk against the built PathWalker
func (a PathWalker) Walk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Any error walking to the file or directory will abort the walk immediately
		if err != nil {
			// go-kernel#1 If we are trying to enter a directory, and it's a permission error then skip the directory
			if info != nil && info.IsDir() && os.IsPermission(err) {
				return filepath.SkipDir
			}
			return err
		}

		// Enter our PathWalker chain
		err = a.Do(path, info)

		// Don't pass this to the Walker as it just means we have stopped processing
		if err == predicateFail {
			return nil
		}
		return err
	})
}
