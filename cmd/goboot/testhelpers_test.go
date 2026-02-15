package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func testRepoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to resolve repo root")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func loadTestFixture(filename string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(testRepoRoot(), "testdata", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to load test fixture %s: %w", filename, err)
	}

	return data, nil
}

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
