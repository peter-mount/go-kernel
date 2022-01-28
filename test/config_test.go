package test

import (
	"github.com/peter-mount/go-kernel"
	"testing"
)

type configService struct {
	conf *config1 `kernel:"config,test1"`
}

type config1 struct {
	Name string `yaml:"name"`
}

type configService2 struct {
	cs1  *configService `kernel:"inject"`
	conf *config1       `kernel:"config,test1"`
}

func TestConfig(t *testing.T) {
	to := &configService{}

	err := kernel.Launch(to)
	if err != nil {
		t.Fatal(err)
	}

	if to.conf == nil {
		t.Fatal("No config injected")
	}

	if to.conf.Name != "test" {
		t.Errorf("Name not read, got %q expected \"test\"", to.conf.Name)
	}
}

func TestConfig_Multi(t *testing.T) {
	cs2 := &configService2{}

	err := kernel.Launch(cs2)
	if err != nil {
		t.Fatal(err)
	}

	// Test injected instance in cs2

	if cs2.conf == nil {
		t.Fatal("No config injected")
	}

	if cs2.conf.Name != "test" {
		t.Errorf("Name not read, got %q expected \"test\"", cs2.conf.Name)
	}

	// Now check instance under cs1

	if cs2.cs1.conf == nil {
		t.Fatal("No dependency injected")
	}

	if cs2.cs1.conf.Name != "test" {
		t.Errorf("Name not read, got %q expected \"test\"", cs2.cs1.conf.Name)
	}

	// Now check instance injected in both is the same one

	if cs2.conf != cs2.cs1.conf {
		t.Errorf("Injected instance should be same but isn't")
	}
}
