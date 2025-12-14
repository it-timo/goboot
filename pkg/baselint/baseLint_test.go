package baselint_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baselint"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

type recordingRegistrar struct {
	linesCalls map[string][]string
	fileCalls  map[string][]string
	linesErr   error
	fileErr    error
}

func (r *recordingRegistrar) RegisterLines(name string, lines []string) error {
	if r.linesCalls == nil {
		r.linesCalls = make(map[string][]string)
	}

	if r.linesErr != nil {
		return r.linesErr
	}

	r.linesCalls[name] = append([]string{}, lines...)

	return nil
}

func (r *recordingRegistrar) RegisterFile(name string, lines []string) error {
	if r.fileCalls == nil {
		r.fileCalls = make(map[string][]string)
	}

	if r.fileErr != nil {
		return r.fileErr
	}

	r.fileCalls[name] = append([]string{}, lines...)

	return nil
}

var _ = Describe("BaseLint Service", func() {
	var (
		tempDir     string
		baseLint    *baselint.BaseLint
		validConfig *config.BaseLintConfig
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "baselint-test-*")
		Expect(err).NotTo(HaveOccurred())

		validConfig = &config.BaseLintConfig{
			SourcePath:     tempDir,
			ProjectName:    "testproject",
			RepoImportPath: "github.com/test/testproject",
			Linters: map[string]*config.Linter{
				goboottypes.LinterGo: {
					Enabled: true,
					Cmd:     goboottypes.DefaultGoLintCmd,
				},
			},
		}

		baseLint = baselint.NewBaseLint(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("NewBaseLint", func() {
		It("creates a new BaseLint instance", func() {
			bl := baselint.NewBaseLint("/some/path")
			Expect(bl).NotTo(BeNil())
		})
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseLint.ID()).To(Equal(goboottypes.ServiceNameBaseLint))
			Expect(baseLint.ID()).To(Equal("base_lint"))
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
			It("accepts BaseLintConfig", func() {
				validConfig.SourcePath = sourceDir
				err := baseLint.SetConfig(validConfig)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid config type", func() {
			It("returns an error for wrong config type", func() {
				wrongConfig := &config.BaseProjectConfig{}
				err := baseLint.SetConfig(wrongConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})
		})

		It("fails when source and target paths are identical", func() {
			validConfig.SourcePath = tempDir
			err := baseLint.SetConfig(validConfig)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SetScriptReceiver", func() {
		It("accepts a registrar implementation", func() {
			// We can't easily test the full registrar without mocking
			// but we can verify the method exists and doesn't panic
			Expect(func() {
				baseLint.SetScriptReceiver(nil)
			}).NotTo(Panic())
		})
	})

	Describe("Run", func() {
		var sourceDir string

		createTemplate := func(name, content string) {
			Expect(os.WriteFile(filepath.Join(sourceDir, name+goboottypes.TemplateSuffix), []byte(content), 0o644)).To(Succeed())
		}

		BeforeEach(func() {
			var err error
			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())

			validConfig.SourcePath = sourceDir
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		It("copies and renders enabled linter templates", func() {
			createTemplate(".golangci.yml", "run: {{ .ProjectName }} {{ .RepoImportPath }}")
			createTemplate(".yamllint.yml", "yaml: {{ .ProjectName }}")

			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo:   {Enabled: true, Cmd: "go-cmd"},
				goboottypes.LinterYAML: {Enabled: true, Cmd: "yaml-cmd"},
			}

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())
			Expect(baseLint.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			golangFile := filepath.Join(targetRoot, ".golangci.yml")
			yamlFile := filepath.Join(targetRoot, ".yamllint.yml")

			Expect(golangFile).To(BeAnExistingFile())
			Expect(yamlFile).To(BeAnExistingFile())

			content, err := os.ReadFile(golangFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring(validConfig.ProjectName))
			Expect(string(content)).To(ContainSubstring(validConfig.RepoImportPath))
		})

		It("skips disabled or unknown linters", func() {
			createTemplate(".golangci.yml", "run: {{ .ProjectName }}")

			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo:   {Enabled: false, Cmd: "go-cmd"},
				"unknown":              {Enabled: true, Cmd: "unknown-cmd"},
				goboottypes.LinterYAML: {Enabled: false, Cmd: "yaml-cmd"},
			}

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())
			Expect(baseLint.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			Expect(filepath.Join(targetRoot, ".golangci.yml")).NotTo(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, ".yamllint.yml")).NotTo(BeAnExistingFile())
		})

		It("registers scripts for enabled linters only", func() {
			createTemplate(".golangci.yml", "go")

			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo:   {Enabled: true, Cmd: "go-cmd"},
				goboottypes.LinterYAML: {Enabled: false, Cmd: "yaml-cmd"},
				"unknown":              {Enabled: true, Cmd: "unknown-cmd"},
			}

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())

			registrar := &recordingRegistrar{}
			baseLint.SetScriptReceiver(registrar)

			Expect(baseLint.Run()).To(Succeed())

			Expect(registrar.linesCalls).To(HaveKey(goboottypes.ServiceNameBaseLint))
			Expect(registrar.linesCalls[goboottypes.ServiceNameBaseLint]).To(ConsistOf("go-cmd"))
			Expect(registrar.fileCalls).To(HaveKey(goboottypes.ScriptFileLint))
			Expect(registrar.fileCalls[goboottypes.ScriptFileLint]).To(ConsistOf("go-cmd"))
		})

		It("returns an error when template file is missing", func() {
			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo: {Enabled: true, Cmd: "go-cmd"},
			}
			// no template created

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())
			err := baseLint.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required template \".golangci.yml"))
		})

		It("propagates template rendering errors", func() {
			createTemplate(".golangci.yml", "{{") // invalid template

			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo: {Enabled: true, Cmd: "go-cmd"},
			}

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())
			err := baseLint.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed template render"))
		})

		It("propagates registrar errors", func() {
			createTemplate(".golangci.yml", "go")
			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo: {Enabled: true, Cmd: "go-cmd"},
			}
			Expect(baseLint.SetConfig(validConfig)).To(Succeed())

			registrar := &recordingRegistrar{linesErr: errors.New("lines boom")}
			baseLint.SetScriptReceiver(registrar)

			err := baseLint.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to register script commands"))
		})

		It("uses default linter command when none is provided", func() {
			createTemplate(".golangci.yml", "go")

			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo: {Enabled: true, Cmd: ""},
			}
			// mimic validation pipeline to populate default commands.
			Expect(validConfig.Validate()).To(Succeed())

			Expect(baseLint.SetConfig(validConfig)).To(Succeed())

			registrar := &recordingRegistrar{}
			baseLint.SetScriptReceiver(registrar)

			Expect(baseLint.Run()).To(Succeed())

			Expect(registrar.linesCalls).To(HaveKey(goboottypes.ServiceNameBaseLint))
			Expect(registrar.linesCalls[goboottypes.ServiceNameBaseLint][0]).To(ContainSubstring("golangci-lint run"))
		})

		It("propagates registrar file errors", func() {
			createTemplate(".golangci.yml", "go")
			validConfig.Linters = map[string]*config.Linter{
				goboottypes.LinterGo: {Enabled: true, Cmd: "go-cmd"},
			}
			Expect(baseLint.SetConfig(validConfig)).To(Succeed())

			registrar := &recordingRegistrar{fileErr: errors.New("file boom")}
			baseLint.SetScriptReceiver(registrar)

			err := baseLint.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to register script file"))
		})
	})
})
