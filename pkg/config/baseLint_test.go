package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseLintConfig", func() {
	var (
		baseLint *config.BaseLintConfig
		tempDir  string
		testPath = "github.com/test/testproject"
	)

	BeforeEach(func() {
		baseLint = &config.BaseLintConfig{
			SourcePath:     "./templates/lint_base",
			ProjectName:    "testproject",
			RepoImportPath: testPath,
			Linters: map[string]*config.Linter{
				goboottypes.LinterGo: {
					Enabled: true,
				},
				goboottypes.LinterYAML: {
					Enabled: true,
					Cmd:     "custom-yamllint .",
				},
			},
			AllowedPackages: []string{
				"github.com/allowed/pkg",
			},
		}

		var err error
		tempDir, err = os.MkdirTemp("", "baselint-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseLint.ID()).To(Equal(goboottypes.ServiceNameBaseLint))
			Expect(baseLint.ID()).To(Equal("base_lint"))
		})
	})

	Describe("Validate", func() {
		Context("with valid config", func() {
			It("validates successfully", func() {
				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("fills in default linter commands for enabled linters without custom commands", func() {
				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())

				// Go linter should have default command filled in
				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(Equal(goboottypes.DefaultGoLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(ContainSubstring("golangci-lint"))

				// YAML linter should keep custom command
				Expect(baseLint.Linters[goboottypes.LinterYAML].Cmd).To(Equal("custom-yamllint ."))

				// AllowedPackages should be set
				Expect(baseLint.AllowedPackages).To(ContainElement("github.com/allowed/pkg"))
			})
		})

		Context("with missing required fields", func() {
			It("errors when sourcePath is missing", func() {
				baseLint.SourcePath = ""
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("errors when projectName is missing", func() {
				baseLint.ProjectName = ""
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			It("errors when repoImportPath is missing", func() {
				baseLint.RepoImportPath = ""
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("repoImportPath"))
			})

			It("errors when multiple fields are missing", func() {
				baseLint.SourcePath = ""
				baseLint.ProjectName = ""
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})
		})

		Context("with whitespace-only values", func() {
			It("treats whitespace-only sourcePath as missing", func() {
				baseLint.SourcePath = "    "
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("treats whitespace-only projectName as missing", func() {
				baseLint.ProjectName = "\t\n"
				err := baseLint.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})
		})

		Context("when applying default linter commands", func() {
			It("fills in all standard default commands", func() {
				baseLint.Linters = map[string]*config.Linter{
					goboottypes.LinterGo:    {Enabled: true},
					goboottypes.LinterYAML:  {Enabled: true},
					goboottypes.LinterMake:  {Enabled: true},
					goboottypes.LinterMD:    {Enabled: true},
					goboottypes.LinterShell: {Enabled: true},
					goboottypes.LinterSHFMT: {Enabled: true},
				}

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(Equal(goboottypes.DefaultGoLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterYAML].Cmd).To(Equal(goboottypes.DefaultYMLLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterMake].Cmd).To(Equal(goboottypes.DefaultMakeLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterMD].Cmd).To(Equal(goboottypes.DefaultMDLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterShell].Cmd).To(Equal(goboottypes.DefaultShellLintCmd))
				Expect(baseLint.Linters[goboottypes.LinterSHFMT].Cmd).To(Equal(goboottypes.DefaultSHFMTCmd))
			})

			It("does not overwrite custom commands", func() {
				customCmd := "my-custom-linter"
				baseLint.Linters[goboottypes.LinterGo].Cmd = customCmd

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(Equal(customCmd))
			})

			It("does not fill commands for disabled linters", func() {
				baseLint.Linters = map[string]*config.Linter{
					goboottypes.LinterGo: {Enabled: false},
				}

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(BeEmpty())
			})

			It("leaves unknown enabled linters without a command", func() {
				baseLint.Linters = map[string]*config.Linter{
					"unknown": {Enabled: true},
				}

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(baseLint.Linters["unknown"].Cmd).To(BeEmpty())
			})

			It("keeps commands on disabled linters intact", func() {
				baseLint.Linters = map[string]*config.Linter{
					goboottypes.LinterGo: {Enabled: false, Cmd: "do-not-run"},
				}

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(Equal("do-not-run"))
			})

			It("is idempotent across multiple validations", func() {
				customCmd := "my-custom-linter"
				baseLint.Linters[goboottypes.LinterGo].Cmd = customCmd

				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())

				err = baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(baseLint.Linters[goboottypes.LinterGo].Cmd).To(Equal(customCmd))
			})

			It("handles a nil linter map", func() {
				baseLint.Linters = nil
				err := baseLint.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_lint.yml")
		})

		Context("with valid YAML file", func() {
			It("loads the configuration successfully", func() {
				yamlContent := `sourcePath: ./templates/lint
linters:
  golang:
    enabled: true
  yaml:
    enabled: true
    cmd: yamllint .
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseLintConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.SourcePath).To(Equal("./templates/lint"))
				Expect(newConfig.RepoImportPath).To(Equal(testPath))
				Expect(newConfig.Linters).To(HaveKey("golang"))
				Expect(newConfig.Linters["golang"].Enabled).To(BeTrue())
				Expect(newConfig.Linters["yaml"].Cmd).To(Equal("yamllint ."))
			})
		})

		Context("with non-existent file", func() {
			It("returns an error", func() {
				err := baseLint.ReadConfig("/nonexistent/path.yml", testPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid YAML", func() {
			It("returns an error", func() {
				invalidYAML := `sourcePath: ./templates
projectName: [invalid yaml structure
`
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseLintConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Linter struct", func() {
		Context("when creating linter instances", func() {
			It("can have custom commands", func() {
				linter := &config.Linter{
					Cmd:     "custom-lint-command",
					Enabled: true,
				}

				Expect(linter.Cmd).To(Equal("custom-lint-command"))
				Expect(linter.Enabled).To(BeTrue())
			})

			It("can be disabled", func() {
				linter := &config.Linter{
					Cmd:     "some-command",
					Enabled: false,
				}

				Expect(linter.Enabled).To(BeFalse())
			})
		})
	})
})
