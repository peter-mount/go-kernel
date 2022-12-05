# go-kernel - A service microkernel for golang applications.

go-kernel is a simple microkernel which can run multiple services within the same go runtime. The kernel provides both
lifecycle
support for those services and dependency management so that services start and stop in the correct order.

Each service can declare any dependencies with other services which are then automatically injected
at startup. This ensures that services are started in the correct order.
It also allows for code reuse without most of the required boilerplate code usually required.

Full documentation is available online as [html](https://area51.dev/go/go-kernel/) as well
as [pdf](https://area51.dev/static/book/go-kernel.pdf)
however a brief overview is listed below.

## Services

A service is simply a struct which exposes functions to other services or performs some action during the applications'
lifetime.

    package myapp
    
    // An example service
    type Example struct {
        // Example of injecting another service as a dependency
        Config *conf.Config `kernel:"inject"`
        // Example of injecting a common task worker queue provided by the kernel
        Worker task.Queue `kernel:"worker"`
        // Example if declaring an indirect dependency
        _ *PostCSS `kernel:"inject"`
        // Example of creating a command line parameter
        server *bool `kernel:"flag,s,Run hugo in server mode"`
    }
    
    // Kernel lifecycle, this gets called during the start phase of the application
    // This is optional, as are all of the life cycle functions.
    func (s *MyService) Start() error {
        return nil
    }
    
    // Example of exposed function dependencies can use
    func (s *MyService) Lookup( name string ) interface{} {
        return nil
    }

## Bootstrap

Every application requires a simple bootstrap.
This bootstrap consists of a single call to `kernel.Launch()` providing one or more services that you want to start.
The kernel will then begin with these, injecting any dependencies required.

The order they are listed here will be used to determine which one starts first,
although any dependencies within them will be deployed first and can override this sequence.

The following is an example from the code that generates [area51.dev](https://area51.dev/):

    package main
    
    import (
        "github.com/peter-mount/documentation/tools/hugo"
        "github.com/peter-mount/documentation/tools/pdf"
        "github.com/peter-mount/go-kernel/v2"
        "log"
    )
    
    func main() {
        err := kernel.Launch(
            &hugo.Hugo{},
            &pdf.PDF{},
        )
        if err != nil {
            log.Fatal(err)
        }
    }

Here it defines two services which must be started. Each one is declared by passing an empty instance of the required
Service's to
the Launch method.

## Lifecycle

This is documented fully in the main documentation, but the following table lists them in the order the kernel will call
them.

All of these are optional.

| Seq | Lifecycle | Description                                                      | Example                                        |
| --- |-----------|------------------------------------------------------------------|------------------------------------------------|
| 1 | Inject    | Resolve dependencies                                             | field injection with `kernel:inject` tag       |
| 2 | Init      | Declare command line flags                                       | func (s *Example) Init(k *kernel.Kernel) error |
| 3 | PostInit  | Verify state is correct, e.g. check command line flags are valid | func (s *Example) PostInit() error             |
| 4 | Start     | Starts the service, open any resources like Databases etc        | func (s *Service) Start() error                |
| 5 | Run       | Run any tasks the service requires                               | func (s *Service) Run() error                  |
| 6 | Stop      | Stops a service                                                  | func (s *Service) Stop()                       |

As the kernel runs through each lifecycle stage, if any method returns an error then the kernel will stop at that point.
If the failure occurs in the Start or Run stages then the Stop stage will be invoked ensuring that any services that
have been started at that point will be stopped.

### Notes
* The `Inject` and `Init` lifecycles are actually the same one but done in that order. I cannot show this in markdown for this page.
* You must *NOT* create any external resources like opening files, databases etc. before the `Start` stage.
* `Start()` is optional. If a service does not implement this function then it's marked as implicitly started.
* `Run()` is deprecated for most purposes. If a service is to perform some task like scan a disk for files it should use the Task api with the worker queue.
* `Stop()` is optional and, it is perfectly valid to implement `Stop()` without a corresponding `Start()` function.
When a service has started (with or without a `Start()` function) it is marked as started so if it implements `Stop()` then that method will be called to clean up the service.
