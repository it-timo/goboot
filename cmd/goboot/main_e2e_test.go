package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
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

func parseYAMLFile(path string) map[string]any {
	data, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	var doc map[string]any

	err = yaml.Unmarshal(data, &doc)
	Expect(err).NotTo(HaveOccurred(), "invalid YAML: %s", path)

	return doc
}

func assertGeneratedGitLabCIValid(projectRoot string) {
	rootDoc := parseYAMLFile(filepath.Join(projectRoot, ".gitlab-ci.yml"))
	Expect(rootDoc).To(HaveKey("stages"))
	Expect(rootDoc).To(HaveKey("include"))

	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "commands.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "versions.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "lint.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "test.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "build.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "container.yml"))
	parseYAMLFile(filepath.Join(projectRoot, ".gitlab/ci", "security.yml"))
}

func assertYAMLFilesValid(projectRoot string, relPaths []string) {
	for _, relPath := range relPaths {
		parseYAMLFile(filepath.Join(projectRoot, filepath.FromSlash(relPath)))
	}
}

func assertFilesExist(root string, relPaths []string) {
	for _, relPath := range relPaths {
		Expect(filepath.Join(root, filepath.FromSlash(relPath))).To(BeAnExistingFile())
	}
}

func assertFilesNotExist(root string, relPaths []string) {
	for _, relPath := range relPaths {
		Expect(filepath.Join(root, filepath.FromSlash(relPath))).NotTo(BeAnExistingFile())
	}
}

