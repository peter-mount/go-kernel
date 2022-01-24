package kernel

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
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
			if err := k.injectField(sk, f, sf, tv); err != nil {
				return err
			}
		}
	}
	return nil
}

// injectField handles the injection of a specific field
func (k *Kernel) injectField(tag string, f int, sf reflect.StructField, tv reflect.Value) error {
	// Run through each param in the tag
	for _, tagTerm := range strings.Split(tag, ",") {
		switch tagTerm {
		case "-":
			// Do nothing. Think how json/xml uses this. If it's the first entry then this ignores this field
			return nil

		case "inject":
			// inject a dependency
			if err := k.injectService(f, sf, tv); err != nil {
				return err
			}

		default:
			// Fail with an unsupported tag value
			return fmt.Errorf("unsupported kernel tag %q", tagTerm)
		}
	}
	return nil
}

// injectService injects a dependency into the service structure
func (k *Kernel) injectService(f int, sf reflect.StructField, tv reflect.Value) error {

	if sf.Type.Kind() != reflect.Ptr {
		return fmt.Errorf("injection failed \"%s %s\" not a pointer to a Service", sf.Name, sf.Type)
	}

	inst := reflect.New(sf.Type.Elem()).Interface()
	if sInst, ok := inst.(Service); ok {
		// Add the service in the traditional way, returning us the deployed instance
		resolvedService, err := k.AddService(sInst)
		if err != nil {
			return err
		}

		// Get the Value for the field in the service struct
		tf := tv.Elem().Field(f)

		// Some magic, this provides us write access the field even if it's unexported.
		// Without this tf.Set() will fail if the field is unexported
		// see: https://stackoverflow.com/a/43918797/6734016
		tf = reflect.NewAt(tf.Type(), unsafe.Pointer(tf.UnsafeAddr())).Elem()

		// Convert our resolved service into a Value then convert to the field's type
		vv := reflect.ValueOf(resolvedService)
		tf.Set(vv.Convert(sf.Type))
	} else {
		return fmt.Errorf("injection failed \"%s %s\" not a Service", sf.Name, sf.Type)
	}

	return nil
}
