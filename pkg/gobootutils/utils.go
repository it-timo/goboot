/*
Package gobootutils provides shared helpers for filesystem-safe generation flows.
*/
package gobootutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

// EnsureDir creates relPath in root (like MkdirAll) while staying root-scoped.
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

// isPathEscapingRoot reports whether a cleaned path attempts upward traversal.
func isPathEscapingRoot(clean string) bool {
	escapeTry := clean == ".." ||
		strings.HasPrefix(clean, "../") ||
		strings.Contains(clean, "/../")

	return escapeTry
}

// ensurePathExists ensures a single path segment exists in root.
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

// CloseFileWithErr closes a file and intentionally ignores close errors.
func CloseFileWithErr(curFile *os.File) {
	_ = curFile.Close()
}

// ComparePaths resolves both paths and enforces either equality or inequality.
// If forceDiffer is true, equal paths return an error; otherwise unequal paths return an error.
func ComparePaths(first, second string, forceDiffer bool) error {
	firstAbs, err := filepath.Abs(filepath.Clean(first))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute first path: %w", err)
	}

	secondAbs, err := filepath.Abs(filepath.Clean(second))
	if err != nil {
		return fmt.Errorf("failed to resolve absolute second path: %w", err)
	}

	if forceDiffer {
		if firstAbs == secondAbs {
			return fmt.Errorf("first and second path must be different: %q == %q", firstAbs, secondAbs)
		}

		return nil
	}

	if firstAbs != secondAbs {
		return fmt.Errorf("first and second path must be the same: %q != %q", firstAbs, secondAbs)
	}

	return nil
}

// CreateRootDir creates `<targetDir>/<name>` and opens it as *os.Root.
func CreateRootDir(targetDir, name string) (*os.Root, error) {
	cleanPath, err := containedRootPath(targetDir, name)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cleanPath, goboottypes.DirPerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create root dir: %w", err)
	}

	curRoot, err := os.OpenRoot(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open root: %w", err)
	}

	return curRoot, nil
}

func containedRootPath(targetDir, name string) (string, error) {
	targetAbs, err := filepath.Abs(filepath.Clean(targetDir))
	if err != nil {
		return "", fmt.Errorf("failed to resolve target dir: %w", err)
	}

	if strings.TrimSpace(name) == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("invalid root dir name %q", name)
	}

	cleanPath, err := filepath.Abs(filepath.Join(targetAbs, name))
	if err != nil {
		return "", fmt.Errorf("failed to resolve root dir: %w", err)
	}

	relPath, err := filepath.Rel(targetAbs, cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to compare root dir with target dir: %w", err)
	}

	if relPath == "." || strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) || relPath == ".." {
		return "", fmt.Errorf("invalid root dir name %q: path escapes target dir", name)
	}

	return cleanPath, nil
}

// EnforceTemplateSourceLimits rejects template sources that exceed file count or total byte limits.
func EnforceTemplateSourceLimits(sourcePath string, maxFiles int, maxTotalBytes int64) error {
	if maxFiles <= 0 || maxTotalBytes <= 0 {
		return nil
	}

	var (
		fileCount  int
		totalBytes int64
	)

	err := filepath.WalkDir(sourcePath, func(_ string, dirEntry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if dirEntry.IsDir() {
			return nil
		}

		fileCount++
		if fileCount > maxFiles {
			return fmt.Errorf("template source exceeds file limit: %d > %d", fileCount, maxFiles)
		}

		info, err := dirEntry.Info()
		if err != nil {
			return fmt.Errorf("failed to read template file metadata: %w", err)
		}

		totalBytes += info.Size()
		if totalBytes > maxTotalBytes {
			return fmt.Errorf("template source exceeds byte limit: %d > %d", totalBytes, maxTotalBytes)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed template source validation: %w", err)
	}

	return nil
}
