package test

import (
	"github.com/peter-mount/go-kernel"
	"testing"
)

const (
	// Values in config.yaml and config_include.yaml
	test1Name = "test"
	test2Name = "included test"
)

type configService struct {
	conf         *config1 `kernel:"config,test1"`
	includedConf *config1 `kernel:"config,test2"`
}

type config1 struct {
	Name string `yaml:"name"`
}

func (c *config1) test(t *testing.T, expected string) {
	if c == nil {
		t.Fatal("No dependency injected")
	}

	if c.Name != expected {
		t.Errorf("Name not read, got %q expected %q", c.Name, expected)
	}
}

type configService2 struct {
	cs1  *configService `kernel:"inject"`
	conf *config1       `kernel:"config,test1"`
}

func TestConfig_Multi(t *testing.T) {
	cs2 := &configService2{}

	err := kernel.Launch(cs2)
	if err != nil {
		t.Fatal(err)
	}

	// Test injected instance in cs2
	cs2.conf.test(t, test1Name)

	// Now check instance under cs1
	cs2.cs1.conf.test(t, test1Name)

	// Now check instance injected in both is the same one

	if cs2.conf != cs2.cs1.conf {
		t.Errorf("Injected instance should be same but isn't")
	}

	// Test included config
	cs2.cs1.includedConf.test(t, test2Name)
}
