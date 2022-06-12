package kernel

import "github.com/peter-mount/go-kernel/util"

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
