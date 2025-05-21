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
	"path/filepath"
	"strings"
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

	// Security check: prevent path escape
	if clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		log.Printf("[SECURITY] attempted directory escape: %q", relPath)
		return fmt.Errorf("invalid directory path %q: path escapes root", relPath)
	}

	parts := strings.Split(clean, string(os.PathSeparator))
	current := ""

	for _, part := range parts {
		// Skip empty parts (shouldn't happen after Clean but just in case)
		if part == "" {
			continue
		}

		current = filepath.Join(current, part)

		_, err := root.Stat(current)
		if err == nil {
			continue // exists
		}

		if os.IsNotExist(err) {
			err = root.Mkdir(current, perm)
			if err != nil {
				return fmt.Errorf("failed to create directory %q: %w", current, err)
			}
		} else {
			return fmt.Errorf("failed to stat directory %q: %w", current, err)
		}
	}

	return nil
}
