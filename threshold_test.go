package main_test

import (
	"os"
	"os/exec"
	"testing"

	main "github.com/digio/covermate"
	"github.com/stretchr/testify/assert"
)

func TestCheckThreshold_NoCoverfile(t *testing.T) {
	os.Remove("testdata/coverage.out")

	err := main.CheckThreshold("testdata/bad.out", 100.0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCheckThreshold(t *testing.T) {
	os.Remove("testdata/coverage.out")
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
	cmd.Dir = "testdata"
	assert.NoError(t, cmd.Run())

	err := main.CheckThreshold("testdata/coverage.out", 100.0)
	assert.Error(t, err)

	err = main.CheckThreshold("testdata/coverage.out", 10.0)
	assert.NoError(t, err)
}
