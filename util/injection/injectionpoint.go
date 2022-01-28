package injection

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

type Point struct {
	f  int
	sf reflect.StructField
	tv reflect.Value
	t  reflect.Type
}

func (ip *Point) Type() reflect.Type {
	return ip.t
}

func (ip *Point) StructField() reflect.StructField {
	return ip.sf
}

func Of(f int, sf reflect.StructField, tv reflect.Value) (*Point, error) {
	ip := &Point{
		f:  f,
		sf: sf,
		tv: tv,
	}

	switch ip.sf.Type.Kind() {
	case reflect.Interface:
		ip.t = ip.sf.Type

	case reflect.Ptr:
		ip.t = ip.sf.Type.Elem()

	default:
		log.Println(ip.sf.Type, ip.sf.Type.Kind())
		return nil, ip.Errorf("must be a pointer")
	}

	return ip, nil
}

func (ip *Point) Errorf(f string, args ...interface{}) error {
	tvt := ip.tv.Elem().Type()
	return fmt.Errorf("injection failed \"%s %s\" in %s/%s: %s",
		ip.sf.Name, ip.sf.Type,
		tvt.PkgPath(), tvt.Name(),
		fmt.Sprintf(f, args...))
}

func (ip *Point) Error(err error) error {
	return ip.Errorf("%v", err)
}

// Get returns an accessible Value for a field in another Value
func (ip *Point) Get() reflect.Value {
	// Get the Value for the field in the service struct
	tf := ip.tv.Elem().Field(ip.f)

	// Some magic, this provides us write access the field even if it's unexported.
	// Without this tf.Set() will fail if the field is unexported
	// see: https://stackoverflow.com/a/43918797/6734016
	return reflect.NewAt(tf.Type(), unsafe.Pointer(tf.UnsafeAddr())).Elem()
}

// Set sets a field in a value with a specific instance of an interface
func (ip *Point) Set(val interface{}) {
	// Convert our resolved service into a Value then convert to the field's type
	vv := reflect.ValueOf(val)
	ip.Get().Set(vv.Convert(ip.sf.Type))
}

// New creates a new instance of a Type
func (ip *Point) New() interface{} {
	if ip.t == nil {
		// Must call EnforcePointer() first
		panic(ip.Errorf("not yet dereferenced"))
	}
	return reflect.New(ip.t).Interface()
}
