package regeneration

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func readManifest(projectPath string) (*Manifest, error) {
	manifestPath := filepath.Join(projectPath, ManifestFileName)
	content, err := os.ReadFile(manifestPath) // #nosec G304 -- projectPath is validated by the caller.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read regeneration manifest: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)

	var manifest Manifest

	err = decoder.Decode(&manifest)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to decode regeneration manifest: %w", err)
	}

	err = validateManifest(manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func validateManifest(manifest Manifest) error {
	if manifest.SchemaVersion != ManifestSchemaVersion {
		return fmt.Errorf("unsupported regeneration manifest schema: %d", manifest.SchemaVersion)
	}

	seenPaths := make(map[string]struct{}, len(manifest.Files))

	for _, file := range manifest.Files {
		if err := validateManifestPath(file.Path); err != nil {
			return err
		}

		if _, found := seenPaths[file.Path]; found {
			return fmt.Errorf("duplicate path in regeneration manifest: %q", file.Path)
		}

		seenPaths[file.Path] = struct{}{}

		digest, err := hex.DecodeString(file.SHA256)
		if err != nil || len(digest) != sha256.Size {
			return fmt.Errorf("invalid SHA-256 for manifest path %q", file.Path)
		}

		if file.Mode > 0o777 {
			return fmt.Errorf("invalid file mode for manifest path %q: %o", file.Path, file.Mode)
		}
	}

	return nil
}

func validateManifestPath(path string) error {
	cleanPath := filepath.Clean(path)
	if path == "" || path == "." || filepath.IsAbs(path) || cleanPath != path ||
		strings.HasPrefix(cleanPath, ".."+string(os.PathSeparator)) || cleanPath == ".." ||
		path == ManifestFileName {
		return fmt.Errorf("invalid path in regeneration manifest: %q", path)
	}

	return nil
}

func encodeManifest(manifest Manifest) ([]byte, error) {
	sort.Slice(manifest.Files, func(firstIndex, secondIndex int) bool {
		return manifest.Files[firstIndex].Path < manifest.Files[secondIndex].Path
	})
	manifest.Services = sortedServices(manifest.Services)

	content, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to encode regeneration manifest: %w", err)
	}

	return content, nil
}

func manifestFileMap(manifest *Manifest) map[string]ManifestFile {
	if manifest == nil {
		return map[string]ManifestFile{}
	}

	files := make(map[string]ManifestFile, len(manifest.Files))
	for _, file := range manifest.Files {
		files[file.Path] = file
	}

	return files
}
