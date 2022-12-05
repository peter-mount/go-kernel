package walk

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// Directory names.
var dirNames = []string{"dir1", "dir2/inaccessible", "dir2", "dir3", "dir4"}

const (
	disownEntry   = 2 // Entry in dirNames to disown so it, and it's content cannot be accessed
	symLinkTarget = 3 // Dir to be target of symlink
	symLinkParent = 4 // Dir containing the symlink
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
		//fmt.Printf("Creating dir %s\n", dn)
		err = os.MkdirAll(dn, os.ModePerm)

		if err == nil && i == disownEntry {
			// Set permissions so we cannot enter this directory
			//fmt.Printf("Disowning %s\n", dn)
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
			//fmt.Printf("%s %#o\n", path, info.Mode())
			return nil
		}).
		Walk(rootDir)

	if err != nil {
		t.Fatalf("Path walker test failed: %v", err)
	}
}

func TestLinkFollow(t *testing.T) {
	rootDir, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	linkName := "link"
	symlink := path.Join(rootDir, dirNames[symLinkParent], linkName)
	target := path.Join("..", dirNames[symLinkTarget])
	//fmt.Printf("link %s -> %s", symlink, target)
	err = os.Symlink(target, symlink)
	if err != nil {
		t.Fatal(err)
	}

	testName := "test"
	err = ioutil.WriteFile(path.Join(rootDir, dirNames[symLinkTarget], testName), []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	count := 0

	err = NewPathWalker().
		Then(func(filePath string, info os.FileInfo) error {
			if path.Base(filePath) == testName {
				count = count + 1
				//fmt.Printf("Match %d %q %s\n", count, linkName, filePath)
			}
			return nil
		}).
		IsFile().
		FollowSymlinks().
		Walk(rootDir)

	if err != nil {
		t.Fatalf("Path walker test failed: %v", err)
	}

	if count != 2 {
		t.Fatalf("Did not visit %q 2 times, visited it %d", linkName, count)
	}
}
