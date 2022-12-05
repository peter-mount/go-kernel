package moduletest

import (
	"github.com/peter-mount/go-kernel"
	"github.com/peter-mount/go-kernel/test/moduletest/internal/modtest"
	"testing"
)

// TestPackageAutodeploy simply launches an empty kernel.
// The package import for modtest should do the actual registering of the services
func TestPackageAutodeploy(t *testing.T) {

	// Launch an empty kernel. The import of modtest should have deployed the service
	err := kernel.Launch()
	if err != nil {
		// This will happen if no services were deployed
		t.Fatal(err)
	}

	// This should fail if the service did not run.
	// If the service did not even deploy then Launch() would have failed already
	if !modtest.Run {
		t.Fatal("service did not autodeploy")
	}
}
