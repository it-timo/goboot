package regeneration

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fileState struct {
	digest string
	mode   fs.FileMode
}

type preparedPlan struct {
	plan        Plan
	manifestRaw []byte
}

// Apply plans a staged project and transactionally updates the target unless
// DryRun is set. The returned plan is always safe to display to users.
func Apply(request Request) (Plan, error) {
	prepared, err := prepare(request)
	if err != nil {
		return Plan{}, err
	}

	if prepared.plan.HasConflicts() {
		return prepared.plan, ConflictError{Count: prepared.plan.Count(ActionConflict)}
	}

	if request.DryRun {
		return prepared.plan, nil
	}

	err = applyTransaction(request, prepared)
	if err != nil {
		return prepared.plan, err
	}

	return prepared.plan, nil
}

func prepare(request Request) (preparedPlan, error) {
	policy, err := ParsePolicy(string(request.Policy))
	if err != nil {
		return preparedPlan{}, err
	}
	request.Policy = policy

	err = validateProjectPaths(request.StagedProject, request.TargetProject)
	if err != nil {
		return preparedPlan{}, err
	}

	desired, err := scanProject(request.StagedProject)
	if err != nil {
		return preparedPlan{}, fmt.Errorf("failed to scan staged project: %w", err)
	}

	current, err := scanOptionalProject(request.TargetProject)
	if err != nil {
		return preparedPlan{}, fmt.Errorf("failed to scan target project: %w", err)
	}

	previousManifest, err := readManifest(request.TargetProject)
	if err != nil {
		return preparedPlan{}, err
	}

	plan, ownedFiles := buildPlan(desired, current, manifestFileMap(previousManifest), policy)
	typeConflicts, err := pathTypeConflicts(request.StagedProject, request.TargetProject)
	if err != nil {
		return preparedPlan{}, err
	}

	plan = mergeTypeConflicts(plan, typeConflicts)

	manifest := Manifest{
		SchemaVersion:    ManifestSchemaVersion,
		GeneratorVersion: request.GeneratorVersion,
		Profile:          request.Profile,
		Services:         sortedServices(request.Services),
		Files:            ownedFiles,
	}

	manifestRaw, err := encodeManifest(manifest)
	if err != nil {
		return preparedPlan{}, err
	}

	plan = addManifestChange(plan, request.TargetProject, manifestRaw)

	return preparedPlan{
		plan:        plan,
		manifestRaw: manifestRaw,
	}, nil
}

func mergeTypeConflicts(plan Plan, conflicts []Change) Plan {
	if len(conflicts) == 0 {
		return plan
	}

	conflictPaths := make(map[string]struct{}, len(conflicts))
	for _, conflict := range conflicts {
		conflictPaths[conflict.Path] = struct{}{}
	}

	filtered := make([]Change, 0, len(plan.Changes)+len(conflicts))
	for _, change := range plan.Changes {
		if _, found := conflictPaths[change.Path]; !found {
			filtered = append(filtered, change)
		}
	}

	plan.Changes = append(filtered, conflicts...)

	return plan
}

func pathTypeConflicts(stagedProject, targetProject string) ([]Change, error) {
	desiredTypes, err := scanPathTypes(stagedProject)
	if err != nil {
		return nil, err
	}

	currentTypes, err := scanOptionalPathTypes(targetProject)
	if err != nil {
		return nil, err
	}

	var conflicts []Change

	for path, desiredDirectory := range desiredTypes {
		currentDirectory, found := currentTypes[path]
		if !found || currentDirectory == desiredDirectory {
			continue
		}

		conflicts = append(conflicts, Change{
			Path:   path,
			Action: ActionConflict,
			Reason: "generated and existing path types differ",
		})
	}

	return conflicts, nil
}

func scanOptionalPathTypes(projectPath string) (map[string]bool, error) {
	_, err := os.Lstat(projectPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]bool{}, nil
		}

		return nil, err
	}

	return scanPathTypes(projectPath)
}

func scanPathTypes(projectPath string) (map[string]bool, error) {
	pathTypes := make(map[string]bool)

	err := filepath.WalkDir(projectPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relativePath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		if relativePath == "." || relativePath == ManifestFileName {
			return nil
		}

		pathTypes[relativePath] = entry.IsDir()

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pathTypes, nil
}

func validateProjectPaths(stagedProject, targetProject string) error {
	stagedAbsolute, err := filepath.Abs(filepath.Clean(stagedProject))
	if err != nil {
		return fmt.Errorf("failed to resolve staged project path: %w", err)
	}

	targetAbsolute, err := filepath.Abs(filepath.Clean(targetProject))
	if err != nil {
		return fmt.Errorf("failed to resolve target project path: %w", err)
	}

	if projectPathsOverlap(stagedAbsolute, targetAbsolute) {
		return errors.New("staged and target project paths must not overlap")
	}

	return nil
}

func projectPathsOverlap(firstPath, secondPath string) bool {
	if firstPath == secondPath {
		return true
	}

	if pathContains(firstPath, secondPath) {
		return true
	}

	return pathContains(secondPath, firstPath)
}

func pathContains(parentPath, childPath string) bool {
	relativePath, err := filepath.Rel(parentPath, childPath)
	if err != nil {
		return false
	}

	return relativePath != ".." && !strings.HasPrefix(relativePath, ".."+string(os.PathSeparator))
}

func scanOptionalProject(projectPath string) (map[string]fileState, error) {
	info, err := os.Lstat(projectPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]fileState{}, nil
		}

		return nil, err
	}

	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return nil, fmt.Errorf("target project must be a real directory: %q", projectPath)
	}

	return scanProject(projectPath)
}

