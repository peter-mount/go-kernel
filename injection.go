package kernel

import (
	"github.com/peter-mount/go-kernel/v2/util/injection"
	"reflect"
	"strings"
)

// inject implements injection using field tag's
func (k *Kernel) inject(v interface{}) error {
	if v == nil {
		return nil
	}

	tv := reflect.ValueOf(v)

	t := tv.Type()
	if t == nil {
		return nil
	}

	// Run through each field in the service
	elem := t.Elem()
	numField := elem.NumField()
	for f := 0; f < numField; f++ {
		sf := elem.Field(f)
		if sk, ok := sf.Tag.Lookup("kernel"); ok {
			ip, err := injection.Of(f, sf, tv)
			if err == nil {
				err = k.injectField(sk, ip)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type kernelInjector func(tags []string, ip *injection.Point) error

// injectField handles the injection of a specific field
func (k *Kernel) injectField(tag string, ip *injection.Point) error {
	tags := strings.Split(tag, ",")

	var injector kernelInjector
	switch tags[0] {
	case "-":
		// Do nothing. Think how json/xml uses this. If it's the first entry then this ignores this field
		return nil

	case "inject":
		// inject a dependency
		injector = k.injectService

	case "worker":
		// Inject the default task.Queue(f, sf, tv)
		injector = k.injectWorker

	case "flag":
		injector = k.injectFlag

	case "config":
		injector = k.injectConfig

	default:
		// Fail with an unsupported tag value
		return ip.Errorf("unsupported kernel tag %q", tags[0])
	}

	if injector != nil {
		if err := injector(tags[1:], ip); err != nil {
			return err
		}
	}

	return nil
}

// injectService injects a dependency into the service structure
func (k *Kernel) injectService(_ []string, ip *injection.Point) error {

	// See if we already have the service deployed.
	// At this point it could be either a Service or an API
	t := ip.Type()
	n := getServiceName(t)
	if resolvedService, exists := k.index[n]; exists {
		ip.Set(resolvedService)
		return nil
	}

	inst := ip.New()
	if sInst, ok := inst.(Service); ok {
		// Add the service in the traditional way, returning us the deployed instance
		resolvedService, err := k.AddService(sInst)
		if err != nil {
			return err
		}

		ip.Set(resolvedService)
	} else {
		return ip.Errorf("not a Service")
	}

	return nil
}

func (k *Kernel) injectWorker(_ []string, ip *injection.Point) error {
	resolvedService, err := k.AddService(&Worker{})
	if err != nil {
		return err
	}
	ip.Set(resolvedService)
	return nil
}
