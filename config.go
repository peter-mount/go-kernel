package kernel

import (
	"bufio"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/injection"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

// dynamicConfig is an extensible yaml based config file format.
//
// Services can add their own config handlers to this instance so that
// when it loads the configuration from a file it updates each service
// before they start.
//
// In the config file, the yaml consists of objects, one per service.
type dynamicConfig struct {
	filename *string                 `kernel:"flag,config,Configuration file,config.yaml"`
	entries  map[string]*configEntry // Map of entries
	files    map[string]interface{}  // Map used to prevent infinite loop loading files
}

type configEntry struct {
	name            string             // Name of entry in yaml file
	config          interface{}        // Config to inject into
	injectionPoints []*injection.Point // Injection points to receive this config
}

// Add a named config entry. Returns an Error if the name is already in use
func (dc *dynamicConfig) add(name string, ip *injection.Point) error {
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

func (dc *dynamicConfig) Start() error {
	dc.files = make(map[string]interface{})
	return dc.processFile(*dc.filename)
}

const (
	includePrefix = "#include "
)

func (dc *dynamicConfig) processFile(filename string) error {
	absFileName, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	// Prevent loading the same file twice
	if _, exists := dc.files[absFileName]; exists {
		return fmt.Errorf("already read %q, possible infinite loop", absFileName)
	}
	dc.files[absFileName] = true

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > 0 {
			c := line[0]
			switch {

			case strings.HasPrefix(line, includePrefix):
				// Process any existing block
				if err := dc.processBlock(lines); err != nil {
					return err
				}

				// Clear the current block, so we start a fresh once the included file has been read
				lines = nil

				// Everything after includePrefix is the filename, trim excess white space
				nextFilename := strings.TrimSpace(line[len(includePrefix):])
				l := len(nextFilename)

				// If filename is not empty and is wrapped with " then load it
				if l > 2 && nextFilename[0] == '"' && nextFilename[l-1] == '"' {
					err = dc.processFile(nextFilename[1 : l-1])
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("invalid #include %q in %q", line, filename)
				}

			case c == '#':
				// Strip comments

			case (c >= 'A' && c < 'Z') || (c >= 'a' && c < 'z'):
				// Line starts with a character then it's the start of a block

				// Process any existing block
				if err := dc.processBlock(lines); err != nil {
					return err
				}

				// Start the block with this line
				lines = []string{line}

			default:
				// Append to the line list
				lines = append(lines, line)
			}
		}
	}

	// Handle last config block
	return dc.processBlock(lines)
}

func (dc *dynamicConfig) processBlock(lines []string) error {
	if len(lines) > 0 {
		i := strings.Index(lines[0], ":")
		if i < 0 {
			return fmt.Errorf("invalid config parameter %s", strings.Join(lines, "\\n"))
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
	sv, err := k.AddService(&dynamicConfig{})
	if err != nil {
		return err
	}
	dc := sv.(*dynamicConfig)

	// Add the injection point to the section
	return dc.add(configSectionName, ip)
}
