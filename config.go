package kernel

import (
	"bufio"
	"fmt"
	"github.com/peter-mount/go-kernel/util/injection"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// DynamicConfig is an extensible yaml based config file format.
//
// Services can add their own config handlers to this instance so that
// when it loads the configuration from a file it updates each service
// before they start.
//
// In the config file, the yaml consists of objects, one per service.
type DynamicConfig struct {
	filename *string                 `kernel:"flag,config,Configuration file,config.yaml"`
	entries  map[string]*configEntry // Map of entries
}

type configEntry struct {
	name            string             // Name of entry in yaml file
	config          interface{}        // Config to inject into
	injectionPoints []*injection.Point // Injection points to receive this config
}

// Add a named config entry. Returns an Error if the name is already in use
func (dc *DynamicConfig) add(name string, ip *injection.Point) error {
	if dc.entries == nil {
		dc.entries = make(map[string]*configEntry)
	}

	e, exists := dc.entries[name]

	// Create the entry on first use
	if !exists {
		e = &configEntry{name: name}
		dc.entries[name] = e

		// The shared instance
		e.config = ip.New()
	}

	// Inject the shared instance
	ip.Set(e.config)

	return nil
}

func (dc *DynamicConfig) Start() error {
	f, err := os.Open(*dc.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// TODO this needs to be first char is alpha
		if len(line) > 0 {
			if line[0] > 64 {
				if err := dc.processLines(lines); err != nil {
					return err
				}
				lines = nil
			}

			// Strip comments
			if line[0] != '#' {
				lines = append(lines, line)
			}
		}
	}

	// Handle last config block
	if err := dc.processLines(lines); err != nil {
		return err
	}

	return nil
}

func (dc *DynamicConfig) processLines(lines []string) error {
	if len(lines) > 0 {
		i := strings.Index(lines[0], ":")
		if i < 0 {
			return fmt.Errorf("invalid")
		}

		// Name of section, keep everything after : but trimmed incase it's a simple config entry
		n := lines[0][:i]
		lines[0] = strings.TrimSpace(lines[0][i+1:])

		if e, exists := dc.entries[n]; exists {
			b := []byte(strings.Join(lines, "\n"))

			if err := yaml.Unmarshal(b, e.config); err != nil {
				return err
			}
		} else {
			// TODO display warning here of an unused section in the read yaml? It isn't an error however
		}
	}
	return nil
}

func (k *Kernel) injectConfig(tags []string, ip *injection.Point) error {

	var configSectionName string
	if len(tags) > 0 {
		// First field in tag is the section name
		configSectionName = tags[0]
	}

	if configSectionName == "" {
		// Default to field name for the yaml section name
		configSectionName = ip.StructField().Name
	}

	// lazy init service
	sv, err := k.AddService(&DynamicConfig{})
	if err != nil {
		return err
	}
	dc := sv.(*DynamicConfig)

	// Add the injection point to the section
	return dc.add(configSectionName, ip)
}
