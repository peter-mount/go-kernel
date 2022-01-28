package test

import (
	"github.com/peter-mount/go-kernel"
	"testing"
)

// testDeployService1 is an injectable Service used in tests
type testDeployService1 struct {
	service2 *testDeployService2 `kernel:"inject"`
}

func (t *testDeployService1) Name() string {
	return "testDeployService1"
}

// Calc will either panic or return 42
func (t *testDeployService1) Calc() int {
	return 6 * t.service2.Calc()
}

// testDeployService2 injected into testDeployService1
type testDeployService2 struct {
}

func (t *testDeployService2) Name() string {
	return "testDeployService2"
}

func (t *testDeployService2) Calc() int {
	return 7
}

// TestService_Inject tests the kernel injection mechanism.
// This is in a separate package as we do not want the Kernel to have direct access
// to unexported fields within these two test services
func TestService_Inject(t *testing.T) {

	s := &testDeployService1{}

	err := kernel.Launch(s)
	if err != nil {
		t.Fatal(err)
	}

	if s.service2 == nil {
		t.Fatal("No injected service")
	} else {
		result := s.Calc()
		if result != 42 {
			t.Errorf("Got %d expected 42", result)
		}
	}
}

// testDeployService3 tries to inject a non Service instance
type testDeployService3 struct {
	service2 *testDeployNonService3 `kernel:"inject"`
}

func (t *testDeployService3) Name() string {
	return "testDeployService3"
}

// testDeployNonService3 is a struct but not a service
type testDeployNonService3 struct {
}

// TestService_Inject tests the kenel injection mechanism.
// This is in a separate package as we do not want the Kernal to have direct access
// to unexported fields within these two test services
func TestService_InjectNonService(t *testing.T) {

	// Note originally this would fail but now any struct is deployable
	s := &testDeployService3{}

	err := kernel.Launch(s)
	if err != nil {
		t.Fatal(err)
	}
}

// testDeployService4 tries to inject a non-struct
type testDeployService4 struct {
	service2 *string `kernel:"inject"`
}

func (t *testDeployService4) Name() string {
	return "testDeployService4"
}

// TestService_Inject tests the kenel injection mechanism.
// This is in a separate package as we do not want the Kernal to have direct access
// to unexported fields within these two test services
func TestService_InjectNonStruct(t *testing.T) {

	s := &testDeployService4{}

	err := kernel.Launch(s)
	if err == nil {
		t.Fatal("No error returned")
	}
	if err.Error() != "Cannot deploy non-service" {
		t.Fatal("Unexpected error returned: " + err.Error())
	}
}

// testDeployService4 tries to inject a non-struct
type testDeployService5 struct {
	service2 string `kernel:"inject"`
}

func (t *testDeployService5) Name() string {
	return "testDeployService5"
}

// TestService_Inject tests the kenel injection mechanism.
// This is in a separate package as we do not want the Kernal to have direct access
// to unexported fields within these two test services
func TestService_InjectNonPointer(t *testing.T) {

	s := &testDeployService5{}

	err := kernel.Launch(s)
	if err == nil {
		t.Fatal("No error returned")
	}
	if err.Error() != "injection failed \"service2 string\" in github.com/peter-mount/go-kernel/test/testDeployService5: must be a pointer" {
		t.Fatal("Unexpected error returned: " + err.Error())
	}
}
