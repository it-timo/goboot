package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseTestConfig", func() {
	var (
		baseTest *config.BaseTestConfig
		tempDir  string
		testPath = "github.com/test/testproject"
	)

	BeforeEach(func() {
		baseTest = &config.BaseTestConfig{
			SourcePath:     "./templates/test_base",
			ProjectName:    "testproject",
			RepoImportPath: testPath,
			UseStyle:       goboottypes.TestStyleGinkgo,
		}

		var err error
		tempDir, err = os.MkdirTemp("", "basetest-config-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseTest.ID()).To(Equal(goboottypes.ServiceNameBaseTest))
			Expect(baseTest.ID()).To(Equal("base_test"))
		})
	})

	Describe("Validate", func() {
		Context("with valid config", func() {
			It("validates successfully", func() {
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("fills in derived fields", func() {
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseTest.CapsProjectName).To(Equal("TESTPROJECT"))
				Expect(baseTest.LowerProjectName).To(Equal("testproject"))
			})

			It("accepts ginkgo test style", func() {
				baseTest.UseStyle = goboottypes.TestStyleGinkgo
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("accepts go test style", func() {
				baseTest.UseStyle = goboottypes.TestStyleGo
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with missing required fields", func() {
			It("errors when sourcePath is missing", func() {
				baseTest.SourcePath = ""
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("errors when projectName is missing", func() {
				baseTest.ProjectName = ""
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			It("errors when repoImportPath is missing", func() {
				baseTest.RepoImportPath = ""
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("repoImportPath"))
			})

			It("errors when useStyle is missing", func() {
				baseTest.UseStyle = ""
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("useStyle"))
			})

			It("errors when multiple fields are missing", func() {
				baseTest.SourcePath = ""
				baseTest.ProjectName = ""
				baseTest.UseStyle = ""
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
				Expect(err.Error()).To(ContainSubstring("projectName"))
				Expect(err.Error()).To(ContainSubstring("useStyle"))
			})
		})

		Context("with whitespace-only values", func() {
			It("treats whitespace-only sourcePath as missing", func() {
				baseTest.SourcePath = "    "
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("treats whitespace-only projectName as missing", func() {
				baseTest.ProjectName = "\t\n"
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			It("treats whitespace-only repoImportPath as missing", func() {
				baseTest.RepoImportPath = "   "
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("repoImportPath"))
			})

			It("treats whitespace-only useStyle as missing", func() {
				baseTest.UseStyle = "  \t  "
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("useStyle"))
			})
		})

		Context("with invalid useStyle", func() {
			It("errors for unknown test style", func() {
				baseTest.UseStyle = "pytest"
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("useStyle must be"))
				Expect(err.Error()).To(ContainSubstring(goboottypes.TestStyleGinkgo))
				Expect(err.Error()).To(ContainSubstring(goboottypes.TestStyleGo))
			})

			It("errors for empty string after trimming", func() {
				baseTest.UseStyle = "   "
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
			})

			It("errors for case-sensitive mismatch", func() {
				baseTest.UseStyle = "Ginkgo" // Should be "ginkgo"
				err := baseTest.Validate()
				Expect(err).To(HaveOccurred())
			})
		})

		//nolint:dupl // Accept duplicated logic for clarity and extensibility.
		Context("when filling derived fields", func() {
			It("converts projectName to uppercase for CapsProjectName", func() {
				baseTest.ProjectName = "myAwesomeProject"
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseTest.CapsProjectName).To(Equal("MYAWESOMEPROJECT"))
			})

			It("converts projectName to lowercase for LowerProjectName", func() {
				baseTest.ProjectName = "MyAwesomeProject"
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseTest.LowerProjectName).To(Equal("myawesomeproject"))
			})

			It("handles already lowercase projectName", func() {
				baseTest.ProjectName = "alreadylower"
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseTest.LowerProjectName).To(Equal("alreadylower"))
				Expect(baseTest.CapsProjectName).To(Equal("ALREADYLOWER"))
			})

			It("handles already uppercase projectName", func() {
				baseTest.ProjectName = "ALREADYUPPER"
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseTest.CapsProjectName).To(Equal("ALREADYUPPER"))
				Expect(baseTest.LowerProjectName).To(Equal("alreadyupper"))
			})
		})
	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_test.yml")
		})

		Context("with valid YAML file", func() {
			It("loads the configuration successfully with ginkgo style", func() {
				yamlContent := `sourcePath: ./templates/test
useStyle: ginkgo
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseTestConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.SourcePath).To(Equal("./templates/test"))
				Expect(newConfig.UseStyle).To(Equal(goboottypes.TestStyleGinkgo))
				Expect(newConfig.RepoImportPath).To(Equal(testPath))
				Expect(newConfig.ProjectName).To(Equal("")) // yaml:"-" field
			})

			It("loads the configuration successfully with go style", func() {
				yamlContent := `sourcePath: ./templates/test_go
useStyle: go
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseTestConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.SourcePath).To(Equal("./templates/test_go"))
				Expect(newConfig.UseStyle).To(Equal(goboottypes.TestStyleGo))
				Expect(newConfig.RepoImportPath).To(Equal(testPath))
			})

			It("populates RepoImportPath from parameter", func() {
				yamlContent := `sourcePath: ./templates/test
useStyle: ginkgo
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				customPath := "github.com/custom/repo"
				newConfig := &config.BaseTestConfig{}
				err = newConfig.ReadConfig(configPath, customPath)
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.RepoImportPath).To(Equal(customPath))
			})
		})

		Context("with non-existent file", func() {
			It("returns an error", func() {
				err := baseTest.ReadConfig("/nonexistent/path.yml", testPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid YAML", func() {
			It("returns an error", func() {
				invalidYAML := `sourcePath: ./templates
useStyle: [invalid yaml structure
`
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseTestConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Real-world scenarios", func() {
		Context("when setting up Ginkgo test suite", func() {
			It("provides all necessary fields for template rendering", func() {
				baseTest.UseStyle = goboottypes.TestStyleGinkgo
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				// Check that all template-required fields are set
				Expect(baseTest.ProjectName).NotTo(BeEmpty())
				Expect(baseTest.LowerProjectName).NotTo(BeEmpty())
				Expect(baseTest.CapsProjectName).NotTo(BeEmpty())
				Expect(baseTest.RepoImportPath).NotTo(BeEmpty())
				Expect(baseTest.UseStyle).To(Equal(goboottypes.TestStyleGinkgo))
			})
		})

		Context("when setting up standard Go test suite", func() {
			It("provides all necessary fields for template rendering", func() {
				baseTest.UseStyle = goboottypes.TestStyleGo
				err := baseTest.Validate()
				Expect(err).NotTo(HaveOccurred())

				// Check that all template-required fields are set
				Expect(baseTest.ProjectName).NotTo(BeEmpty())
				Expect(baseTest.LowerProjectName).NotTo(BeEmpty())
				Expect(baseTest.CapsProjectName).NotTo(BeEmpty())
				Expect(baseTest.RepoImportPath).NotTo(BeEmpty())
				Expect(baseTest.UseStyle).To(Equal(goboottypes.TestStyleGo))
			})
		})
	})
})
