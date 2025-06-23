/*
Package utils provides reusable helper functions for safe and scoped filesystem
operations and related logic used throughout the goboot project.

It is designed to encapsulate commonly necessary behaviors, such as recursive
directory creation under a secure root, without duplicating logic across services.

This package emphasizes reusability and minimal external assumptions.
*/
package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/types"
)

// EnsureDir safely creates a nested directory path relative to the given os.Root.
//
// It emulates the behavior of os.MkdirAll but operates entirely within the provided *os.Root context
// to ensure isolation from the host filesystem.
//
// Each segment of the path is verified or created in sequence.
// If a segment already exists, it is skipped.
//
// Parameters:
//   - relPath: The relative directory path to ensure existing (e.g., "a/b/c").
//   - root: The secured *os.Root within which all directories are created.
//   - perm: The permission mode to apply when creating new directories.
//
// Returns an error if any directory creation or stat check fails.
func EnsureDir(relPath string, root *os.Root, perm os.FileMode) error {
	if relPath == "." || relPath == "" {
		return nil
	}

	clean := filepath.Clean(relPath)

	if isPathEscapingRoot(clean) {
		return fmt.Errorf("invalid directory path %q: path escapes root", relPath)
	}

	parts := strings.Split(clean, string(os.PathSeparator))
	current := ""

	for _, part := range parts {
		if part == "" {
			continue
		}

		current = filepath.Join(current, part)

		err := ensurePathExists(current, root, perm)
		if err != nil {
			return err
		}
	}

	return nil
}

// isPathEscapingRoot checks whether the cleaned relative path attempts to escape the root scope.
//
// Returns true if the path includes segments like "../" that indicate upward traversal outside the root.
func isPathEscapingRoot(clean string) bool {
	escapeTry := clean == ".." ||
		strings.HasPrefix(clean, "../") ||
		strings.Contains(clean, "/../")

	if escapeTry {
		log.Printf("[SECURITY] attempted directory escape: %q", clean)
	}

	return escapeTry
}

// ensurePathExists checks if the given path exists within the root and creates it if missing.
//
// If the directory already exists, the call is a no-op.
//
// Returns an error if stat or mkdir operations fail.
func ensurePathExists(path string, root *os.Root, perm os.FileMode) error {
	_, err := root.Stat(path)
	if err == nil {
		return nil // already exists.
	}

	if os.IsNotExist(err) {
		err = root.Mkdir(path, perm)
		if err != nil {
			return fmt.Errorf("failed to create directory %q: %w", path, err)
		}

		return nil
	}

	return fmt.Errorf("failed to stat directory %q: %w", path, err)
}

// CloseFileWithErr closes the given file and logs a warning if an error occurs.
//
// It does not return the error, as this function is intended for use in defer statements
// where failure to close should not interrupt the execution flow.
func CloseFileWithErr(curFile *os.File) {
	err := curFile.Close()
	if err != nil {
		log.Printf("failed to close file: %s", err)
	}
}

// ComparePaths resolves and compares two filesystem paths after cleaning and normalization.
//
// This function ensures deterministic and platform-safe comparison of two paths.
// It converts both paths to absolute, cleaned versions before comparing.
//
// If `forceDiffer` is true:
//   - The function returns an error if both paths resolve to the same absolute location.
//   - This is useful to assert that source and target paths are not accidentally identical.
//
// If `forceDiffer` is false:
//   - The function returns an error if the paths are not equal.
//   - This is useful to assert that two paths resolve to the exact same location.
//
// This utility should be used when validating user-defined paths in config or templates.
//
// Returns:
//   - `nil` if the comparison is valid based on the `forceDiffer` flag.
//   - `error` if the resolution fails or the comparison logic fails.
func ComparePaths(first, second string, forceDiffer bool) error {
	// Resolve and clean the first path.
	firstAbs, err := filepath.Abs(filepath.Clean(first))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute first path: %w", err)
	}

	// Resolve and clean the second path.
	secondAbs, err := filepath.Abs(filepath.Clean(second))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute second path: %w", err)
	}

	// Enforce inequality if required.
	if forceDiffer {
		if firstAbs == secondAbs {
			return fmt.Errorf("first and second path must be different: %q == %q", firstAbs, secondAbs)
		}

		return nil
	}

	// Enforce equality if required.
	if firstAbs != secondAbs {
		return fmt.Errorf("first and second path must be the same: %q != %q", firstAbs, secondAbs)
	}

	return nil
}

// CreateRootDir creates the root directory for the project.
//
// It uses the project name to create a new directory under the target directory if not exist.
//
// The directory is created with the correct permissions and ownership.
//
// It also opens the directory as a Root for further processing.
func CreateRootDir(targetDir, name string) (*os.Root, error) {
	curPath := filepath.Join(targetDir, name)

	cleanPath, err := filepath.Abs(path.Clean(curPath))
	if err != nil {
		return nil, fmt.Errorf("failed to clean path: %w", err)
	}

	err = os.MkdirAll(cleanPath, types.DirPerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create root dir: %w", err)
	}

	curRoot, err := os.OpenRoot(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open root: %w", err)
	}

	return curRoot, nil
}
