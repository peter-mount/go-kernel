package log

import (
	"fmt"
	"io"
	"sync"
	"time"
)

var verbose bool
var showTime bool
var showNano bool
var buf []byte
var writer io.Writer
var mutex sync.Mutex

func IsVerbose() bool { return verbose }

func output(s string) error {
	now := time.Now()

	mutex.Lock()
	defer mutex.Unlock()

	buf = buf[:0]

	if showTime {
		year, month, day := now.Date()
		itoa(&buf, year, 4)
		buf = append(buf, '/')
		itoa(&buf, int(month), 2)
		buf = append(buf, '/')
		itoa(&buf, day, 2)
		buf = append(buf, ' ')

		hour, min, sec := now.Clock()
		itoa(&buf, hour, 2)
		buf = append(buf, ':')
		itoa(&buf, min, 2)
		buf = append(buf, ':')
		itoa(&buf, sec, 2)

		if showNano {
			buf = append(buf, '.')
			itoa(&buf, now.Nanosecond()/1e3, 6)
		}

		buf = append(buf, ' ')
	}

	buf = append(buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}
	_, err := writer.Write(buf)
	return err
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
func Println(v ...interface{}) {
	if verbose {
		_ = output(fmt.Sprintln(v...))
	}
}

func Printf(f string, v ...interface{}) {
	if verbose {
		_ = output(fmt.Sprintf(f, v...))
	}
}

func IfVerbose(f func()) {
	if verbose {
		f()
	}
}
