package regeneration

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

func applyTransaction(request Request, prepared preparedPlan) error {
	targetParent := filepath.Dir(request.TargetProject)

	err := os.MkdirAll(targetParent, goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to create target parent: %w", err)
	}

	candidatePath, err := os.MkdirTemp(targetParent, ".goboot-candidate-*")
	if err != nil {
		return fmt.Errorf("failed to create transaction candidate: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(candidatePath)
		}
	}()

	err = prepareCandidate(request, prepared, candidatePath)
	if err != nil {
		return err
	}

	err = swapProject(candidatePath, request.TargetProject)
	if err != nil {
		return err
	}

	committed = true

	return nil
}

func prepareCandidate(request Request, prepared preparedPlan, candidatePath string) error {
	err := setCandidateMode(candidatePath, request.StagedProject, request.TargetProject)
	if err != nil {
		return err
	}

	if projectDirectoryExists(request.TargetProject) {
		err = copyTree(request.TargetProject, candidatePath)
		if err != nil {
			return fmt.Errorf("failed to copy current project into transaction: %w", err)
		}
	}

	err = applyChanges(candidatePath, request.StagedProject, prepared.plan)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(candidatePath, ManifestFileName), prepared.manifestRaw, goboottypes.FilePerm)
	if err != nil {
		return fmt.Errorf("failed to write ownership manifest: %w", err)
	}

	err = os.Chmod(filepath.Join(candidatePath, ManifestFileName), goboottypes.FilePerm)
	if err != nil {
		return fmt.Errorf("failed to set ownership manifest mode: %w", err)
	}

	return nil
}

func setCandidateMode(candidatePath, stagedProject, targetProject string) error {
	modeSource := stagedProject
	if projectDirectoryExists(targetProject) {
		modeSource = targetProject
	}

	info, err := os.Stat(modeSource)
	if err != nil {
		return fmt.Errorf("failed to inspect project root mode: %w", err)
	}

	err = os.Chmod(candidatePath, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to set transaction root mode: %w", err)
	}

	return nil
}

func projectDirectoryExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func applyChanges(candidatePath, stagedProject string, plan Plan) error {
	for _, change := range plan.Changes {
		if change.Path == ManifestFileName {
			continue
		}

		switch change.Action {
		case ActionCreate, ActionUpdate:
			err := copyFile(
				filepath.Join(stagedProject, change.Path),
				filepath.Join(candidatePath, change.Path),
			)
			if err != nil {
				return fmt.Errorf("failed to apply %s for %q: %w", change.Action, change.Path, err)
			}
		case ActionDelete:
			err := os.Remove(filepath.Join(candidatePath, change.Path))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to delete stale file %q: %w", change.Path, err)
			}
		case ActionPreserve, ActionUnchanged:
			continue
		case ActionConflict:
			return errors.New("refusing to apply a plan containing conflicts")
		}
	}

	return nil
}

func copyTree(sourceRoot, targetRoot string) error {
	err := filepath.WalkDir(sourceRoot, func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symbolic links are not supported in regeneration trees: %q", sourcePath)
		}

		relativePath, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return fmt.Errorf("failed to resolve copied project path: %w", err)
		}

		if relativePath == "." {
			return nil
		}

		targetPath := filepath.Join(targetRoot, relativePath)

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to inspect copied project entry: %w", err)
		}

		if entry.IsDir() {
			err = os.MkdirAll(targetPath, info.Mode().Perm())
			if err != nil {
				return fmt.Errorf("failed to create copied project directory: %w", err)
			}

			return nil
		}

		if !info.Mode().IsRegular() {
			return fmt.Errorf("non-regular files are not supported in regeneration trees: %q", sourcePath)
		}

		return copyFile(sourcePath, targetPath)
	})
	if err != nil {
		return fmt.Errorf("failed to walk copied project tree: %w", err)
	}

	return nil
}

func copyFile(sourcePath, targetPath string) error {
	err := os.MkdirAll(filepath.Dir(targetPath), goboottypes.DirPerm)
	if err != nil {
		return fmt.Errorf("failed to create copied file parent: %w", err)
	}

	source, err := os.Open(sourcePath) // #nosec G304 -- paths are confined to validated transaction trees.
	if err != nil {
		return fmt.Errorf("failed to open copied source file: %w", err)
	}

	info, err := source.Stat()
	if err != nil {
		_ = source.Close()

		return fmt.Errorf("failed to inspect copied source file: %w", err)
	}

	// #nosec G304 -- target is confined to the transaction tree.
	target, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		_ = source.Close()

		return fmt.Errorf("failed to open copied target file: %w", err)
	}

	_, copyErr := io.Copy(target, source)
	sourceCloseErr := source.Close()
	targetCloseErr := target.Close()

	if copyErr != nil {
		return fmt.Errorf("failed to copy file contents: %w", copyErr)
	}

	if sourceCloseErr != nil {
		return fmt.Errorf("failed to close copied source file: %w", sourceCloseErr)
	}

	if targetCloseErr != nil {
		return fmt.Errorf("failed to close copied target file: %w", targetCloseErr)
	}

	err = os.Chmod(targetPath, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to set copied target mode: %w", err)
	}

	return nil
}

func swapProject(candidatePath, targetProject string) error {
	if !projectDirectoryExists(targetProject) {
		err := os.Rename(candidatePath, targetProject)
		if err != nil {
			return fmt.Errorf("failed to commit new project transaction: %w", err)
		}

		return nil
	}

	backupPath, err := reserveBackupPath(filepath.Dir(targetProject))
	if err != nil {
		return err
	}

	err = os.Rename(targetProject, backupPath)
	if err != nil {
		return fmt.Errorf("failed to move current project to transaction backup: %w", err)
	}

	err = os.Rename(candidatePath, targetProject)
	if err != nil {
		rollbackErr := os.Rename(backupPath, targetProject)
		if rollbackErr != nil {
			return fmt.Errorf("failed to commit transaction and rollback: %w", errors.Join(err, rollbackErr))
		}

		return fmt.Errorf("failed to commit transaction; previous project restored: %w", err)
	}

	err = os.RemoveAll(backupPath)
	if err != nil {
		return fmt.Errorf("transaction committed but backup cleanup failed: %w", err)
	}

	return nil
}

func reserveBackupPath(parent string) (string, error) {
	backupPath, err := os.MkdirTemp(parent, ".goboot-backup-*")
	if err != nil {
		return "", fmt.Errorf("failed to reserve transaction backup: %w", err)
	}

	err = os.Remove(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to prepare transaction backup: %w", err)
	}

	return backupPath, nil
}
