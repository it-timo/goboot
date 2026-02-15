package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func repoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to resolve repo root")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// loadTestFixture reads a YAML file from the testdata directory.
func loadTestFixture(filename string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(repoRoot(), "testdata", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to load test fixture %s: %w", filename, err)
	}

	return data, nil
}

// loadTestFixtureWithVars reads a YAML file and replaces template variables.
// Variables are specified as ${VARNAME} in the fixture and provided as a map.
func loadTestFixtureWithVars(filename string, vars map[string]string) ([]byte, error) {
	data, err := loadTestFixture(filename)
	if err != nil {
		return nil, err
	}

	content := string(data)
	for key, val := range vars {
		content = strings.ReplaceAll(content, "${"+key+"}", val)
	}

	return []byte(content), nil
}
