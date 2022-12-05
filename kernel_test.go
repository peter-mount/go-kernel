package kernel

import "testing"

type testService struct {
	start bool
	run   bool
	stop  bool
	seq   int
}

func (s *testService) Name() string {
	return "testService"
}

func (s *testService) Start() error {
	s.start = true
	return nil
}

func (s *testService) Run() error {
	s.run = true
	return nil
}

func (s *testService) Stop() {
	s.stop = true
}

type testService2 struct {
	t          *testing.T
	dependency *testService
	start      bool
	run        bool
	stop       bool
}

func (s *testService2) Name() string {
	return "testService2"
}

func (s *testService2) Init(k *Kernel) error {
	service, err := k.AddService(&testService{})
	if err != nil {
		return err
	}
	s.dependency = (service).(*testService)

	return nil
}

func (s *testService2) Start() error {
	if !s.dependency.start {
		s.t.Errorf("service 1 not started")
	}
	s.start = true
	return nil
}

func (s *testService2) Stop() {
	if s.dependency.stop {
		s.t.Errorf("service 1 stopped before service 2")
	}
	s.stop = true
}

func TestKernel_Launch(t *testing.T) {
	s := &testService{}

	err := Launch(s)
	if err != nil {
		t.Errorf("Launch failed: %v", err)
	}

	if !s.start {
		t.Errorf("Test service did not start")
	}

	if !s.run {
		t.Errorf("Test service did not run")
	}

	if !s.stop {
		t.Errorf("Test service did not stop")
	}
}

func TestKernel_LaunchOrder(t *testing.T) {
	s := &testService2{t: t}

	err := Launch(
		&testService{seq: 1}, // seq1 so we know if we pick up this instance or a new one
		s,
	)
	if err != nil {
		t.Errorf("Launch failed: %v", err)
	}

	if s.dependency.seq != 1 {
		s.t.Errorf("service 1 was not original one deployed")
	}
}

// TestEmptyKernel should fail if launching with no services does _not_ return an error
func TestEmptyKernel(t *testing.T) {

	err := Launch()
	if err == nil {
		t.Fatal("No error returned")
	}

}
