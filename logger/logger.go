package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

// ReportLogger used with the Logger package to write logs to the console and a secondary
// source, usually a printer
type Logger struct {
	mu         sync.Mutex   // ensures atomic writes; protects the following fields
	title      string       // Title of the report
	flag       int          // properties
	buf        []byte       // for accumulating text to write
	logger     *log.Logger  // Logger to also send log calls to
	textBuf    bytes.Buffer // Text buffer for final report
	lt         time.Time    // Last time entry logged
	started    bool         // Set after the first log entry
	indent     []byte       // Line prefix to indent new/broken lines
	lineWidth  int          // width of a line
	lineIndent int          // Size of indent
	escpos     bool         // Escpos mode
	startTime  time.Time    // Start time
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(title string, flag int, width int, l *log.Logger) *Logger {
	return &Logger{
		title:     title,
		flag:      flag,
		logger:    l,
		lineWidth: width,
		startTime: time.Now(),
	}
}

func NewReport(title string, l *log.Logger) *Logger {
	return New(title, log.LstdFlags, 80, l)
}

func (l *Logger) EscPos() *Logger {
	if !l.started {
		l.escpos = true
		l.lineWidth = 42
	}
	return l
}

func (l *Logger) Buffer() *bytes.Buffer {
	if l.started && l.escpos {
		l.textBuf.Write([]byte{'\n', '\n', '\n', '\x1B', 'M', 2, '\x1B', 'V', 'A', '0', '\xFA'})
	}
	return &l.textBuf
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

func (l *Logger) write(b ...byte) {
	_, _ = l.textBuf.Write(b)
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {

	if !l.started {
		if l.escpos {
			l.write('\x1B', '@')
		}
		if l.title != "" {
			if l.escpos {
				l.write('\x1B', 'M', '\x00', '\x1B', 'E', 1)
			}
			l.write('\n')
			t := l.title
			if len(t) >= l.lineWidth {
				t = t[:l.lineWidth]
			}
			l.write([]byte(t)...)
			l.write('\n', '\n')
			if l.escpos {
				l.write('\x1B', 'E', 0)
			}
		}
		if l.escpos {
			l.write('\x1B', 'M', '\x01')
		}
	}

	if l.flag&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flag&log.LUTC != 0 {
			t = t.UTC()
		}

		// If date then only show at start of the report or if it changes
		if l.flag&log.Ldate != 0 {
			year, month, day := t.Date()
			lyear, lmonth, lday := l.lt.Date()
			if lyear != year && lmonth != month && lday != day {
				itoa(buf, year, 4)
				*buf = append(*buf, '/')
				itoa(buf, int(month), 2)
				*buf = append(*buf, '/')
				itoa(buf, day, 2)
				*buf = append(*buf, '\n')
			}
		}

		if l.flag&(log.Ltime|log.Lmicroseconds) != 0 {
			bs := len(*buf)

			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)

			if l.flag&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}

			*buf = append(*buf, ' ')

			if !l.started {
				width := len(*buf) - bs
				l.indent = bytes.Repeat([]byte{' '}, width)
				l.lineWidth = l.lineWidth - width
			}
		}

		// Keep track of the new latest time
		l.lt = t
	}

	l.started = true
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, os string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = l.buf[:0]

	l.formatHeader(&l.buf, now, file, line)

	s := os
	lf := len(s) == 0 || s[len(s)-1] != '\n'
	lc := false
	for len(s) > 0 {
		if lc {
			l.buf = append(l.buf, l.indent...)
		} else {
			lc = true
		}

		if len(s) < l.lineWidth {
			l.buf = append(l.buf, s[:]...)
			s = ""
		} else {
			spc := strings.Index(s, "\n")
			if spc < 0 {
				spc = len(s)
			}
			if spc > l.lineWidth {
				spc = l.lineWidth
			}

			// Find last space before lf and use that
			p := spc
			for ; p >= 0 && s[p] != ' '; p-- {
			}
			if p > 0 {
				spc = p
			}

			if spc > len(s) {
				spc = len(s)
			}

			l.buf = append(l.buf, s[0:spc]...)
			l.buf = append(l.buf, '\n')
			if spc < len(s) {
				s = s[spc+1:]
			} else {
				s = ""
			}
		}
	}

	if lf {
		l.buf = append(l.buf, '\n')
	}

	_, err := l.textBuf.Write(l.buf)
	if err != nil {
		return err
	}

	if l.logger != nil {
		return l.logger.Output(calldepth, os)
	}

	return nil
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	_ = l.Output(2, fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	_ = l.Output(2, fmt.Sprintln(v...))
}

// Flags returns the output flags for the logger.
func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
	l.logger.SetFlags(flag)
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		return l.logger.Prefix()
	}
	return ""
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		l.logger.SetPrefix(prefix)
	}
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logger != nil {
		return l.logger.Writer()
	}
	return nil
}
