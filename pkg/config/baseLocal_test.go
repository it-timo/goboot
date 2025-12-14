package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseLocalConfig", func() {
	var (
		baseLocal *config.BaseLocalConfig
		tempDir   string
	)

	BeforeEach(func() {
		baseLocal = &config.BaseLocalConfig{
			SourcePath:  "./templates/local_base",
			ProjectName: "testproject",
			FileList:    []string{"Makefile", "Taskfile.yml", ".editorconfig"},
		}

		var err error
		tempDir, err = os.MkdirTemp("", "baselocal-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseLocal.ID()).To(Equal(goboottypes.ServiceNameBaseLocal))
			Expect(baseLocal.ID()).To(Equal("base_local"))
		})
	})

	Describe("Validate", func() {
		Context("with valid config", func() {
			It("validates successfully", func() {
				err := baseLocal.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with missing required fields", func() {
			It("errors when sourcePath is missing", func() {
				baseLocal.SourcePath = ""
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("errors when projectName is missing", func() {
				baseLocal.ProjectName = ""
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			It("errors when fileList is empty", func() {
				baseLocal.FileList = []string{}
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fileList"))
			})

			It("errors when fileList is nil", func() {
				baseLocal.FileList = nil
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fileList"))
			})

			It("errors when multiple fields are missing", func() {
				baseLocal.SourcePath = ""
				baseLocal.ProjectName = ""
				baseLocal.FileList = nil
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
				Expect(err.Error()).To(ContainSubstring("projectName"))
				Expect(err.Error()).To(ContainSubstring("fileList"))
			})
		})

		Context("with whitespace-only values", func() {
			It("treats whitespace-only sourcePath as missing", func() {
				baseLocal.SourcePath = "   \t\n"
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("treats whitespace-only projectName as missing", func() {
				baseLocal.ProjectName = "  "
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})
		})

		Context("with valid fileList", func() {
			It("accepts single file", func() {
				baseLocal.FileList = []string{"Makefile"}
				err := baseLocal.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("accepts multiple files", func() {
				baseLocal.FileList = []string{
					"Makefile",
					"Taskfile.yml",
					".editorconfig",
					".gitignore",
					"scripts/lint.sh",
				}
				err := baseLocal.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid fileList entries", func() {
			It("errors on blank entries", func() {
				baseLocal.FileList = []string{"Makefile", " "}
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("blank entries"))
			})

			It("errors on duplicate entries", func() {
				baseLocal.FileList = []string{"Makefile", "Makefile"}
				err := baseLocal.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("duplicates"))
			})
		})
	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_local.yml")
		})

		Context("with valid YAML file", func() {
			It("loads the configuration successfully", func() {
				yamlContent := `sourcePath: ./templates/local
fileList:
  - Makefile
  - Taskfile.yml
  - .editorconfig
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseLocalConfig{}
				err = newConfig.ReadConfig(configPath, "")
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.SourcePath).To(Equal("./templates/local"))
				Expect(newConfig.ProjectName).To(Equal(""))
				Expect(newConfig.FileList).To(HaveLen(3))
				Expect(newConfig.FileList).To(ContainElement("Makefile"))
				Expect(newConfig.FileList).To(ContainElement("Taskfile.yml"))
				Expect(newConfig.FileList).To(ContainElement(".editorconfig"))
			})

			It("loads single file in fileList", func() {
				yamlContent := `sourcePath: ./templates
projectName: proj
fileList:
  - single.txt
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseLocalConfig{}
				err = newConfig.ReadConfig(configPath, "")
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.FileList).To(HaveLen(1))
				Expect(newConfig.FileList[0]).To(Equal("single.txt"))
			})
		})

		Context("with non-existent file", func() {
			It("returns an error", func() {
				err := baseLocal.ReadConfig("/nonexistent/path.yml", "")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid YAML", func() {
			It("returns an error", func() {
				invalidYAML := `sourcePath: ./templates
projectName: [invalid structure
fileList:
  - not properly closed
`
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseLocalConfig{}
				err = newConfig.ReadConfig(configPath, "")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FileList behavior", func() {
		Context("when validating fileList content", func() {
			It("accepts paths with subdirectories", func() {
				baseLocal.FileList = []string{
					"scripts/lint.sh",
					"scripts/format.sh",
					".github/workflows/ci.yml",
				}
				err := baseLocal.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("accepts dotfiles", func() {
				baseLocal.FileList = []string{
					".editorconfig",
					".gitignore",
					".dockerignore",
				}
				err := baseLocal.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
