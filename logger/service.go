package logger

import (
	"flag"
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Kernel service to expose the logger.
// You can either use the command line arguments or provide the printer definition as an IPP url in the IPP environment
// variable.
//
// e.g. ipp://host/printers/printer
//      ipp://user:pass@host/printers/printer
//
// To use tls then use ipps instead of ipp. Yes it's not standard but the service will recognise that it should use TLS.
//
// Also append ?escpos=true if the remote printer is an escpos printer, e.g. a Receipt printer.
//
// e.g. ipp://host/printers/printer?escpos=true
//      ipp://user:pass@host/printers/printer?escpos=true
//
type LoggerService struct {
	escpos     *bool
	ipp        *string
	ippEnabled bool
	host       string
	port       int
	printer    string
	username   string
	password   string
	useTLS     bool
}

const (
	printersPrefix = "/printers/"
)

func (s *LoggerService) Name() string {
	return "LoggerService"
}
func (s *LoggerService) Init(k *kernel.Kernel) error {
	s.escpos = flag.Bool("logescpos", false, "Set ESC/POS for report format")
	s.ipp = flag.String("logipp", "", "IPP address of printer")
	return nil
}

func (s *LoggerService) PostInit() error {
	if s.ipp == nil || *s.ipp == "" {
		*s.ipp = os.Getenv("IPP")
	}

	if s.ipp != nil && *s.ipp != "" {

		log.Println(*s.ipp)
		u, err := url.Parse(*s.ipp)
		if err != nil {
			return err
		}

		s.ippEnabled = true

		s.useTLS = strings.HasPrefix(u.Scheme, "s")

		user := u.User
		if user != nil {
			s.username = user.Username()
			s.password, _ = user.Password()
		}

		s.host = u.Hostname()

		if u.Port() != "" {
			s.port, err = strconv.Atoi(u.Port())
			if err != nil {
				return err
			}
		} else if u.Scheme == "ipp" || u.Scheme == "ipps" {
			s.port = 631
		} else if s.useTLS {
			s.port = 443
		} else {
			s.port = 80
		}

		s.printer = u.Path
		if strings.HasPrefix(s.printer, printersPrefix) {
			s.printer = s.printer[len(printersPrefix):]
		}

		if u.Query().Get("escpos") != "" {
			*s.escpos = true
		}

	}

	return nil
}

func (s *LoggerService) submitLog(l *Logger, err error) error {
	if s.ippEnabled {
		client := ipp.NewIPPClient(s.host, s.port, s.username, s.password, s.useTLS)

		now := time.Now()

		// Optional Error log line
		failure := ""
		if err != nil {
			failure = "\n    Error: " + err.Error()
		}

		host, err := os.Hostname()
		if err != nil {
			host = "unknown"
		}

		l.write([]byte(fmt.Sprintf("\n"+
			"   Report: %s\n"+
			"  Started: %s\n"+
			"Completed: %s\n"+
			" Duration: %s%s\n"+
			"     Host: %s\n",
			l.title,
			l.startTime.Format(time.RFC1123),
			now.Format(time.RFC1123),
			// Printer can't handle "µ" so replace with "u"
			strings.ReplaceAll(now.Sub(l.startTime).String(), "µ", "u"),
			failure,
			host,
		))...)

		buffer := l.Buffer()

		doc := ipp.Document{
			Document: buffer,
			Name:     l.title,
			Size:     buffer.Len(),
			MimeType: ipp.MimeTypeOctetStream,
		}

		jobAttributes := make(map[string]interface{})
		jobAttributes[ipp.OperationAttributeJobName] = l.title

		_, err = client.PrintJob(doc, s.printer, jobAttributes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LoggerService) Report(title string, f func(l *Logger) error) error {
	l := NewReport(title, log.New(os.Stdout, "", log.LstdFlags))
	if *s.escpos {
		l.EscPos()
	}

	err := f(l)
	if err != nil {
		l.Println("Failure:", err)
	}

	err2 := s.submitLog(l, err)
	if err2 != nil {
		return err2
	}

	return err
}
