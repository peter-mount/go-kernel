// Kernel is a simple microkernel that allows for Service's to be deployed within
// an application.
//
// It manages the complete lifecycle of the application with muliple stages each
// called in sequence: Init, PostInit, Start & Run. Once the kernel gets to the
// Start phase then any error will cause the Stop phase to be invoked to allow
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
// If a service has dependencies then it should implement Init() and call AddService
// to add them - the kernel will handle the rest.
package kernel

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Service to be deployed within the Kernel
type Service interface {
	// Name returns the unique name of this service
	Name() string
}

// A Service that is expected to run in the Run lifecycle phase
type RunnableService interface {
	// Run executes the service
	Run() error
}

// A Service that expects to be called in the Init lifecycle phase
type InitialisableService interface {
	// Init initialises a Service when it's added to the Kernel
	Init(*Kernel) error
}

// A Service that expects to be called in the PostInit lifecycle phase
type PostInitialisableService interface {
	// Init initialises a Service when it's added to the Kernel
	PostInit() error
}

// A Service that expects to be called in the Start lifecycle phase
type StartableService interface {
	// Start called when the Kernel starts but before services Run
	Start() error
}

// A Service that expects to be called when the kernel shutsdown if it's in the
// Start or Run lifecycle phases
type StoppableService interface {
	Stop()
}

// Kernel is the core container for deployed services
type Kernel struct {
	// The deployed services
	services []Service
	// The services that are running & need to be shut down
	stopList []StoppableService
	// Used to prevent circular dependencies
	dependencies map[string]interface{}
	// mark the kernel as read only
	readOnly bool
}

// Launch is a convenience method to launch a single service.
// This does the boiler plate work and requires the single service adds any
// dependencies within it's Init() method, if any
func Launch(services ...Service) error {
	k := &Kernel{}
	k.dependencies = make(map[string]interface{})

	for _, s := range services {
		if _, err := k.AddService(s); err != nil {
			return err
		}
	}

	k.readOnly = true

	flag.Parse()

	if err := k.postinit(); err != nil {
		return err
	}

	// Listen to signals & close the db before exiting
	// SIGINT for ^C, SIGTERM for docker stopping the container
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println("Signal", sig)

		k.stop()

		log.Println("Application terminated")

		os.Exit(0)
	}()

	defer k.stop()

	if err := k.start(); err != nil {
		return err
	}

	return k.run()
}

// AddService adds a service to the kernel
func (k *Kernel) AddService(s Service) (Service, error) {
	if k.readOnly {
		return nil, fmt.Errorf("Cannot add %s as Kernel is read only", s.Name())
	}

	// Prevent circular dependencies
	if _, exists := k.dependencies[s.Name()]; exists {
		return nil, fmt.Errorf("Circular dependency %s", s.Name())
	}

	// Check we don't already have it
	for _, e := range k.services {
		if e.Name() == s.Name() {
			return e, nil
		}
	}

	// This will prevent circular dependencies
	k.dependencies[s.Name()] = nil
	defer delete(k.dependencies, s.Name())

	// Init the service, it can add dependencies here
	if is, ok := s.(InitialisableService); ok {
		if err := is.Init(k); err != nil {
			return nil, err
		}
	}

	// Finally add the service to the end of the startup list
	k.services = append(k.services, s)

	return s, nil
}

func (k *Kernel) postinit() error {
	for _, s := range k.services {
		// Start the service
		if pi, ok := s.(PostInitialisableService); ok {
			if err := pi.PostInit(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k *Kernel) start() error {
	for _, s := range k.services {
		// Start the service
		if ss, ok := s.(StartableService); ok {
			if err := (ss).Start(); err != nil {
				return err
			}
		}

		// Add to stop list if necessary
		if ss, ok := s.(StoppableService); ok {
			k.stopList = append(k.stopList, ss)
		}
	}

	return nil
}

func (k *Kernel) stop() {
	for i := len(k.stopList) - 1; i >= 0; i-- {
		k.stopList[i].Stop()
	}
}

func (k *Kernel) run() error {
	for _, s := range k.services {
		if rs, ok := s.(RunnableService); ok {
			if err := rs.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
