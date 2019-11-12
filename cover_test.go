package main_test

import (
	"os"
	"os/exec"
	"testing"

	main "github.com/digio/covermate"
	"github.com/stretchr/testify/assert"
)

func TestCheckCoverage_NoCoverfile(t *testing.T) {
	os.Remove("testdata/coverage.out")

	err := main.CheckCoverage("testdata/bad.out", "nocover")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCheckCoverage_FullCoverage(t *testing.T) {
	os.Remove("testdata/coverage.out")
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out",
		"github.com/digio/covermate/testdata/fullcover")
	cmd.Dir = "testdata"
	cmd.Stdout = os.Stdout
	assert.NoError(t, cmd.Run())

	err := main.CheckCoverage("testdata/coverage.out", "nocover")
	assert.NoError(t, err)
}

func TestCheckCoverage_Mixed(t *testing.T) {
	os.Remove("testdata/coverage.out")
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
	cmd.Dir = "testdata"
	assert.NoError(t, cmd.Run())

	err := main.CheckCoverage("testdata/coverage.out", "nocover")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Non-tagged code exists without coverage")
}

func TestCheckCoverage_Switch(t *testing.T) {
	os.Remove("testdata/coverage.out")
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "github.com/digio/covermate/testdata/switcher")
	cmd.Dir = "testdata"
	assert.NoError(t, cmd.Run())

	err := main.CheckCoverage("testdata/coverage.out", "nocover")
	assert.NoError(t, err)
}
