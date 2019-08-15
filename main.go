package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		filename  = kingpin.Flag("filename", "coverage report from go test").Default("coverage.out").String()
		tag       = kingpin.Flag("tag", "comment tag to exclude blocks from mandatory coverage").Default("nocover").String()
		threshold = kingpin.Flag("threshold", "minimum required overall coverage").Default("-1").Float64()
	)
	kingpin.Parse()

	if err := CheckCoverage(*filename, *tag); err != nil {
		// nocover
		fmt.Println(err)
		os.Exit(1)
	}
	if *threshold > 0 {
		if err := CheckThreshold(*filename, *threshold); err != nil {
			// nocover
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
