package test

import (
	"flag"
	"github.com/peter-mount/go-kernel"
	"testing"
)

type testFlagService struct {
	tBool  *bool    `kernel:"flag,bool"`
	tInt   *int     `kernel:"flag,intFlag"`
	tInt2  *int64   `kernel:"flag,longFlag,,42"`
	tFloat *float64 `kernel:"flag,float,A 64 bit float,3.1415926"`
}

func (t *testFlagService) Name() string {
	return "testFlagService"
}

func TestFlag_Inject(t *testing.T) {

	s := &testFlagService{}

	err := kernel.Launch(s)
	if err != nil {
		t.Fatal(err)
	}

	if s.tBool == nil {
		t.Errorf("bool flag missing")
	}

	flag.PrintDefaults()
}