func scanProject(projectPath string) (map[string]fileState, error) {
	files := make(map[string]fileState)

	err := filepath.WalkDir(projectPath, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symbolic links are not supported in regeneration trees: %q", path)
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return fmt.Errorf("non-regular files are not supported in regeneration trees: %q", path)
		}

		relativePath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		if relativePath == ManifestFileName {
			return nil
		}

		digest, err := fileDigest(path)
		if err != nil {
			return err
		}

		files[relativePath] = fileState{digest: digest, mode: info.Mode().Perm()}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304 -- paths are discovered under validated project roots.
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, copyErr := io.Copy(hash, file)
	closeErr := file.Close()

	if copyErr != nil {
		return "", copyErr
	}

	if closeErr != nil {
		return "", closeErr
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func buildPlan(
	desired map[string]fileState,
	current map[string]fileState,
	previous map[string]ManifestFile,
	policy Policy,
) (Plan, []ManifestFile) {
	paths := combinedPaths(desired, previous)
	changes := make([]Change, 0, len(paths))
	ownedFiles := make([]ManifestFile, 0, len(desired))

	for _, path := range paths {
		change, owned := planPath(path, desired, current, previous, policy)
		if change != nil {
			changes = append(changes, *change)
		}

		if owned {
			ownedFiles = append(ownedFiles, ManifestFile{
				Path:   path,
				SHA256: desired[path].digest,
				Mode:   uint32(desired[path].mode.Perm()),
			})
		}
	}

	return Plan{Changes: changes}, ownedFiles
}

func combinedPaths(desired map[string]fileState, previous map[string]ManifestFile) []string {
	pathSet := make(map[string]struct{}, len(desired)+len(previous))
	for path := range desired {
		pathSet[path] = struct{}{}
	}

	for path := range previous {
		pathSet[path] = struct{}{}
	}

	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}

	sort.Strings(paths)

	return paths
}

func planPath(
	path string,
	desired map[string]fileState,
	current map[string]fileState,
	previous map[string]ManifestFile,
	policy Policy,
) (*Change, bool) {
	desiredFile, desiredFound := desired[path]
	currentFile, currentFound := current[path]
	previousFile, previouslyOwned := previous[path]

	if desiredFound {
		return planDesiredPath(path, desiredFile, currentFile, currentFound, previousFile, previouslyOwned, policy)
	}

	return planStalePath(path, currentFile, currentFound, previousFile, previouslyOwned, policy)
}

func planDesiredPath(
	path string,
	desired fileState,
	current fileState,
	currentFound bool,
	previous ManifestFile,
	previouslyOwned bool,
	policy Policy,
) (*Change, bool) {
	if !currentFound {
		return &Change{Path: path, Action: ActionCreate, Reason: "generated file is new"}, true
	}

	if previouslyOwned && desired.digest == current.digest && desired.mode == current.mode {
		return &Change{Path: path, Action: ActionUnchanged, Reason: "content and mode already match"}, true
	}

	if previouslyOwned && current.digest == previous.SHA256 && uint32(current.mode.Perm()) == previous.Mode {
		return &Change{Path: path, Action: ActionUpdate, Reason: "previously generated file is unchanged by the user"}, true
	}

	reason := "existing file is not owned by goboot"
	if previouslyOwned {
		reason = "previously generated file was modified after generation"
	}

	switch policy {
	case PolicyReplace:
		return &Change{Path: path, Action: ActionUpdate, Reason: reason + "; replace policy selected"}, true
	case PolicyPreserve:
		return &Change{Path: path, Action: ActionPreserve, Reason: reason + "; preserve policy selected"}, false
	default:
		return &Change{Path: path, Action: ActionConflict, Reason: reason}, false
	}
}

func planStalePath(
	path string,
	current fileState,
	currentFound bool,
	previous ManifestFile,
	previouslyOwned bool,
	policy Policy,
) (*Change, bool) {
	if !previouslyOwned || !currentFound {
		return nil, false
	}

	if policy == PolicyPreserve {
		return &Change{Path: path, Action: ActionPreserve, Reason: "stale generated file preserved by policy"}, false
	}

	currentUnchanged := current.digest == previous.SHA256 && uint32(current.mode.Perm()) == previous.Mode
	if currentUnchanged || policy == PolicyReplace {
		return &Change{Path: path, Action: ActionDelete, Reason: "file is no longer generated"}, false
	}

	return &Change{Path: path, Action: ActionConflict, Reason: "stale generated file was modified after generation"}, false
}

func addManifestChange(plan Plan, targetProject string, manifestRaw []byte) Plan {
	action := ActionCreate
	reason := "ownership manifest is new"

	current, err := os.ReadFile(filepath.Join(targetProject, ManifestFileName)) // #nosec G304 -- target path is validated.
	if err == nil {
		if bytes.Equal(current, manifestRaw) {
			action = ActionUnchanged
			reason = "ownership manifest already matches"
		} else {
			action = ActionUpdate
			reason = "ownership metadata changed"
		}
	}

	plan.Changes = append(plan.Changes, Change{Path: ManifestFileName, Action: action, Reason: reason})
	sort.Slice(plan.Changes, func(firstIndex, secondIndex int) bool {
		return plan.Changes[firstIndex].Path < plan.Changes[secondIndex].Path
	})

	return plan
}
