package main

import (
	"fmt"
	"math"

	"golang.org/x/tools/cover"
)

// CheckThreshold calculates the proportion of covered statements referenced by the coverage
// report at `filename`.  This figure is then compared with limit to determine if overall coverage is adequate.
func CheckThreshold(filename string, limit float64) error {
	var t, c int
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return err
	}
	for _, p := range profiles {
		for _, block := range p.Blocks {
			t += block.NumStmt
			if block.Count > 0 {
				c += block.NumStmt
			}
		}
	}

	// avoid div0
	if t == 0 {
		t = 1
	}
	overall := math.Round(1000.0*float64(c)/float64(t)) / 10.0
	if overall < limit {
		return fmt.Errorf("Coverage %.1f is below required threshold %.1f", overall, limit)
	}
	return nil
}
