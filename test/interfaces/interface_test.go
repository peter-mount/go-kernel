// interfaces tests the new functionality where a service can be
// registered against an interface. Doing this allows services to protect their
// inner workings from callers
package interfaces

import (
	"fmt"
	"github.com/peter-mount/go-kernel"
	"strings"
	"testing"
)

// API interface implemented by Service1 & Service2.
//
// It implements a single method which returns the ID of the underlying service,
// 1 for Service1 & 2 for Service2.
//
type API interface {
	Get() int
}

// Service1 implements API and returns 1
type Service1 struct{}

func (s *Service1) Get() int { return 1 }

// Service2 implements API and returns 2
type Service2 struct{}

func (s *Service2) Get() int { return 2 }

// API3 interface implemented by Service3.
//
// Its sole method returns true if Service3 has run with no errors.
type API3 interface {
	HasRun() bool
}

// Service3 implements API3, has a dependency on the service implementing the
// API interface.
type Service3 struct {
	Api    API `kernel:"inject"`
	hasRun bool
}

// HasRun returns true only if Run() has run to completion.
// This allows us to check that Service3 has been deployed.
func (s *Service3) HasRun() bool {
	return s.hasRun
}

// Run checks that the injected service is the correct one (Service1), then sets
// hasRun to true to indicate success
func (s *Service3) Run() error {
	result := s.Api.Get()
	if result != 1 {
		return fmt.Errorf("expected service 1 got %d", result)
	}

	s.hasRun = true
	return nil
}

// Wrapper around RegisterAPI but traps any panic's which are required by TestInterfaceLookup
// as we want to ignore a service being registered twice
func registerService(t *testing.T, s kernel.Service) {
	// The second RegisterAPI call should fail with an already registered error
	// but here we want that so don't fail if it does
	defer func() {
		if err := recover(); err != nil {
			err2 := err.(error)
			msg := err2.Error()
			if !strings.HasSuffix(msg, "already registered") {
				t.Fatal(err)
			}
		}
	}()
	kernel.RegisterAPI((*API)(nil), s)
}

// TestInterfaceLookup is the real test.
//
// Here it registers Service1 against the API interface.
//
// It then tries to register Service2 against the same API. This must fail with an
// expected panic so fail the test if that panic does NOT occur.
//
// We then register Service3 with the API3 interface to test that multiple interfaces
// can be registered.
//
// Finally, we launch the kernel. this should run Service3 to completion.
// If it returns an error then it's a failure.
// e.g. Service2 was present not Service1 when registering the API.
//
// If Service3 did not run (i.e. did not deploy) then we fail the test.
func TestInterfaceLookup(t *testing.T) {

	// Register Service1 against API
	kernel.RegisterAPI((*API)(nil), &Service1{})

	// Try to register service 2 as API.
	// This would normally panic as it's already registered, but we want this hence
	// we use the helper function to capture the panic
	registerService(t, &Service2{})

	// Now register service3. This also tests that we are registering API & API3 as two
	// separate interfaces and not just registering under a common one
	service3 := &Service3{}
	kernel.RegisterAPI((*API3)(nil), service3)

	// Now run Service3 which depends on API. It will fail if Service2 is present and not 1
	err := kernel.Launch()
	if err != nil {
		t.Fatal(err)
	}

	// The test passes if the service did run.
	if !service3.HasRun() {
		t.Fatal("Service3 did not run")
	}
}
