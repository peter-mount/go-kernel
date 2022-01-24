package kernel

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

// injectFlag - kernel:"flag:name:desc:default" - default is optional
func (k *Kernel) injectFlag(tags []string, f int, sf reflect.StructField, tv reflect.Value) error {
	t := sf.Type
	if sf.Type.Kind() != reflect.Ptr {
		return fmt.Errorf("flag %s not pointer", sf.Name)
	}
	t = t.Elem()

	switch t.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(getFlagDefault(tags, sf, "false"))
		if err != nil {
			return err
		}
		setVal(f, sf, tv, flag.Bool(getFlagName(tags, sf), v, getFlagDesc(tags, sf)))

	case reflect.String:
		v := getFlagDefault(tags, sf, "")
		setVal(f, sf, tv, flag.String(getFlagName(tags, sf), v, getFlagDesc(tags, sf)))

	case reflect.Int:
		v, err := strconv.ParseInt(getFlagDefault(tags, sf, "0"), 10, 64)
		if err != nil {
			return err
		}
		setVal(f, sf, tv, flag.Int(getFlagName(tags, sf), int(v), getFlagDesc(tags, sf)))

	case reflect.Int64:
		v, err := strconv.ParseInt(getFlagDefault(tags, sf, "0"), 10, 64)
		if err != nil {
			return err
		}
		setVal(f, sf, tv, flag.Int64(getFlagName(tags, sf), v, getFlagDesc(tags, sf)))

	case reflect.Float64:
		v, err := strconv.ParseFloat(getFlagDefault(tags, sf, "0.0"), 64)
		if err != nil {
			return err
		}
		setVal(f, sf, tv, flag.Float64(getFlagName(tags, sf), v, getFlagDesc(tags, sf)))

	default:
		return fmt.Errorf("unsupported flag type %q on %q", t, sf.Name)
	}

	return nil
}

func getFlagName(tags []string, sf reflect.StructField) string {
	if len(tags) > 0 && tags[0] != "" {
		return tags[0]
	}
	return sf.Name
}

func getFlagDesc(tags []string, sf reflect.StructField) string {
	if len(tags) > 1 && tags[1] != "" {
		return tags[1]
	}
	return getFlagName(tags, sf)
}

func getFlagDefault(tags []string, _ reflect.StructField, d string) string {
	if len(tags) > 2 && tags[2] != "" {
		return tags[2]
	}
	return d
}
