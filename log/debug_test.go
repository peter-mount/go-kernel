package log

import (
	"strings"
	"testing"
)

func TestVerboseFlags(t *testing.T) {
	output := ""
	bools := []bool{false, true}
	for _, vn := range bools {
		for _, vt := range bools {
			for _, v := range bools {
				var n []string
				if v {
					n = append(n, "v")
				}
				if vt {
					n = append(n, "vt")
				}
				if vn {
					n = append(n, "vn")
				}
				if len(n) == 0 {
					n = append(n, "none")
				}

				t.Run(strings.Join(n, " "), func(t *testing.T) {

					verbose = false
					showTime = false
					showNano = false

					// Keep a local copy in this scope as the variables will change before
					// this code runs
					aV, aVT, aVN := v, vt, vn

					// Expected values
					eV := aV || aVT || aVN
					eVT := aVT || aVN
					eVN := aVN

					debug := &debug{
						Verbose: &aV,
						Times:   &aVT,
						Nano:    &aVN,
						Output:  &output,
					}
					if err := debug.Start(); err != nil {
						t.Errorf("debug.Start() returned error %v", err)
					} else if verbose != eV || showTime != eVT || showNano != eVN {
						t.Errorf("Unexpected state: got verbose %v showTime %v showNano %v wanted verbose %v showTime %v showNano %v",
							verbose, showTime, showNano,
							eV, eVT, eVN)
					}
				})
			}
		}
	}
}
