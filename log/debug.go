package log

import (
	"github.com/peter-mount/go-kernel/v2"
	"os"
)

func init() {
	kernel.Register(&debug{})
}

// Dummy service which adds the -v flag
type debug struct {
	Verbose *bool   `kernel:"flag,v,Verbose"`
	Times   *bool   `kernel:"flag,vt,Verbose with time stamps"`
	Nano    *bool   `kernel:"flag,vn,Verbose with time stamps to nanosecond"`
	Output  *string `kernel:"flag,vo,Verbose to file instead of stderr"`
	output  *os.File
}

func (d *debug) Start() error {
	verbose = *d.Verbose || *d.Times || *d.Nano || *d.Output != ""
	showTime = *d.Times || *d.Nano
	showNano = *d.Nano

	if *d.Output != "" {
		switch *d.Output {
		case "stderr":
			writer = os.Stderr
		case "stdout":
			writer = os.Stdin
		default:
			f, err := os.Create(*d.Output)
			if err != nil {
				return err
			}
			d.output = f
			writer = f
		}
	}

	if writer == nil {
		writer = os.Stderr
	}

	return nil
}

func (d *debug) Stop() {
	if d.output != nil {
		_ = d.output.Close()
	}
}
