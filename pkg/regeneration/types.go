/*
Package regeneration plans and transactionally applies generated project trees.
*/
package regeneration

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	// ManifestFileName records files owned by goboot in a generated project.
	ManifestFileName = ".goboot-manifest.yml"
	// ManifestSchemaVersion is the current ownership-manifest schema.
	ManifestSchemaVersion = 1
)

// Policy controls how generation handles existing files that goboot cannot
// safely identify as unchanged managed output.
type Policy string

const (
	// PolicyManaged updates unchanged owned files and rejects user-modified collisions.
	PolicyManaged Policy = "managed"
	// PolicyReplace replaces collisions and removes stale files previously owned by goboot.
	PolicyReplace Policy = "replace"
	// PolicyPreserve keeps collisions and relinquishes ownership of preserved files.
	PolicyPreserve Policy = "preserve"
)

// ParsePolicy validates and normalizes a regeneration policy.
func ParsePolicy(raw string) (Policy, error) {
	policy := Policy(strings.ToLower(strings.TrimSpace(raw)))

	switch policy {
	case PolicyManaged, PolicyReplace, PolicyPreserve:
		return policy, nil
	default:
		return "", fmt.Errorf("unsupported regeneration policy %q", raw)
	}
}

// Action describes one planned filesystem outcome.
type Action string

const (
	// ActionCreate writes a new managed file.
	ActionCreate Action = "create"
	// ActionUpdate replaces an existing managed file.
	ActionUpdate Action = "update"
	// ActionDelete removes a stale managed file.
	ActionDelete Action = "delete"
	// ActionPreserve keeps an existing file without claiming ownership.
	ActionPreserve Action = "preserve"
	// ActionConflict rejects a path that cannot be changed safely.
	ActionConflict Action = "conflict"
	// ActionUnchanged records a managed file with identical contents and mode.
	ActionUnchanged Action = "unchanged"
)

// Change describes one planned path outcome.
type Change struct {
	Path   string
	Action Action
	Reason string
}

// Plan is a deterministic description of a regeneration transaction.
type Plan struct {
	Changes []Change
}

// HasConflicts reports whether applying the plan would overwrite user work
// under the selected policy.
func (plan Plan) HasConflicts() bool {
	for _, change := range plan.Changes {
		if change.Action == ActionConflict {
			return true
		}
	}

	return false
}

// Count returns the number of changes with the requested action.
func (plan Plan) Count(action Action) int {
	count := 0

	for _, change := range plan.Changes {
		if change.Action == action {
			count++
		}
	}

	return count
}

// String renders a stable human-readable plan.
func (plan Plan) String() string {
	var output strings.Builder

	for _, change := range plan.Changes {
		output.WriteString(strings.ToUpper(string(change.Action)))
		_ = output.WriteByte('\t')
		output.WriteString(change.Path)

		if change.Reason != "" {
			_ = output.WriteByte('\t')
			output.WriteString(change.Reason)
		}

		_ = output.WriteByte('\n')
	}

	_, _ = fmt.Fprintf(
		&output,
		"summary: create=%d update=%d delete=%d preserve=%d conflict=%d unchanged=%d",
		plan.Count(ActionCreate),
		plan.Count(ActionUpdate),
		plan.Count(ActionDelete),
		plan.Count(ActionPreserve),
		plan.Count(ActionConflict),
		plan.Count(ActionUnchanged),
	)

	return output.String()
}

// Request contains all inputs for a generation comparison and transaction.
type Request struct {
	StagedProject    string
	TargetProject    string
	GeneratorVersion string
	Profile          string
	Services         []string
	Policy           Policy
	DryRun           bool
}

// Manifest records the generator inputs and files owned after a successful run.
type Manifest struct {
	SchemaVersion    int            `yaml:"schemaVersion"`
	GeneratorVersion string         `yaml:"generatorVersion"`
	Profile          string         `yaml:"profile"`
	Services         []string       `yaml:"services"`
	Files            []ManifestFile `yaml:"files"`
}

// ManifestFile records one owned path and the digest last written by goboot.
type ManifestFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
	Mode   uint32 `yaml:"mode"`
}

// ConflictError reports that a managed transaction would overwrite user work.
type ConflictError struct {
	Count int
}

// Error implements error.
func (conflict ConflictError) Error() string {
	return fmt.Sprintf("regeneration plan contains %d conflict(s)", conflict.Count)
}

// IsConflict reports whether err represents regeneration conflicts.
func IsConflict(err error) bool {
	var conflict ConflictError

	return errors.As(err, &conflict)
}

func sortedServices(services []string) []string {
	sorted := append([]string(nil), services...)
	sort.Strings(sorted)

	return sorted
}
