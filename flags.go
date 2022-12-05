package kernel

import (
	"flag"
	"github.com/peter-mount/go-kernel/v2/util/injection"
	"reflect"
	"strconv"
)

// injectFlag - kernel:"flag:name:desc:default" - default is optional
func (k *Kernel) injectFlag(tags []string, ip *injection.Point) error {
	switch ip.Type().Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(getFlagDefault(tags, "false"))
		if err != nil {
			return ip.Error(err)
		}
		ip.Set(flag.Bool(getFlagName(tags, ip), v, getFlagDesc(tags, ip)))

	case reflect.String:
		v := getFlagDefault(tags, "")
		ip.Set(flag.String(getFlagName(tags, ip), v, getFlagDesc(tags, ip)))

	case reflect.Int:
		v, err := strconv.ParseInt(getFlagDefault(tags, "0"), 10, 64)
		if err != nil {
			return ip.Error(err)
		}
		ip.Set(flag.Int(getFlagName(tags, ip), int(v), getFlagDesc(tags, ip)))

	case reflect.Int64:
		v, err := strconv.ParseInt(getFlagDefault(tags, "0"), 10, 64)
		if err != nil {
			return ip.Error(err)
		}
		ip.Set(flag.Int64(getFlagName(tags, ip), v, getFlagDesc(tags, ip)))

	case reflect.Float64:
		v, err := strconv.ParseFloat(getFlagDefault(tags, "0.0"), 64)
		if err != nil {
			return ip.Error(err)
		}
		ip.Set(flag.Float64(getFlagName(tags, ip), v, getFlagDesc(tags, ip)))

	default:
		return ip.Errorf("unsupported flag type %q", ip.Type())
	}

	return nil
}

func getFlagName(tags []string, ip *injection.Point) string {
	if len(tags) > 0 && tags[0] != "" {
		return tags[0]
	}
	return ip.StructField().Name
}

func getFlagDesc(tags []string, ip *injection.Point) string {
	if len(tags) > 1 && tags[1] != "" {
		return tags[1]
	}
	return getFlagName(tags, ip)
}

func getFlagDefault(tags []string, d string) string {
	if len(tags) > 2 && tags[2] != "" {
		return tags[2]
	}
	return d
}
