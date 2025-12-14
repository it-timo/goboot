package basetest_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/basetest"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

// mockRegistrar implements goboottypes.Registrar for testing.
type mockRegistrar struct {
	registeredLines map[string][]string
	registeredFiles map[string][]string
	registerErr     error
}

func newMockRegistrar() *mockRegistrar {
	return &mockRegistrar{
		registeredLines: make(map[string][]string),
		registeredFiles: make(map[string][]string),
	}
}

func (m *mockRegistrar) RegisterLines(id string, lines []string) error {
	if m.registerErr != nil {
		return m.registerErr
	}

	m.registeredLines[id] = lines

	return nil
}

func (m *mockRegistrar) RegisterFile(id string, lines []string) error {
	if m.registerErr != nil {
		return m.registerErr
	}

	m.registeredFiles[id] = lines

	return nil
}

var _ = Describe("BaseTest", func() {
	var (
		tmpUserDir string
		tmpSrcDir  string
		baseTest   *basetest.BaseTest
		cfg        *config.BaseTestConfig
		mockReg    *mockRegistrar
	)

	BeforeEach(func() {
		var err error
		// Create a temp directory to simulate user project target
		tmpUserDir, err = os.MkdirTemp("", "goboot-target-*")
		Expect(err).NotTo(HaveOccurred())

		// Create a temp directory to simulate source templates
		tmpSrcDir, err = os.MkdirTemp("", "goboot-src-*")
		Expect(err).NotTo(HaveOccurred())

		mockReg = newMockRegistrar()
	})

	AfterEach(func() {
		if tmpUserDir != "" {
			Expect(os.RemoveAll(tmpUserDir)).To(Succeed())
		}
		if tmpSrcDir != "" {
			Expect(os.RemoveAll(tmpSrcDir)).To(Succeed())
		}
	})

	Describe("NewBaseTest", func() {
		It("creates a new instance with target directory", func() {
			bt := basetest.NewBaseTest(tmpUserDir)
			Expect(bt).NotTo(BeNil())
		})

		It("returns instance with correct ID", func() {
			bt := basetest.NewBaseTest(tmpUserDir)
			Expect(bt.ID()).To(Equal(goboottypes.ServiceNameBaseTest))
		})
	})

	Describe("ID", func() {
		BeforeEach(func() {
			baseTest = basetest.NewBaseTest(tmpUserDir)
		})

		It("returns the correct service identifier", func() {
			Expect(baseTest.ID()).To(Equal(goboottypes.ServiceNameBaseTest))
			Expect(baseTest.ID()).To(Equal("base_test"))
		})
	})

	Describe("SetConfig", func() {
		BeforeEach(func() {
			baseTest = basetest.NewBaseTest(tmpUserDir)
			cfg = &config.BaseTestConfig{
				SourcePath:       tmpSrcDir,
				ProjectName:      "MyProject",
				TestCMD:          goboottypes.DefaultGoTestCMD,
				RepoImportPath:   "github.com/example/myproject",
				UseStyle:         goboottypes.TestStyleGinkgo,
				CapsProjectName:  "MYPROJECT",
				LowerProjectName: "myproject",
			}
		})

		Context("with valid config", func() {
			It("sets the config successfully", func() {
				err := baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("accepts validated config", func() {
				err := cfg.Validate()
				Expect(err).NotTo(HaveOccurred())

				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid config type", func() {
			It("returns error for wrong config type", func() {
				wrongConfig := &config.BaseLintConfig{}
				err := baseTest.SetConfig(wrongConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})
		})

		Context("with path comparison", func() {
			It("errors when source and target are the same", func() {
				// Set source path to target path
				cfg.SourcePath = tmpUserDir
				err := baseTest.SetConfig(cfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed path comparison"))
			})

			It("accepts different source and target paths", func() {
				// Default setup has different paths
				err := baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("SetScriptReceiver", func() {
		BeforeEach(func() {
			baseTest = basetest.NewBaseTest(tmpUserDir)
		})

		It("sets the registrar successfully", func() {
			baseTest.SetScriptReceiver(mockReg)
			// No error, just sets the receiver
		})

		It("accepts nil registrar", func() {
			baseTest.SetScriptReceiver(nil)
			// Should not panic
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			baseTest = basetest.NewBaseTest(tmpUserDir)
		})

		Context("with basic file structure", func() {
			BeforeEach(func() {
				// Create a simple template file
				err := os.WriteFile(filepath.Join(tmpSrcDir, "README.md.tmpl"), []byte("# {{.ProjectName}} Tests"), 0644)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "MyProject",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/example/myproject",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "MYPROJECT",
					LowerProjectName: "myproject",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("copies and renders files successfully", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check file exists (template suffix removed)
				expectedFile := filepath.Join(tmpUserDir, "MyProject", "README.md")
				Expect(expectedFile).To(BeAnExistingFile())

				// Check content was rendered
				content, err := os.ReadFile(expectedFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(Equal("# MyProject Tests"))
			})
		})

		Context("with directory structure and path templates", func() {
			BeforeEach(func() {
				// Create directory with template in name
				err := os.MkdirAll(filepath.Join(tmpSrcDir, "{{.LowerProjectName}}utils"), 0755)
				Expect(err).NotTo(HaveOccurred())

				// Create file inside templated directory
				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "{{.LowerProjectName}}utils", "utils_test.go.tmpl"),
					[]byte("package utils\n// Test for {{.ProjectName}}"),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "MyProject",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/example/myproject",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "MYPROJECT",
					LowerProjectName: "myproject",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("renders directory names from templates", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check directory was renamed
				expectedDir := filepath.Join(tmpUserDir, "MyProject", "myprojectutils")
				Expect(expectedDir).To(BeADirectory())
			})

			It("renders file paths and content", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check file exists with template suffix removed
				expectedFile := filepath.Join(tmpUserDir, "MyProject", "myprojectutils", "utils_test.go")
				Expect(expectedFile).To(BeAnExistingFile())

				// Check content was rendered
				content, err := os.ReadFile(expectedFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("// Test for MyProject"))
			})
		})

		Context("with multiple files and directories", func() {
			BeforeEach(func() {
				// Create multiple directories and files
				err := os.MkdirAll(filepath.Join(tmpSrcDir, "pkg", "config"), 0755)
				Expect(err).NotTo(HaveOccurred())

				err = os.MkdirAll(filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}"), 0755)
				Expect(err).NotTo(HaveOccurred())

				// Files in different locations
				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "pkg", "config", "config_test.go.tmpl"),
					[]byte("package config\n// {{.CapsProjectName}} config tests"),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}", "main_test.go.tmpl"),
					[]byte("package {{.LowerProjectName}}\n"),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "TestApp",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/test/app",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "TESTAPP",
					LowerProjectName: "testapp",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates complete directory structure", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check all directories exist
				Expect(filepath.Join(tmpUserDir, "TestApp", "pkg", "config")).To(BeADirectory())
				Expect(filepath.Join(tmpUserDir, "TestApp", "pkg", "testapp")).To(BeADirectory())
			})

			It("renders all files correctly", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check config test file
				configTest := filepath.Join(tmpUserDir, "TestApp", "pkg", "config", "config_test.go")
				Expect(configTest).To(BeAnExistingFile())
				content, err := os.ReadFile(configTest)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("TESTAPP config tests"))

				// Check main test file
				mainTest := filepath.Join(tmpUserDir, "TestApp", "pkg", "testapp", "main_test.go")
				Expect(mainTest).To(BeAnExistingFile())
				content, err = os.ReadFile(mainTest)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("package testapp"))
			})
		})

		Context("with empty source directory", func() {
			BeforeEach(func() {
				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "EmptyProject",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/test/empty",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "EMPTYPROJECT",
					LowerProjectName: "emptyproject",
				}
				err := baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates root directory even with no files", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				rootDir := filepath.Join(tmpUserDir, "EmptyProject")
				Expect(rootDir).To(BeADirectory())
			})
		})

		Context("with script registration", func() {
			BeforeEach(func() {
				// Create a simple file
				err := os.WriteFile(filepath.Join(tmpSrcDir, "test.go"), []byte("package test"), 0644)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "ScriptTest",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/test/script",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "SCRIPTTEST",
					LowerProjectName: "scripttest",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("registers test scripts when receiver is set", func() {
				baseTest.SetScriptReceiver(mockReg)
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check that scripts were registered
				Expect(mockReg.registeredLines).To(HaveKey(goboottypes.ServiceNameBaseTest))
				Expect(mockReg.registeredLines[goboottypes.ServiceNameBaseTest]).To(ContainElement(goboottypes.DefaultGoTestCMD))
			})

			It("skips script registration when receiver is nil", func() {
				baseTest.SetScriptReceiver(nil)
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Should not panic, just skip registration
			})
		})

		Context("with template errors", func() {
			It("returns error for invalid path template syntax", func() {
				// Create file with invalid template in path
				err := os.MkdirAll(filepath.Join(tmpSrcDir, "{{.InvalidField"), 0755)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "BadPath",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/test/bad",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "BADPATH",
					LowerProjectName: "badpath",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())

				err = baseTest.Run()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to render path"))
			})

			It("returns error for invalid content template syntax", func() {
				// Create file with invalid template in content
				err := os.WriteFile(
					filepath.Join(tmpSrcDir, "bad.go.tmpl"),
					[]byte("{{.NoSuchField}}"),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "BadContent",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/test/bad",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "BADCONTENT",
					LowerProjectName: "badcontent",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())

				err = baseTest.Run()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to render"))
			})
		})
	})

	Describe("Real-world scenarios", func() {
		BeforeEach(func() {
			baseTest = basetest.NewBaseTest(tmpUserDir)
		})

		Context("when setting up Ginkgo test suite", func() {
			BeforeEach(func() {
				// Create realistic Ginkgo test structure
				err := os.MkdirAll(filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}"), 0755)
				Expect(err).NotTo(HaveOccurred())

				// Suite file
				suiteContent := `package {{.LowerProjectName}}_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test{{.CapsProjectName}}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "{{.ProjectName}} Suite")
}
`
				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}", "{{.LowerProjectName}}_suite_test.go.tmpl"),
					[]byte(suiteContent),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				// Test file
				testContent := `package {{.LowerProjectName}}_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.ProjectName}}", func() {
	It("works", func() {
		Expect(true).To(BeTrue())
	})
})
`
				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}", "{{.LowerProjectName}}_test.go.tmpl"),
					[]byte(testContent),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "MyApp",
					TestCMD:          goboottypes.DefaultGoTestCMD,
					RepoImportPath:   "github.com/example/myapp",
					UseStyle:         goboottypes.TestStyleGinkgo,
					CapsProjectName:  "MYAPP",
					LowerProjectName: "myapp",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates complete Ginkgo test suite structure", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check suite file
				suiteFile := filepath.Join(tmpUserDir, "MyApp", "pkg", "myapp", "myapp_suite_test.go")
				Expect(suiteFile).To(BeAnExistingFile())

				content, err := os.ReadFile(suiteFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("package myapp_test"))
				Expect(string(content)).To(ContainSubstring("func TestMYAPP(t *testing.T)"))
				Expect(string(content)).To(ContainSubstring("MyApp Suite"))

				// Check test file
				testFile := filepath.Join(tmpUserDir, "MyApp", "pkg", "myapp", "myapp_test.go")
				Expect(testFile).To(BeAnExistingFile())

				content, err = os.ReadFile(testFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("Describe(\"MyApp\""))
			})
		})

		Context("when setting up standard Go test suite", func() {
			BeforeEach(func() {
				// Create realistic standard Go test structure
				err := os.MkdirAll(filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}"), 0755)
				Expect(err).NotTo(HaveOccurred())

				testContent := `package {{.LowerProjectName}}

import (
	"testing"
)

func Test{{.ProjectName}}Functionality(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		// Test implementation
	})
}
`
				err = os.WriteFile(
					filepath.Join(tmpSrcDir, "pkg", "{{.LowerProjectName}}", "{{.LowerProjectName}}_test.go.tmpl"),
					[]byte(testContent),
					0644,
				)
				Expect(err).NotTo(HaveOccurred())

				cfg = &config.BaseTestConfig{
					SourcePath:       tmpSrcDir,
					ProjectName:      "StandardApp",
					RepoImportPath:   "github.com/example/standardapp",
					UseStyle:         goboottypes.TestStyleGo,
					CapsProjectName:  "STANDARDAPP",
					LowerProjectName: "standardapp",
				}
				err = baseTest.SetConfig(cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("creates complete standard Go test structure", func() {
				err := baseTest.Run()
				Expect(err).NotTo(HaveOccurred())

				// Check test file
				testFile := filepath.Join(tmpUserDir, "StandardApp", "pkg", "standardapp", "standardapp_test.go")
				Expect(testFile).To(BeAnExistingFile())

				content, err := os.ReadFile(testFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("package standardapp"))
				Expect(string(content)).To(ContainSubstring("func TestStandardAppFunctionality(t *testing.T)"))
			})
		})
	})
})
