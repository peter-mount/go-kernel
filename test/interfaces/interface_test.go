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

type API interface {
	Get() int
}

type Service1 struct{}

func (s *Service1) Get() int { return 1 }

type Service2 struct{}

func (s *Service2) Get() int { return 2 }

type Service3 struct {
	api API `kernel:"inject"`
}

// gets set to true when the service is run
var service3Run bool

func (s *Service3) Run() error {
	result := s.api.Get()
	if result != 1 {
		return fmt.Errorf("expected service 1 got %d", result)
	}

	service3Run = true
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

func TestInterfaceLookup(t *testing.T) {

	// Register Service1 against API
	kernel.RegisterAPI((*API)(nil), &Service1{})

	// Try to register service 2 as API.
	// This would normally panic as it's already registered, but we want this hence
	// we use the helper function to capture the panic
	registerService(t, &Service2{})

	// Now run Service3 which depends on API. It will fail if Service2 is present and not 1
	err := kernel.Launch(&Service3{})
	if err != nil {
		t.Fatal(err)
	}

	if !service3Run {
		t.Fatal("Service3 did not run")
	}
}
