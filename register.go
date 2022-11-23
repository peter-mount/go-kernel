package kernel

import (
	"errors"
	"fmt"
	"github.com/peter-mount/go-kernel/util"
	"reflect"
)

// Our singleton instance of the kernel
var instance *Kernel

func init() {
	// Ensure we have an initialised kernel.
	// Note this is a separate function as other code cannot access init()
	resetKernel()
}

func resetKernel() {
	instance = &Kernel{
		dependencies: util.NewSyncSet(),
		services:     util.NewList(),
		stopList:     util.NewList(),
		index:        make(map[string]Service),
	}
}

// Register will add the specified services to the kernel.
// If the kernel has been started then this will panic.
//
// This is normally used within a packages' init() function to automatically deploy services.
//
func Register(services ...Service) {
	err := assertInstanceAmendable()
	if err == nil {
		err = instance.DependsOn(services...)
	}

	if err != nil {
		panic(err)
	}
}

// RegisterAPI registers an API
func RegisterAPI(api interface{}, service Service) {
	err := assertInstanceAmendable()
	if err != nil {
		panic(err)
	}

	// api must be an interface
	kt := reflect.TypeOf(api).Elem()
	if kt.Kind() != reflect.Interface {
		panic(errors.New("cannot register non-interface"))
	}

	name := getServiceName(kt)

	resolvedService, err := instance.addService(name, service, true)
	if err != nil {
		panic(err)
	}

	if resolvedService != service {
		panic(fmt.Errorf("service %s already registered", name))
	}
}
