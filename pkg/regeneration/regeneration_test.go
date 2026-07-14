package regeneration_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/it-timo/goboot/pkg/regeneration"
)

func writeFile(root, relativePath, content string, mode os.FileMode) {
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	Expect(os.MkdirAll(filepath.Dir(path), 0o755)).To(Succeed())
	Expect(os.WriteFile(path, []byte(content), mode)).To(Succeed())
}

func readFile(root, relativePath string) string {
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relativePath)))
	Expect(err).NotTo(HaveOccurred())

	return string(content)
}

func readManifest(projectPath string) regeneration.Manifest {
	content, err := os.ReadFile(filepath.Join(projectPath, regeneration.ManifestFileName))
	Expect(err).NotTo(HaveOccurred())

	var manifest regeneration.Manifest
	Expect(yaml.Unmarshal(content, &manifest)).To(Succeed())

	return manifest
}

var _ = Describe("Regeneration transactions", func() {
	var (
		tempDir       string
		stagedProject string
		targetProject string
		request       regeneration.Request
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		stagedProject = filepath.Join(tempDir, "staged", "Example")
		targetProject = filepath.Join(tempDir, "target", "Example")

		Expect(os.MkdirAll(stagedProject, 0o755)).To(Succeed())

		request = regeneration.Request{
			StagedProject:    stagedProject,
			TargetProject:    targetProject,
			GeneratorVersion: "v0.7.0",
			Profile:          "standard",
			Services:         []string{"base_test", "base_project"},
			Policy:           regeneration.PolicyManaged,
		}
	})

	It("creates a project and deterministic ownership manifest", func() {
		writeFile(stagedProject, "README.md", "generated\n", 0o644)
		writeFile(stagedProject, "scripts/check.sh", "#!/bin/sh\n", 0o755)

		stagedInfo, err := os.Stat(filepath.Join(stagedProject, "scripts", "check.sh"))
		Expect(err).NotTo(HaveOccurred())

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionCreate)).To(Equal(3))
		Expect(readFile(targetProject, "README.md")).To(Equal("generated\n"))

		manifest := readManifest(targetProject)
		Expect(manifest.SchemaVersion).To(Equal(regeneration.ManifestSchemaVersion))
		Expect(manifest.GeneratorVersion).To(Equal("v0.7.0"))
		Expect(manifest.Services).To(Equal([]string{"base_project", "base_test"}))
		Expect(manifest.Files).To(HaveLen(2))
		Expect(manifest.Files[0].Path).To(Equal("README.md"))
		Expect(manifest.Files[1].Path).To(Equal(filepath.Join("scripts", "check.sh")))
		Expect(manifest.Files[1].Mode).To(Equal(uint32(stagedInfo.Mode().Perm())))
	})

	It("produces a dry-run plan without creating the target", func() {
		writeFile(stagedProject, "README.md", "planned\n", 0o644)

		request.DryRun = true

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionCreate)).To(Equal(2))
		Expect(targetProject).NotTo(BeAnExistingFile())
	})

	It("does not claim an identical file without ownership evidence", func() {
		writeFile(stagedProject, "README.md", "same bytes\n", 0o644)
		writeFile(targetProject, "README.md", "same bytes\n", 0o644)

		plan, err := regeneration.Apply(request)
		Expect(regeneration.IsConflict(err)).To(BeTrue())
		Expect(plan.Count(regeneration.ActionConflict)).To(Equal(1))
		Expect(filepath.Join(targetProject, regeneration.ManifestFileName)).NotTo(BeAnExistingFile())
	})

	It("updates an unchanged file owned by the previous manifest", func() {
		writeFile(stagedProject, "README.md", "first\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())

		writeFile(stagedProject, "README.md", "second\n", 0o644)

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionUpdate)).To(Equal(2))
		Expect(readFile(targetProject, "README.md")).To(Equal("second\n"))
	})

	It("rejects user modifications under the managed policy", func() {
		writeFile(stagedProject, "README.md", "first\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())

		writeFile(targetProject, "README.md", "user edit\n", 0o644)
		writeFile(stagedProject, "README.md", "second\n", 0o644)

		plan, err := regeneration.Apply(request)
		Expect(regeneration.IsConflict(err)).To(BeTrue())
		Expect(plan.Count(regeneration.ActionConflict)).To(Equal(1))
		Expect(readFile(targetProject, "README.md")).To(Equal("user edit\n"))
	})

	It("detects user changes to an owned file mode", func() {
		if runtime.GOOS == "windows" {
			Skip("Windows does not preserve POSIX permission bits")
		}

		writeFile(stagedProject, "scripts/check.sh", "#!/bin/sh\n", 0o755)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())

		Expect(os.Chmod(filepath.Join(targetProject, "scripts", "check.sh"), 0o644)).To(Succeed())

		_, err = regeneration.Apply(request)
		Expect(regeneration.IsConflict(err)).To(BeTrue())
	})

	It("overwrites modified files only under the replace policy", func() {
		writeFile(stagedProject, "README.md", "first\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())

		writeFile(targetProject, "README.md", "user edit\n", 0o644)
		writeFile(stagedProject, "README.md", "replacement\n", 0o644)

		request.Policy = regeneration.PolicyReplace

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionUpdate)).To(Equal(2))
		Expect(readFile(targetProject, "README.md")).To(Equal("replacement\n"))
	})

	It("preserves collisions and relinquishes their ownership", func() {
		writeFile(stagedProject, "README.md", "first\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())

		writeFile(targetProject, "README.md", "user edit\n", 0o644)
		writeFile(stagedProject, "README.md", "second\n", 0o644)

		request.Policy = regeneration.PolicyPreserve

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionPreserve)).To(Equal(1))
		Expect(readFile(targetProject, "README.md")).To(Equal("user edit\n"))
		Expect(readManifest(targetProject).Files).To(BeEmpty())
	})

	It("deletes unchanged stale generated files", func() {
		writeFile(stagedProject, "keep.txt", "keep\n", 0o644)
		writeFile(stagedProject, "stale.txt", "stale\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Remove(filepath.Join(stagedProject, "stale.txt"))).To(Succeed())

		plan, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Count(regeneration.ActionDelete)).To(Equal(1))
		Expect(filepath.Join(targetProject, "stale.txt")).NotTo(BeAnExistingFile())
	})

	It("rejects deletion of a modified stale file", func() {
		writeFile(stagedProject, "stale.txt", "stale\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		writeFile(targetProject, "stale.txt", "user edit\n", 0o644)
		Expect(os.Remove(filepath.Join(stagedProject, "stale.txt"))).To(Succeed())

		plan, err := regeneration.Apply(request)
		Expect(regeneration.IsConflict(err)).To(BeTrue())
		Expect(plan.Count(regeneration.ActionConflict)).To(Equal(1))
		Expect(readFile(targetProject, "stale.txt")).To(Equal("user edit\n"))
	})

	It("preserves unowned files across a successful transaction", func() {
		Expect(os.MkdirAll(targetProject, 0o755)).To(Succeed())
		writeFile(targetProject, "notes.txt", "user owned\n", 0o600)
		writeFile(stagedProject, "README.md", "generated\n", 0o644)

		originalInfo, err := os.Stat(filepath.Join(targetProject, "notes.txt"))
		Expect(err).NotTo(HaveOccurred())

		_, err = regeneration.Apply(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(readFile(targetProject, "notes.txt")).To(Equal("user owned\n"))

		info, err := os.Stat(filepath.Join(targetProject, "notes.txt"))
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode().Perm()).To(Equal(originalInfo.Mode().Perm()))
	})

	It("rejects a malformed ownership manifest without changing the project", func() {
		Expect(os.MkdirAll(targetProject, 0o755)).To(Succeed())
		writeFile(targetProject, regeneration.ManifestFileName, "schemaVersion: invalid\n", 0o644)
		writeFile(targetProject, "notes.txt", "untouched\n", 0o644)
		writeFile(stagedProject, "README.md", "generated\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(err).To(HaveOccurred())
		Expect(readFile(targetProject, "notes.txt")).To(Equal("untouched\n"))
		Expect(filepath.Join(targetProject, "README.md")).NotTo(BeAnExistingFile())
	})

	It("rejects symbolic links in staged or target trees", func() {
		if runtime.GOOS == "windows" {
			Skip("symbolic-link permissions differ on Windows")
		}

		writeFile(stagedProject, "source.txt", "source\n", 0o644)
		Expect(os.Symlink("source.txt", filepath.Join(stagedProject, "link.txt"))).To(Succeed())

		_, err := regeneration.Apply(request)
		Expect(err).To(MatchError(ContainSubstring("symbolic links are not supported")))
		Expect(targetProject).NotTo(BeAnExistingFile())
	})

	It("leaves the target untouched for path type conflicts", func() {
		Expect(os.MkdirAll(filepath.Join(targetProject, "collision"), 0o755)).To(Succeed())
		writeFile(targetProject, "notes.txt", "untouched\n", 0o644)
		writeFile(stagedProject, "collision", "generated file\n", 0o644)

		_, err := regeneration.Apply(request)
		Expect(regeneration.IsConflict(err)).To(BeTrue())
		Expect(readFile(targetProject, "notes.txt")).To(Equal("untouched\n"))
		Expect(filepath.Join(targetProject, regeneration.ManifestFileName)).NotTo(BeAnExistingFile())
	})

	It("recognizes typed conflict errors through wrapping", func() {
		wrapped := errors.New("ordinary error")
		Expect(regeneration.IsConflict(wrapped)).To(BeFalse())
	})
})
