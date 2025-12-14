package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// repoRoot returns the repository root based on this test file location.
func repoRoot(t GinkgoTInterface) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to resolve current file location")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// writeConfig helper ensures test files are written with the expected permissions.
func writeConfig(path, content string) {
	err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644)
	Expect(err).NotTo(HaveOccurred())
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	return string(data)
}

var _ = Describe("End-to-end goboot runs", func() {
	It("scaffolds a full project with all services enabled (ginkgo style)", func() {
		defer withFakeGo()()
		tempDir := GinkgoT().TempDir()
		projectName := "E2EGinkgo"
		repoURL := "github.com/example/e2e-ginkgo"
		targetDir := filepath.Join(tempDir, "target")
		projectRoot := filepath.Join(targetDir, projectName)
		root := repoRoot(GinkgoT())

		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		writeConfig(baseProjectCfg, fmt.Sprintf(`
sourcePath: %s
usedGoVersion: "1.22.5"
usedNodeVersion: "20.0.0"
releaseCurrentWindow: "Q1 2026"
releaseUpcomingWindow: "Q3 2026"
releaseLongTerm: "2029"
author: "E2E Author"
gitProvider: "github"
gitUser: "example"
`, filepath.Join(root, "templates", "project_base")))

		baseLintCfg := filepath.Join(tempDir, "base_lint.yml")
		writeConfig(baseLintCfg, fmt.Sprintf(`
sourcePath: %s
linters:
  golang:
    cmd: |-
      {{DOCKER_RUN}} golangci/golangci-lint:v2.7.1 golangci-lint run ./...
    enabled: true
  yaml:
    cmd: |-
      {{DOCKER_RUN}} pipelinecomponents/yamllint:0.35.9 yamllint .
    enabled: true
  markdown:
    cmd: |-
      {{DOCKER_RUN}} ghcr.io/igorshubovych/markdownlint-cli:v0.46.0 markdownlint "**/*.md"
    enabled: true
  shellcheck:
    cmd: |-
      {{DOCKER_RUN}} koalaman/shellcheck:v0.11.0 -x {{SH_FILES}}
    enabled: true
  shfmt:
    cmd: |-
      {{DOCKER_RUN}} mvdan/shfmt:v3.12.0 -d -i 2 -ci {{SH_FILES}}
    enabled: true
allowedPackages:
  - gopkg.in/yaml.v3
`, filepath.Join(root, "templates", "lint_base")))

		baseTestCfg := filepath.Join(tempDir, "base_test.yml")
		writeConfig(baseTestCfg, fmt.Sprintf(`
sourcePath: %s
useStyle: "ginkgo"
testCmd: |-
  go test ./...
`, filepath.Join(root, "templates", "test_base")))

		baseLocalCfg := filepath.Join(tempDir, "base_local.yml")
		writeConfig(baseLocalCfg, fmt.Sprintf(`
sourcePath: %s
fileList:
  - make
  - task
  - script
  - commit
`, filepath.Join(root, "templates", "local_base")))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		writeConfig(gobootCfg, fmt.Sprintf(`
projectName: %s
repoUrl: %s
targetPath: %s
services:
  - id: base_project
    confPath: %s
    enabled: true
  - id: base_lint
    confPath: %s
    enabled: true
  - id: base_test
    confPath: %s
    enabled: true
  - id: base_local
    confPath: %s
    enabled: true
`, projectName, repoURL, targetDir, baseProjectCfg, baseLintCfg, baseTestCfg, baseLocalCfg))

		Expect(run([]string{"--config", gobootCfg})).To(Succeed())

		Expect(projectRoot).To(BeADirectory())
		Expect(filepath.Join(projectRoot, ".golangci.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".yamllint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".markdownlint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".shellcheckrc")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "Taskfile.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "lint.sh")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "test.sh")).To(BeAnExistingFile())

		goMod := readFile(filepath.Join(projectRoot, "go.mod"))
		Expect(goMod).To(ContainSubstring("module " + repoURL))
		Expect(goMod).NotTo(ContainSubstring("{{"))

		makefile := readFile(filepath.Join(projectRoot, "Makefile"))
		Expect(makefile).To(ContainSubstring("Makefile — Developer Targets"))
		Expect(makefile).To(ContainSubstring("PROJECT := E2EGinkgo"))
		Expect(makefile).To(ContainSubstring("golangci/golangci-lint"))
		Expect(makefile).NotTo(ContainSubstring("{{"))

		lintScript := readFile(filepath.Join(projectRoot, "scripts", "lint.sh"))
		Expect(lintScript).To(ContainSubstring("E2EGinkgo"))
		Expect(lintScript).To(ContainSubstring("golangci/golangci-lint"))
		Expect(lintScript).NotTo(ContainSubstring("{{"))

		testScript := readFile(filepath.Join(projectRoot, "scripts", "test.sh"))
		Expect(testScript).To(ContainSubstring("go test ./..."))
		Expect(testScript).NotTo(ContainSubstring("{{"))

		ginkgoSuite := readFile(filepath.Join(projectRoot, "pkg", "e2eginkgo", "e2eginkgo_suite_test.go"))
		Expect(ginkgoSuite).To(ContainSubstring("RunSpecs"))
		Expect(ginkgoSuite).To(ContainSubstring("E2EGinkgo Suite"))
	})

	It("supports go-style tests and selectively enabled linters", func() {
		defer withFakeGo()()
		tempDir := GinkgoT().TempDir()
		projectName := "E2EGoStyle"
		repoURL := "github.com/example/e2e-gostyle"
		targetDir := filepath.Join(tempDir, "out")
		projectRoot := filepath.Join(targetDir, projectName)
		root := repoRoot(GinkgoT())

		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		writeConfig(baseProjectCfg, fmt.Sprintf(`
sourcePath: %s
usedGoVersion: "1.21.9"
usedNodeVersion: "18.0.0"
releaseCurrentWindow: "Q2 2026"
releaseUpcomingWindow: "Q4 2026"
releaseLongTerm: "2030"
author: "Go Style Author"
gitProvider: "gitlab"
gitUser: "go-style"
`, filepath.Join(root, "templates", "project_base")))

		baseLintCfg := filepath.Join(tempDir, "base_lint.yml")
		writeConfig(baseLintCfg, fmt.Sprintf(`
sourcePath: %s
linters:
  golang:
    cmd: echo go-lint
    enabled: true
  yaml:
    cmd: ""
    enabled: true
  markdown:
    cmd: echo markdown
    enabled: false
  shellcheck:
    cmd: echo shellcheck
    enabled: false
allowedPackages:
  - github.com/example/safe
`, filepath.Join(root, "templates", "lint_base")))

		baseTestCfg := filepath.Join(tempDir, "base_test.yml")
		writeConfig(baseTestCfg, fmt.Sprintf(`
sourcePath: %s
useStyle: "go"
testCmd: ""
`, filepath.Join(root, "templates", "test_base")))

		baseLocalCfg := filepath.Join(tempDir, "base_local.yml")
		writeConfig(baseLocalCfg, fmt.Sprintf(`
sourcePath: %s
fileList:
  - task
  - script
`, filepath.Join(root, "templates", "local_base")))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		writeConfig(gobootCfg, fmt.Sprintf(`
projectName: %s
repoUrl: %s
targetPath: %s
services:
  - id: base_project
    confPath: %s
    enabled: true
  - id: base_lint
    confPath: %s
    enabled: true
  - id: base_test
    confPath: %s
    enabled: true
  - id: base_local
    confPath: %s
    enabled: true
`, projectName, repoURL, targetDir, baseProjectCfg, baseLintCfg, baseTestCfg, baseLocalCfg))

		Expect(run([]string{"--config", gobootCfg})).To(Succeed())

		Expect(projectRoot).To(BeADirectory())
		Expect(filepath.Join(projectRoot, ".golangci.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".yamllint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".markdownlint.yml")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".shellcheckrc")).NotTo(BeAnExistingFile())

		// go style skips ginkgo suite generation.
		_, err := os.Stat(filepath.Join(projectRoot, "pkg", "e2egostyle", "e2egostyle_suite_test.go"))
		Expect(os.IsNotExist(err)).To(BeTrue())

		goStyleTest := readFile(filepath.Join(projectRoot, "pkg", "e2egostyle", "e2egostyle_test.go"))
		Expect(goStyleTest).To(ContainSubstring("package e2egostyle_test"))
		Expect(goStyleTest).To(ContainSubstring(repoURL))
		Expect(goStyleTest).NotTo(ContainSubstring("Describe"))

		lintScript := readFile(filepath.Join(projectRoot, "scripts", "lint.sh"))
		Expect(lintScript).To(ContainSubstring("go-lint"))
		Expect(lintScript).To(ContainSubstring("yamllint"))

		testScript := readFile(filepath.Join(projectRoot, "scripts", "test.sh"))
		Expect(testScript).To(ContainSubstring("go test -race -timeout=5m"))
		Expect(testScript).NotTo(ContainSubstring("{{"))
	})
})