var _ = Describe("End-to-end goboot runs", func() {
	const (
		baseProjectFixtureFile = "cmd_goboot/goboot/service_execution_base_project.yml"
	)

	It("scaffolds a full project with all services enabled (ginkgo style)", func() {
		defer withFakeGo()()

		tempDir := GinkgoT().TempDir()
		projectName := "E2EGinkgo"
		repoURL := "github.com/example/e2e-ginkgo"
		gitProvider := "gitlab"
		targetDir := filepath.Join(tempDir, "target")
		projectRoot := filepath.Join(targetDir, projectName)
		root := repoRoot(GinkgoT())
		projectBaseTemplates := filepath.Join(root, "templates", "project_base")
		lintBaseTemplates := filepath.Join(root, "templates", "lint_base")
		testBaseTemplates := filepath.Join(root, "templates", "test_base")
		localBaseTemplates := filepath.Join(root, "templates", "local_base")
		ciBaseTemplates := filepath.Join(root, "templates", "ci_base")
		dockerBaseTemplates := filepath.Join(root, "templates", "docker_base")
		governanceBaseTemplates := filepath.Join(root, "templates", "governance_base")
		supplyChainBaseTemplates := filepath.Join(root, "templates", "supplychain_base")

		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		baseProjectContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_project.yml", map[string]string{
			"TEMPLATES_PROJECT_BASE": projectBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseProjectCfg, string(baseProjectContent))

		baseLintCfg := filepath.Join(tempDir, "base_lint.yml")
		baseLintContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_lint.yml", map[string]string{
			"TEMPLATES_LINT_BASE": lintBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLintCfg, string(baseLintContent))

		baseTestCfg := filepath.Join(tempDir, "base_test.yml")
		baseTestContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_test.yml", map[string]string{
			"TEMPLATES_TEST_BASE": testBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseTestCfg, string(baseTestContent))

		baseLocalCfg := filepath.Join(tempDir, "base_local.yml")
		baseLocalContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_local.yml", map[string]string{
			"TEMPLATES_LOCAL_BASE": localBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLocalCfg, string(baseLocalContent))

		baseCiCfg := filepath.Join(tempDir, "base_ci.yml")
		baseCiContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_ci.yml", map[string]string{
			"TEMPLATES_CI_BASE": ciBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseCiCfg, string(baseCiContent))

		baseLoggerCfg := filepath.Join(tempDir, "base_logger.yml")
		baseLoggerContent, err := loadTestFixture("cmd_goboot/base/ginkgo/base_logger.yml")
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLoggerCfg, string(baseLoggerContent))

		baseDockerCfg := filepath.Join(tempDir, "base_docker.yml")
		baseDockerContent, err := loadTestFixtureWithVars("cmd_goboot/base/ginkgo/base_docker.yml", map[string]string{
			"TEMPLATES_DOCKER_BASE": dockerBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseDockerCfg, string(baseDockerContent))

		baseGovernanceCfg := filepath.Join(tempDir, "base_governance.yml")
		baseGovernanceContent, err := loadTestFixtureWithVars("cmd_goboot/base/governance.yml", map[string]string{
			"TEMPLATES_GOVERNANCE_BASE": governanceBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseGovernanceCfg, string(baseGovernanceContent))

		baseSupplyChainCfg := filepath.Join(tempDir, "base_supplychain.yml")
		baseSupplyChainContent, err := loadTestFixtureWithVars("cmd_goboot/base/supplychain.yml", map[string]string{
			"TEMPLATES_SUPPLYCHAIN_BASE": supplyChainBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseSupplyChainCfg, string(baseSupplyChainContent))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		gobootContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/e2e.yml", map[string]string{
			fixtureProjectName:     projectName,
			"GIT_PROVIDER":         gitProvider,
			"REPO_URL":             repoURL,
			fixtureTargetDir:       targetDir,
			"BASE_PROJECT_CFG":     baseProjectCfg,
			"BASE_LINT_CFG":        baseLintCfg,
			"BASE_TEST_CFG":        baseTestCfg,
			"BASE_LOGGER_CFG":      baseLoggerCfg,
			"BASE_DOCKER_CFG":      baseDockerCfg,
			"BASE_GOVERNANCE_CFG":  baseGovernanceCfg,
			"BASE_SUPPLYCHAIN_CFG": baseSupplyChainCfg,
			"BASE_LOCAL_CFG":       baseLocalCfg,
			"BASE_CI_CFG":          baseCiCfg,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(gobootCfg, string(gobootContent))

		Expect(run([]string{argConfig, gobootCfg})).To(Succeed())

		Expect(projectRoot).To(BeADirectory())
		Expect(filepath.Join(projectRoot, fileGolangCI)).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".yamllint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".markdownlint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".shellcheckrc")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "Taskfile.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "Dockerfile")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "docker-compose.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".dockerignore")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "CODEOWNERS")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "CONTRIBUTING.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "SECURITY.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "SUPPLY_CHAIN.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".gitlab", "issue_templates", "Bug.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".gitlab", "issue_templates", "Feature.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".gitlab", "merge_request_templates", "Default.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "lint.sh")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "test.sh")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "docker.sh")).To(BeAnExistingFile())

		goMod := readFile(filepath.Join(projectRoot, "go.mod"))
		Expect(goMod).To(ContainSubstring("module " + repoURL))
		Expect(goMod).NotTo(ContainSubstring("{{"))

		makefile := readFile(filepath.Join(projectRoot, "Makefile"))
		Expect(makefile).To(ContainSubstring("Makefile — Developer Targets"))
		Expect(makefile).To(ContainSubstring("PROJECT := E2EGinkgo"))
		Expect(makefile).To(ContainSubstring("golangci/golangci-lint"))
		Expect(makefile).To(ContainSubstring("docker-build"))
		Expect(makefile).NotTo(ContainSubstring("{{"))

		dockerfile := readFile(filepath.Join(projectRoot, "Dockerfile"))
		Expect(dockerfile).To(ContainSubstring("FROM golang:1.26.5"))
		Expect(dockerfile).To(ContainSubstring("COPY go.mod ./"))
		Expect(dockerfile).NotTo(ContainSubstring("COPY go.mod go.sum"))
		Expect(dockerfile).To(ContainSubstring("./cmd/e2eginkgo"))
		Expect(dockerfile).NotTo(ContainSubstring("{{"))

		dockerignore := readFile(filepath.Join(projectRoot, ".dockerignore"))
		Expect(dockerignore).To(ContainSubstring(".env"))
		Expect(dockerignore).To(ContainSubstring(".codex"))
		Expect(dockerignore).To(ContainSubstring("coverage.*"))

		lintScript := readFile(filepath.Join(projectRoot, "scripts", "lint.sh"))
		Expect(lintScript).To(ContainSubstring("E2EGinkgo"))
		Expect(lintScript).To(ContainSubstring("golangci/golangci-lint"))
		Expect(lintScript).NotTo(ContainSubstring("{{"))

		testScript := readFile(filepath.Join(projectRoot, "scripts", "test.sh"))
		Expect(testScript).To(ContainSubstring("go test ./..."))
		Expect(testScript).NotTo(ContainSubstring("{{"))

		dockerScript := readFile(filepath.Join(projectRoot, "scripts", "docker.sh"))
		Expect(dockerScript).To(ContainSubstring("docker build -t e2eginkgo:local ."))
		Expect(dockerScript).To(ContainSubstring("docker compose up --build"))
		Expect(dockerScript).To(ContainSubstring(
			"docker compose config && docker build -t e2eginkgo:local . && docker run --rm e2eginkgo:local -h",
		))
		Expect(dockerScript).NotTo(ContainSubstring("{{"))

		mainGo := readFile(filepath.Join(projectRoot, "cmd", "e2eginkgo", "main.go"))
		Expect(mainGo).To(ContainSubstring("github.com/rs/zerolog"))

		ginkgoSuite := readFile(filepath.Join(projectRoot, "pkg", "e2eginkgo", "e2eginkgo_suite_test.go"))
		Expect(ginkgoSuite).To(ContainSubstring("RunSpecs"))
		Expect(ginkgoSuite).To(ContainSubstring("E2EGinkgo Suite"))

		assertGeneratedGitLabCIValid(projectRoot)
		assertYAMLFilesValid(projectRoot, []string{
			fileGolangCI,
			".yamllint.yml",
			".markdownlint.yml",
			".pre-commit-config.yaml",
			"Taskfile.yml",
			"configs/e2eginkgo.yml",
		})
	})

	It("supports go-style tests and selectively enabled linters", func() {
		defer withFakeGo()()

		tempDir := GinkgoT().TempDir()
		projectName := "E2EGoStyle"
		gitProvider := "gitlab"
		repoURL := "github.com/example/e2e-gostyle"
		targetDir := filepath.Join(tempDir, "out")
		projectRoot := filepath.Join(targetDir, projectName)
		root := repoRoot(GinkgoT())
		projectBaseTemplates := filepath.Join(root, "templates", "project_base")
		lintBaseTemplates := filepath.Join(root, "templates", "lint_base")
		testBaseTemplates := filepath.Join(root, "templates", "test_base")
		localBaseTemplates := filepath.Join(root, "templates", "local_base")
		ciBaseTemplates := filepath.Join(root, "templates", "ci_base")
		dockerBaseTemplates := filepath.Join(root, "templates", "docker_base")
		governanceBaseTemplates := filepath.Join(root, "templates", "governance_base")
		supplyChainBaseTemplates := filepath.Join(root, "templates", "supplychain_base")

		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		baseProjectContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_project.yml", map[string]string{
			"TEMPLATES_PROJECT_BASE": projectBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseProjectCfg, string(baseProjectContent))

		baseLintCfg := filepath.Join(tempDir, "base_lint.yml")
		baseLintContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_lint.yml", map[string]string{
			"TEMPLATES_LINT_BASE": lintBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLintCfg, string(baseLintContent))

		baseTestCfg := filepath.Join(tempDir, "base_test.yml")
		baseTestContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_test.yml", map[string]string{
			"TEMPLATES_TEST_BASE": testBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseTestCfg, string(baseTestContent))

		baseLocalCfg := filepath.Join(tempDir, "base_local.yml")
		baseLocalContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_local.yml", map[string]string{
			"TEMPLATES_LOCAL_BASE": localBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLocalCfg, string(baseLocalContent))

		baseCiCfg := filepath.Join(tempDir, "base_ci.yml")
		baseCiContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_ci.yml", map[string]string{
			"TEMPLATES_CI_BASE": ciBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseCiCfg, string(baseCiContent))

		baseLoggerCfg := filepath.Join(tempDir, "base_logger.yml")
		baseLoggerContent, err := loadTestFixture("cmd_goboot/base/go/base_logger.yml")
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseLoggerCfg, string(baseLoggerContent))

		baseDockerCfg := filepath.Join(tempDir, "base_docker.yml")
		baseDockerContent, err := loadTestFixtureWithVars("cmd_goboot/base/go/base_docker.yml", map[string]string{
			"TEMPLATES_DOCKER_BASE": dockerBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseDockerCfg, string(baseDockerContent))

		baseGovernanceCfg := filepath.Join(tempDir, "base_governance.yml")
		baseGovernanceContent, err := loadTestFixtureWithVars("cmd_goboot/base/governance.yml", map[string]string{
			"TEMPLATES_GOVERNANCE_BASE": governanceBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseGovernanceCfg, string(baseGovernanceContent))

		baseSupplyChainCfg := filepath.Join(tempDir, "base_supplychain.yml")
		baseSupplyChainContent, err := loadTestFixtureWithVars("cmd_goboot/base/supplychain.yml", map[string]string{
			"TEMPLATES_SUPPLYCHAIN_BASE": supplyChainBaseTemplates,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseSupplyChainCfg, string(baseSupplyChainContent))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		gobootContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/e2e.yml", map[string]string{
			fixtureProjectName:     projectName,
			"GIT_PROVIDER":         gitProvider,
			"REPO_URL":             repoURL,
			fixtureTargetDir:       targetDir,
			"BASE_PROJECT_CFG":     baseProjectCfg,
			"BASE_LINT_CFG":        baseLintCfg,
			"BASE_TEST_CFG":        baseTestCfg,
			"BASE_LOGGER_CFG":      baseLoggerCfg,
			"BASE_DOCKER_CFG":      baseDockerCfg,
			"BASE_GOVERNANCE_CFG":  baseGovernanceCfg,
			"BASE_SUPPLYCHAIN_CFG": baseSupplyChainCfg,
			"BASE_LOCAL_CFG":       baseLocalCfg,
			"BASE_CI_CFG":          baseCiCfg,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(gobootCfg, string(gobootContent))

		Expect(run([]string{argConfig, gobootCfg})).To(Succeed())

		Expect(projectRoot).To(BeADirectory())
		Expect(filepath.Join(projectRoot, fileGolangCI)).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".yamllint.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".markdownlint.yml")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".shellcheckrc")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "Dockerfile")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "docker-compose.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "scripts", "docker.sh")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "CODEOWNERS")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, "SUPPLY_CHAIN.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectRoot, ".gitlab", "issue_templates", "Bug.md")).To(BeAnExistingFile())

		// go style skips ginkgo suite generation.
		_, err = os.Stat(filepath.Join(projectRoot, "pkg", "e2egostyle", "e2egostyle_suite_test.go"))
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

		mainGo := readFile(filepath.Join(projectRoot, "cmd", "e2egostyle", "main.go"))
		Expect(mainGo).To(ContainSubstring("log/slog"))
		Expect(mainGo).NotTo(ContainSubstring("github.com/rs/zerolog"))

		assertGeneratedGitLabCIValid(projectRoot)
		assertYAMLFilesValid(projectRoot, []string{
			fileGolangCI,
			".yamllint.yml",
			"Taskfile.yml",
			"configs/e2egostyle.yml",
		})
	})

	It("generates deterministic output for base_project-only scaffolding", func() {
		defer withFakeGo()()

		tempDir := GinkgoT().TempDir()
		targetDir := filepath.Join(tempDir, "out")
		projectRoot := filepath.Join(targetDir, "proj")
		root := repoRoot(GinkgoT())
		projectBaseTemplates := filepath.Join(root, "templates", "project_base")
		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		baseProjectContent, err := loadTestFixtureWithVars(
			baseProjectFixtureFile,
			map[string]string{
				fixtureSourceDir: projectBaseTemplates,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseProjectCfg, string(baseProjectContent))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		gobootContent, err := loadTestFixtureWithVars(
			"cmd_goboot/goboot/service_execution_goboot.yml",
			map[string]string{
				fixtureTargetDir:   targetDir,
				fixtureBaseProject: baseProjectCfg,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		writeConfig(gobootCfg, string(gobootContent))

		Expect(run([]string{argConfig, gobootCfg})).To(Succeed())

		assertFilesExist(projectRoot, []string{
			"LICENSE",
			"NOTICE",
			"README.md",
			"ROADMAP.md",
			"PROJECT_STRUCTURE.md",
			"PROFILE.md",
			"VERSIONING.md",
			"WORKFLOW.md",
			"go.mod",
			"cmd/proj/main.go",
			"configs/proj.yml",
			"pkg/config/proj.go",
			"pkg/proj/proj.go",
			"pkg/projutils/projUtils.go",
		})
		assertFilesNotExist(projectRoot, []string{
			".gitlab-ci.yml",
			fileGolangCI,
			"scripts/lint.sh",
			"scripts/test.sh",
			"pkg/proj/proj_test.go",
		})

		goMod := readFile(filepath.Join(projectRoot, "go.mod"))
		Expect(goMod).To(ContainSubstring("module example.com/x"))
		Expect(goMod).NotTo(ContainSubstring("{{"))
		Expect(readFile(filepath.Join(projectRoot, "README.md"))).NotTo(ContainSubstring("{{"))
		Expect(readFile(filepath.Join(projectRoot, "PROFILE.md"))).To(ContainSubstring("standard"))
	})

	It("fails with a clear error when go mod tidy cannot resolve modules (proxy/path failure)", func() {
		defer withFakeGoScript(`#!/usr/bin/env bash
if [[ "$1" == "mod" && "$2" == "tidy" ]]; then
echo "go: module lookup disabled by GOPROXY=${GOPROXY:-}" >&2
exit 1
fi
exit 0
`)()

		origProxy, hadProxy := os.LookupEnv("GOPROXY")

		Expect(os.Setenv("GOPROXY", "off")).To(Succeed())

		defer func() {
			if hadProxy {
				_ = os.Setenv("GOPROXY", origProxy)

				return
			}

			_ = os.Unsetenv("GOPROXY")
		}()

		tempDir := GinkgoT().TempDir()
		targetDir := filepath.Join(tempDir, "out")
		projectRoot := filepath.Join(targetDir, "proj")
		root := repoRoot(GinkgoT())
		projectBaseTemplates := filepath.Join(root, "templates", "project_base")
		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		baseProjectContent, err := loadTestFixtureWithVars(
			baseProjectFixtureFile,
			map[string]string{
				fixtureSourceDir: projectBaseTemplates,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseProjectCfg, string(baseProjectContent))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		gobootContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/service_execution_goboot.yml", map[string]string{
			fixtureTargetDir:   targetDir,
			fixtureBaseProject: baseProjectCfg,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(gobootCfg, string(gobootContent))

		err = run([]string{argConfig, gobootCfg})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to run go mod tidy"))
		Expect(err.Error()).To(ContainSubstring("exit status 1"))

		Expect(projectRoot).To(BeADirectory())
		Expect(filepath.Join(projectRoot, "go.mod")).To(BeAnExistingFile())
	})

	It("fails fast when template content is not parseable", func() {
		tempDir := GinkgoT().TempDir()
		targetDir := filepath.Join(tempDir, "out")
		sourceDir := filepath.Join(tempDir, "templates")

		Expect(os.MkdirAll(sourceDir, 0o755)).To(Succeed())
		writeConfig(filepath.Join(sourceDir, "README.md.tmpl"), "project: {{ .ProjectName")
		baseProjectCfg := filepath.Join(tempDir, "base_project.yml")
		baseProjectContent, err := loadTestFixtureWithVars(
			baseProjectFixtureFile,
			map[string]string{
				fixtureSourceDir: sourceDir,
			},
		)
		Expect(err).NotTo(HaveOccurred())
		writeConfig(baseProjectCfg, string(baseProjectContent))

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		gobootContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/service_execution_goboot.yml", map[string]string{
			fixtureTargetDir:   targetDir,
			fixtureBaseProject: baseProjectCfg,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(gobootCfg, string(gobootContent))

		err = run([]string{argConfig, gobootCfg})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("service execution failed"))
		Expect(err.Error()).To(ContainSubstring("failed template parse"))
	})

	It("fails when target path is not writable", func() {
		if runtime.GOOS == "windows" {
			Skip("permission model differs on Windows for this scenario")
		}

		tempDir := GinkgoT().TempDir()
		lockedParent := filepath.Join(tempDir, "locked")
		Expect(os.MkdirAll(lockedParent, 0o755)).To(Succeed())

		Expect(os.Chmod(lockedParent, 0o500)).To(Succeed())
		defer func() {
			_ = os.Chmod(lockedParent, 0o755)
		}()

		targetDir := filepath.Join(lockedParent, "out")
		configFile := filepath.Join(tempDir, "goboot.yml")
		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "PermFail",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		writeConfig(configFile, string(yamlContent))

		err = run([]string{argConfig, configFile})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("service registration failed"))
		Expect(err.Error()).To(ContainSubstring("failed to create target directory"))
	})

	It("fails when an enabled lint template file is missing", func() {
		tempDir := GinkgoT().TempDir()
		emptyLintTemplates := filepath.Join(tempDir, "lint_templates")
		Expect(os.MkdirAll(emptyLintTemplates, 0o755)).To(Succeed())

		baseLintCfg := filepath.Join(tempDir, "base_lint.yml")
		writeConfig(baseLintCfg, "sourcePath: "+emptyLintTemplates+"\n"+
			"linters:\n"+
			"  golang:\n"+
			"    cmd: \"\"\n"+
			"    enabled: true\n")

		gobootCfg := filepath.Join(tempDir, "goboot.yml")
		writeConfig(gobootCfg, "projectName: \"MissingLintTemplate\"\n"+
			"repoUrl: \"https://github.com/example/missing-template\"\n"+
			"gitProvider: \"gitlab\"\n"+
			"targetPath: \""+filepath.Join(tempDir, "out")+"\"\n"+
			"services:\n"+
			"  - id: \"base_lint\"\n"+
			"    confPath: \""+baseLintCfg+"\"\n"+
			"    enabled: true\n")

		err := run([]string{argConfig, gobootCfg})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("service execution failed"))
		Expect(err.Error()).To(ContainSubstring("missing required template"))
		Expect(err.Error()).To(ContainSubstring(fileGolangCI))
	})
})
