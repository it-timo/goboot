package baseci_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baseci"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseCI Service", func() {
	const (
		policyStrict              = "strict"
		policyBalanced            = "balanced"
		gitProviderGitHub         = "github"
		gitProviderGitLab         = "gitlab"
		commandGoTest             = "go test ./..."
		commandGolangCILintDocker = "{{DOCKER_RUN}} golangci/golangci-lint:v2.12.2 golangci-lint run ./..."
		jobLint                   = "lint"
		fileScriptsLintYMLTmplStr = `{{ range $c := index .FileScripts "lint.yml" }}{{ $c }}{{ "\n" }}{{ end }}`
	)

	var (
		tempDir     string
		baseCI      *baseci.BaseCI
		validConfig *config.BaseCIConfig
	)

	BeforeEach(func() {
		var err error

		tempDir, err = os.MkdirTemp("", "baseci-test-*")
		Expect(err).NotTo(HaveOccurred())

		validConfig = &config.BaseCIConfig{
			SourcePath:  tempDir,
			ProjectName: "testproject",
			GitProvider: gitProviderGitLab,
			GoVersion:   []string{"1.26"},
			AutoBranches: []string{
				"main",
			},
			ImagePolicy: policyBalanced,
		}

		baseCI = baseci.NewBaseCI(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	ensureProviderDir := func(sourcePath, provider string) {
		providerDir := filepath.Join(sourcePath, strings.ToLower(strings.TrimSpace(provider)))
		Expect(os.MkdirAll(providerDir, 0o755)).To(Succeed())
	}

	Describe("NewBaseCI", func() {
		It("creates a new BaseCI instance", func() {
			ci := baseci.NewBaseCI("/some/path")
			Expect(ci).NotTo(BeNil())
		})
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseCI.ID()).To(Equal(goboottypes.ServiceNameBaseCI))
			Expect(baseCI.ID()).To(Equal("base_ci"))
		})
	})

	Describe("SetConfig", func() {
		var sourceDir string

		BeforeEach(func() {
			var err error

			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		Context("with valid config", func() {
			It("accepts BaseCIConfig", func() {
				validConfig.SourcePath = sourceDir
				ensureProviderDir(validConfig.SourcePath, validConfig.GitProvider)
				err := baseCI.SetConfig(validConfig)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid config type", func() {
			It("returns an error for wrong config type", func() {
				wrongConfig := &config.BaseProjectConfig{}
				err := baseCI.SetConfig(wrongConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})

			It("errors when source and target paths are identical", func() {
				validConfig.SourcePath = tempDir
				err := baseCI.SetConfig(validConfig)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Registrar Interface", func() {
		var sourceDir string

		BeforeEach(func() {
			var err error

			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())

			validConfig.SourcePath = sourceDir
			ensureProviderDir(validConfig.SourcePath, validConfig.GitProvider)
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		Describe("RegisterLines", func() {
			It("accepts script lines for registration", func() {
				lines := []string{commandGoTest}
				err := baseCI.RegisterLines("test_service", lines)
				Expect(err).NotTo(HaveOccurred())
			})

			It("prevents duplicate registrations for the same service", func() {
				Expect(baseCI.RegisterLines("dup_service", []string{"cmd"})).To(Succeed())
				err := baseCI.RegisterLines("dup_service", []string{"cmd"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already registered"))
			})
		})

		Describe("RegisterFile", func() {
			It("accepts file registration", func() {
				lines := []string{commandGoTest}
				err := baseCI.RegisterFile(goboottypes.CIFileLint, lines)
				Expect(err).NotTo(HaveOccurred())
			})

			It("prevents duplicate file registrations", func() {
				Expect(baseCI.RegisterFile(goboottypes.CIFileLint, []string{"echo"})).To(Succeed())
				err := baseCI.RegisterFile(goboottypes.CIFileLint, []string{"echo again"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already registered"))
			})
		})
	})

	Describe("Run", func() {
		var sourceDir string

		createSourceFile := func(relPath, content string) {
			fullPath := filepath.Join(
				sourceDir,
				strings.ToLower(validConfig.GitProvider),
				relPath+goboottypes.TemplateSuffix,
			)
			Expect(os.MkdirAll(filepath.Dir(fullPath), 0o755)).To(Succeed())
			Expect(os.WriteFile(fullPath, []byte(content), 0o644)).To(Succeed())
		}

		BeforeEach(func() {
			var err error

			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())

			validConfig.SourcePath = sourceDir
			ensureProviderDir(validConfig.SourcePath, validConfig.GitProvider)
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		It("copies and renders enabled templates", func() {
			createSourceFile(".gitlab-ci.yml", "root-{{.ProjectName}}-{{len .EnabledJobFiles}}")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(filepath.Join(".gitlab/ci", "lint.yml"),
				"lint-{{len (index .FileScripts \"lint.yml\")}}")

			Expect(baseCI.RegisterFile(goboottypes.CIFileLint, []string{commandGoTest})).To(Succeed())

			Expect(baseCI.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			Expect(filepath.Join(targetRoot, ".gitlab-ci.yml")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, ".gitlab/ci", "lint.yml")).To(BeAnExistingFile())

			content, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab-ci.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring(validConfig.ProjectName))

			scriptContent, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "lint.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(scriptContent)).To(ContainSubstring("1"))
		})

		It("skips job files when no commands are registered", func() {
			createSourceFile(".gitlab-ci.yml", "root")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(filepath.Join(".gitlab/ci", "lint.yml"), "lint")

			Expect(baseCI.Run()).To(Succeed())
			Expect(filepath.Join(tempDir, validConfig.ProjectName, ".gitlab/ci", "lint.yml")).NotTo(BeAnExistingFile())
		})

		It("returns error when a template file is missing", func() {
			err := baseCI.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required template"))
		})

		It("returns error when template rendering fails", func() {
			createSourceFile(".gitlab-ci.yml", "{{")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")

			err := baseCI.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed template render"))
		})

		It("normalizes known image tags to digest variables in strict mode", func() {
			validConfig.ImagePolicy = policyStrict
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())

			createSourceFile(".gitlab-ci.yml", "root")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(
				filepath.Join(".gitlab/ci", "lint.yml"),
				fileScriptsLintYMLTmplStr,
			)

			Expect(baseCI.RegisterFile(goboottypes.CIFileLint, []string{
				commandGolangCILintDocker,
			})).To(Succeed())

			Expect(baseCI.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			content, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "lint.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("$GOLANGCI_LINT_IMAGE"))
			Expect(string(content)).NotTo(ContainSubstring("golangci/golangci-lint:v2.12.2"))
		})

		It("registers config jobs with allowFailure flags and sorted enabled job files", func() {
			validConfig.Jobs = map[string]*config.CIJob{
				jobLint: {
					Commands:     []string{"echo lint"},
					AllowFailure: false,
				},
				"build": {
					Commands:     []string{"echo build"},
					AllowFailure: true,
				},
				"test": {
					Commands:     []string{"echo test"},
					AllowFailure: false,
				},
			}
			Expect(validConfig.Validate()).To(Succeed())
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())

			createSourceFile(".gitlab-ci.yml", `{{ range $f := .EnabledJobFiles }}{{ $f }}{{ "\n" }}{{ end }}`)
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(
				filepath.Join(".gitlab/ci", "lint.yml"),
				`{{ index .FileAllowFailure "lint.yml" }}:{{ len (index .FileScripts "lint.yml") }}`,
			)
			createSourceFile(
				filepath.Join(".gitlab/ci", "build.yml"),
				`{{ index .FileAllowFailure "build.yml" }}:{{ len (index .FileScripts "build.yml") }}`,
			)
			createSourceFile(
				filepath.Join(".gitlab/ci", "test.yml"),
				`{{ index .FileAllowFailure "test.yml" }}:{{ len (index .FileScripts "test.yml") }}`,
			)

			Expect(baseCI.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			rootContent, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab-ci.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(rootContent)).To(Equal("build.yml\nlint.yml\ntest.yml\n"))

			lintContent, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "lint.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(lintContent)).To(Equal("false:1"))

			buildContent, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "build.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(buildContent)).To(Equal("true:1"))

			testContent, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "test.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(testContent)).To(Equal("false:1"))
		})

		It("normalizes config job commands in strict policy mode", func() {
			validConfig.ImagePolicy = policyStrict
			validConfig.Jobs = map[string]*config.CIJob{
				jobLint: {
					Commands: []string{
						commandGolangCILintDocker,
					},
				},
			}
			Expect(validConfig.Validate()).To(Succeed())
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())

			createSourceFile(".gitlab-ci.yml", "root")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(
				filepath.Join(".gitlab/ci", "lint.yml"),
				fileScriptsLintYMLTmplStr,
			)

			Expect(baseCI.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			content, err := os.ReadFile(filepath.Join(targetRoot, ".gitlab/ci", "lint.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("$GOLANGCI_LINT_IMAGE"))
			Expect(string(content)).NotTo(ContainSubstring("golangci/golangci-lint:v2.12.2"))
		})

		It("fails when config job collides with an already registered service job", func() {
			validConfig.Jobs = map[string]*config.CIJob{
				"lint": {
					Commands: []string{"echo lint from config"},
				},
			}
			Expect(validConfig.Validate()).To(Succeed())
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())
			Expect(baseCI.RegisterLines("lint", []string{"echo lint from service"})).To(Succeed())

			createSourceFile(".gitlab-ci.yml", "root")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(filepath.Join(".gitlab/ci", "lint.yml"), "lint")

			err := baseCI.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to register CI jobs"))
			Expect(err.Error()).To(ContainSubstring("already registered in ci"))
		})

		It("fails when config job collides with an already registered job file", func() {
			validConfig.Jobs = map[string]*config.CIJob{
				"lint": {
					Commands: []string{"echo lint from config"},
				},
			}
			Expect(validConfig.Validate()).To(Succeed())
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())
			Expect(baseCI.RegisterFile(goboottypes.CIFileLint, []string{"echo lint from service"})).To(Succeed())

			createSourceFile(".gitlab-ci.yml", "root")
			createSourceFile(filepath.Join(".gitlab/ci", "commands.yml"), "commands")
			createSourceFile(filepath.Join(".gitlab/ci", "versions.yml"), "versions")
			createSourceFile(filepath.Join(".gitlab/ci", "lint.yml"), "lint")

			err := baseCI.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to register CI jobs"))
			Expect(err.Error()).To(ContainSubstring("already registered in ci"))
		})

		It("renders github strict lint workflow with digest env variables", func() {
			validConfig.ImagePolicy = policyStrict
			validConfig.GitProvider = gitProviderGitHub
			ensureProviderDir(validConfig.SourcePath, validConfig.GitProvider)
			Expect(baseCI.SetConfig(validConfig)).To(Succeed())

			createSourceFile(
				filepath.Join(".github", "workflows", "lint.yml"),
				fileScriptsLintYMLTmplStr,
			)

			Expect(baseCI.RegisterFile(goboottypes.CIFileLint, []string{
				commandGolangCILintDocker,
			})).To(Succeed())

			Expect(baseCI.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			content, err := os.ReadFile(filepath.Join(targetRoot, ".github", "workflows", "lint.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("$GOLANGCI_LINT_IMAGE"))
			Expect(string(content)).NotTo(ContainSubstring("golangci/golangci-lint:v2.12.2"))
		})

		It("renders provider and imagePolicy combinations without template leftovers", func() {
			type providerPolicy struct {
				provider string
				policy   string
			}

			cases := []providerPolicy{
				{provider: gitProviderGitLab, policy: policyBalanced},
				{provider: gitProviderGitLab, policy: policyStrict},
				{provider: gitProviderGitLab, policy: "simple"},
				{provider: gitProviderGitHub, policy: policyBalanced},
				{provider: gitProviderGitHub, policy: policyStrict},
				{provider: gitProviderGitHub, policy: "simple"},
			}

			for _, curCase := range cases {
				By(fmt.Sprintf("rendering %s/%s", curCase.provider, curCase.policy))

				sourceCaseDir, err := os.MkdirTemp("", "source-case-*")
				Expect(err).NotTo(HaveOccurred())

				writeTemplate := func(relPath, content string) {
					fullPath := filepath.Join(
						sourceCaseDir,
						curCase.provider,
						relPath+goboottypes.TemplateSuffix,
					)
					Expect(os.MkdirAll(filepath.Dir(fullPath), 0o755)).To(Succeed())
					Expect(os.WriteFile(fullPath, []byte(content), 0o644)).To(Succeed())
				}

				switch curCase.provider {
				case gitProviderGitLab:
					writeTemplate(".gitlab-ci.yml", "stages:\n  - lint")
					writeTemplate(filepath.Join(".gitlab/ci", "commands.yml"), "vars")
					writeTemplate(filepath.Join(".gitlab/ci", "versions.yml"), "vars")
					writeTemplate(
						filepath.Join(".gitlab/ci", "lint.yml"),
						fileScriptsLintYMLTmplStr,
					)
				case gitProviderGitHub:
					writeTemplate(
						filepath.Join(".github", "workflows", "lint.yml"),
						fileScriptsLintYMLTmplStr,
					)
				}

				curCfg := &config.BaseCIConfig{
					SourcePath:  sourceCaseDir,
					ProjectName: fmt.Sprintf("proj-%s-%s", curCase.provider, curCase.policy),
					GitProvider: curCase.provider,
					GoVersion:   []string{"1.26"},
					AutoBranches: []string{
						"main",
						"master",
					},
					ImagePolicy: curCase.policy,
				}

				curCI := baseci.NewBaseCI(tempDir)
				Expect(curCI.SetConfig(curCfg)).To(Succeed())
				Expect(curCI.RegisterFile(goboottypes.CIFileLint, []string{
					commandGolangCILintDocker,
				})).To(Succeed())
				Expect(curCI.Run()).To(Succeed())

				var outPath string

				switch curCase.provider {
				case gitProviderGitLab:
					outPath = filepath.Join(tempDir, curCfg.ProjectName, ".gitlab/ci", "lint.yml")
				case gitProviderGitHub:
					outPath = filepath.Join(tempDir, curCfg.ProjectName, ".github", "workflows", "lint.yml")
				}

				content, err := os.ReadFile(outPath)
				Expect(err).NotTo(HaveOccurred())

				if curCase.policy == policyStrict {
					Expect(string(content)).To(ContainSubstring("$GOLANGCI_LINT_IMAGE"))
				} else {
					Expect(string(content)).To(ContainSubstring("golangci/golangci-lint:v2.12.2"))
				}

				Expect(os.RemoveAll(sourceCaseDir)).To(Succeed())
			}
		})
	})
})
