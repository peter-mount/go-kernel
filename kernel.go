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
	"github.com/peter-mount/go-kernel/util"
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

// A Service that is expected to run in the Run lifecycle phase
type RunnableService interface {
	// Run executes the service
	Run() error
}

// Kernel is the core container for deployed services
type Kernel struct {
	// The deployed services
	//services []Service
	services util.List
	// The services that are running & need to be shut down
	stopList util.List
	// Used to prevent circular dependencies
	dependencies util.Set
	// mark the kernel as read only
	readOnly bool
}

// Launch is a convenience method to launch a single service.
// This does the boiler plate work and requires the single service adds any
// dependencies within it's Init() method, if any
func Launch(services ...Service) error {
	k := &Kernel{
		dependencies: util.NewSyncSet(),
		services:     util.NewList(),
		stopList:     util.NewList(),
	}

	// Add the supplied services in sequence. This creates the dependency graph
	for _, s := range services {
		if _, err := k.AddService(s); err != nil {
			return err
		}
	}

	// From this point nothing else can be added to the Kernel
	k.readOnly = true

	flag.Parse()

	// PostInit services
	if err := k.postInit(); err != nil {
		return err
	}

	// Listen to signals & close the db before exiting
	// SIGINT for ^C, SIGTERM for docker stopping the container
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		log.Println("Signal", sig)

		k.stop()

		log.Println("Application terminated")

		os.Exit(0)
	}()

	// At this point stop all started services on failure or exit
	defer k.stop()

	// Start services
	if err := k.start(); err != nil {
		return err
	}

	// Run services
	return k.run()
}

// AddService adds a service to the kernel
func (k *Kernel) AddService(s Service) (Service, error) {
	if k.readOnly {
		return nil, fmt.Errorf("Cannot add %s as Kernel is read only", s.Name())
	}

	name := s.Name()

	// Prevent circular dependencies
	if k.dependencies.Contains(name) {
		//if _, exists := k.dependencies[s.Name()]; exists {
		return nil, fmt.Errorf("Circular dependency %s", name)
	}

	// Check we don't already have it
	if i := k.services.FindIndexOf(func(e interface{}) bool {
		return (e).(Service).Name() == name
	}); i > -1 {
		return (k.services.Get(i)).(Service), nil
	}

	// This will prevent circular dependencies by using this map
	// to keep track of what's currently being deployed
	k.dependencies.Add(name)
	defer k.dependencies.Remove(name)

	// Init the service, it can add dependencies here
	if is, ok := s.(InitialisableService); ok {
		if err := is.Init(k); err != nil {
			return nil, err
		}
	}

	// Finally add the service to the end of the startup list
	k.services.Add(s)

	return s, nil
}

func (k *Kernel) postInit() error {
	return k.services.ForEachFailFast(func(s interface{}) error {
		if pi, ok := s.(PostInitialisableService); ok {
			if err := pi.PostInit(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (k *Kernel) start() error {
	return k.services.ForEachFailFast(func(s interface{}) error {
		// Start the service
		if ss, ok := s.(StartableService); ok {
			if err := (ss).Start(); err != nil {
				return err
			}
		}

		// Add to stop list if necessary
		if ss, ok := s.(StoppableService); ok {
			k.stopList.Add(ss)
		}
		return nil
	})
}

func (k *Kernel) stop() {
	k.stopList.ReverseIterator().ForEach(func(i interface{}) {
		(i).(StoppableService).Stop()
	})
}

func (k *Kernel) run() error {
	return k.services.ForEachFailFast(func(s interface{}) error {
		if rs, ok := s.(RunnableService); ok {
			if err := rs.Run(); err != nil {
				return err
			}
		}
		return nil
	})
}
