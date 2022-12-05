// Package kernel is a simple microkernel that allows for Service's to be deployed within
// an application.
//
// It manages the complete lifecycle of the application with muliple stages each
// called in sequence: Init, PostInit, Start & Run. Once the kernel gets to the
// Start phase then any Error will cause the Stop phase to be invoked to allow
// any Started service to cleanup.
//
// For most simple applications you can simply use kernel.Launch( s ) where s is
// an uninitiated service and it will create a Kernel add that service and run it.
//
// For more complex applications which need multiple unrelated services deployed
// then it can do by calling NewKernel() to create a new kernel, add each one via
// AddService() and then call Run() - this is what Launch() does internally.
//
// A Service is simply an Object implementing the Service interface and one or more
// of the various lifecycle interfaces.
//
// If a service has injectionPoints then it should implement Init() and call AddService
// to add them - the kernel will handle the rest.
package kernel

import (
	"errors"
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// Service to be deployed within the Kernel
type Service interface {
}

// NamedService is the original Service where the Name() function returns the unique name.
// This is now optional, a service does not require Name() anymore.
type NamedService interface {
	// Name returns the unique name of this service
	Name() string
}

// InitialisableService a Service that expects to be called in the Init lifecycle phase
type InitialisableService interface {
	// Init initialises a Service when it's added to the Kernel
	Init(*Kernel) error
}

// PostInitialisableService a Service that expects to be called in the PostInit lifecycle phase
type PostInitialisableService interface {
	// PostInit initialises a Service when it's added to the Kernel
	PostInit() error
}

// StartableService a Service that expects to be called in the Start lifecycle phase
type StartableService interface {
	// Start called when the Kernel starts but before services Run
	Start() error
}

// StoppableService a Service that expects to be called when the kernel shutsdown if it's in the
// Start or Run lifecycle phases
type StoppableService interface {
	Stop()
}

// RunnableService a Service that is expected to run in the Run lifecycle phase
type RunnableService interface {
	// Run executes the service
	Run() error
}

// Kernel is the core container for deployed services
type Kernel struct {
	services     util.List          // The deployed services
	stopList     util.List          // The services that are running & need to be shut down
	dependencies util.Set           // Used to prevent circular dependencies
	index        map[string]Service // Map of services by name
	readOnly     bool               // mark the kernel as read only
}

// Launch is a convenience method to launch a single service.
// This does the boiler plate work and requires the single service adds any
// injectionPoints within it's Init() method, if any
func Launch(services ...Service) error {

	// Add the supplied services in sequence. This creates the dependency graph
	if err := instance.DependsOn(services...); err != nil {
		return err
	}

	if instance.services.IsEmpty() {
		return errors.New("kernel is empty")
	}

	// From this point nothing else can be added to the Kernel
	instance.readOnly = true
	// When kernel exits, then reset it
	defer resetKernel()

	flag.Parse()

	// PostInit services
	if err := instance.postInit(); err != nil {
		return err
	}

	// Listen to signals & close the db before exiting
	// SIGINT for ^C, SIGTERM for docker stopping the container
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		log.Println("Signal", sig)

		instance.stop()

		log.Println("Application terminated")

		os.Exit(0)
	}()

	// At this point stop all started services on failure or exit
	defer instance.stop()

	// Start services
	if err := instance.start(); err != nil {
		return err
	}

	// Run services
	return instance.run()
}

// DependsOn just adds injectionPoints on other services, it does not return the resolved Service's.
// This is short of _,err:=instance.AddService() for each dependency.
func (k *Kernel) DependsOn(services ...Service) error {
	// Add the supplied services in sequence. This creates the dependency graph
	for _, s := range services {
		if _, err := instance.AddService(s); err != nil {
			return err
		}
	}
	return nil
}

func assertInstanceAmendable() error {
	if instance.readOnly {
		return errors.New("kernel is read only")
	}
	return nil
}

func getServiceName(t reflect.Type) string {
	return t.PkgPath() + "|" + t.Name()
}

// AddService adds a service to the kernel
func (k *Kernel) AddService(s Service) (Service, error) {
	// Generate the service name either via NamedService or reflection
	var name string
	if ns, ok := s.(NamedService); ok {
		name = ns.Name()
	} else {
		name = getServiceName(reflect.ValueOf(s).Elem().Type())
	}

	return k.addService(name, s, false)
}

func (k *Kernel) addService(name string, s Service, api bool) (Service, error) {
	if err := assertInstanceAmendable(); err != nil {
		return nil, err
	}

	// Prevent circular injectionPoints
	if instance.dependencies.Contains(name) {
		//if _, exists := instance.injectionPoints[s.Name()]; exists {
		return nil, fmt.Errorf("Circular dependency %s", name)
	}

	// Check we don't already have it
	if service, exists := instance.index[name]; exists {
		return service, nil
	}

	// At this point we must have a valid service
	if !api && reflect.ValueOf(s).Elem().Type().Kind() != reflect.Struct {
		return nil, errors.New("Cannot deploy non-service")
	}

	// This will prevent circular injectionPoints by using this map
	// to keep track of what's currently being deployed
	instance.dependencies.Add(name)
	defer instance.dependencies.Remove(name)

	// inject injectionPoints using struct field tags
	if err := instance.inject(s); err != nil {
		return nil, err
	}

	// Init the service, it can add injectionPoints here
	if is, ok := s.(InitialisableService); ok {
		if err := is.Init(k); err != nil {
			return nil, err
		}
	}

	// Finally, add the service to the end of the startup list
	instance.services.Add(s)
	instance.index[name] = s

	return s, nil
}

func (k *Kernel) postInit() error {
	return instance.services.ForEachFailFast(func(s interface{}) error {
		if pi, ok := s.(PostInitialisableService); ok {
			if err := pi.PostInit(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (k *Kernel) start() error {
	return instance.services.ForEachFailFast(func(s interface{}) error {
		// Start the service
		if ss, ok := s.(StartableService); ok {
			if err := (ss).Start(); err != nil {
				return err
			}
		}

		// Add to stop list if necessary
		if ss, ok := s.(StoppableService); ok {
			instance.stopList.Add(ss)
		}
		return nil
	})
}

func (k *Kernel) stop() {
	instance.stopList.ReverseIterator().ForEach(func(i interface{}) {
		(i).(StoppableService).Stop()
	})
}

func (k *Kernel) run() error {
	return instance.services.ForEachFailFast(func(s interface{}) error {
		if rs, ok := s.(RunnableService); ok {
			if err := rs.Run(); err != nil {
				return err
			}
		}
		return nil
	})
}
