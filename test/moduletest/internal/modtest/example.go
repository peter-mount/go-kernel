package modtest

import (
	"github.com/peter-mount/go-kernel"
)

var Run bool

type ExampleService struct {
}

func (s *ExampleService) Run() error {
	Run = true
	return nil
}

func init() {
	kernel.Register(&ExampleService{})
}
