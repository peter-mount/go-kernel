package walk

import (
	"fmt"
	"os"
	"path"
	"testing"
)

// Directory names.
var dirNames = []string{"dir1", "dir2/inaccessible", "dir2", "dir3", "dir4"}

const (
	disownEntry = 2 // Entry in dirNames to disown so it, and it's content cannot be accessed
)

// setup creates a test directory testdir with some test files/directories.
// It returns the directory created or an error
func setup() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dirName := path.Join(home, "go-kernel-testdir")

	err = os.RemoveAll(dirName)
	if err != nil && !os.IsNotExist(err) && !os.IsPermission(err) {
		return "", err
	}

	for i, n := range dirNames {
		dn := path.Join(dirName, n)
		fmt.Printf("Creating dir %s\n", dn)
		err = os.MkdirAll(dn, os.ModePerm)

		if err == nil && i == disownEntry {
			// Set permissions so we cannot enter this directory
			fmt.Printf("Disowning %s\n", dn)
			err = os.Chmod(dn, 0)
		}
		if err != nil && !os.IsPermission(err) {
			return "", err
		}
	}

	return dirName, nil
}

/*
 * TestIssue1 tests https://github.com/peter-mount/go-kernel/issues/1
 * utilising some code from https://go.dev/play/p/BN8rtknwFvA linked in that issue.
 *
 * To run this test you need to do the following prep:
 * mkdir /tmp/testdir
 * mkdir /tmp/testdir/rootOwned
 * chmod 700 /tmp/testdir/rootOwned
 *
 * This creates an entry with a subdirectory owned by root and not readable by the test user
 */
func TestIssue1(t *testing.T) {
	rootDir, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	err = NewPathWalker().
		IsFile().
		Then(func(path string, info os.FileInfo) error {
			fmt.Printf("%s %#o\n", path, info.Mode())
			return nil
		}).
		Walk(rootDir)

	if err != nil {
		t.Fatalf("Path walker test failed: %v", err)
	}
}
