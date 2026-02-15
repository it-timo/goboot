package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseCIConfig", func() {
	var (
		baseCI  *config.BaseCIConfig
		tempDir string
	)

	BeforeEach(func() {
		baseCI = &config.BaseCIConfig{
			SourcePath:  "./templates/ci_base",
			ProjectName: "testproject",
			GoVersion:   []string{"1.25", "1.26"},
			GitProvider: goboottypes.GitProviderGitLab,
			Jobs: map[string]*config.CIJob{
				"build": {
					Commands: []string{"go build ./..."},
				},
			},
		}

		var err error
		tempDir, err = os.MkdirTemp("", "baseci-config-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseCI.ID()).To(Equal(goboottypes.ServiceNameBaseCI))
			Expect(baseCI.ID()).To(Equal("base_ci"))
		})
	})

	Describe("Validate", func() {
		Context("with valid config", func() {
			It("validates successfully", func() {
				err := baseCI.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(baseCI.AutoBranches).To(Equal([]string{"main", "master"}))
				Expect(baseCI.ImagePolicy).To(Equal("balanced"))
			})

			It("validates successfully with no jobs configured", func() {
				baseCI.Jobs = nil
				err := baseCI.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with missing required fields", func() {
			It("errors when sourcePath is missing", func() {
				baseCI.SourcePath = ""
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			Context("with invalid optional fields", func() {
				It("errors when autoBranches contains blank entries", func() {
					baseCI.AutoBranches = []string{"main", " "}
					err := baseCI.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("autoBranches contains empty string"))
				})

				It("errors on unsupported image policy", func() {
					baseCI.ImagePolicy = "hardened"
					err := baseCI.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("imagePolicy"))
				})
			})

			It("errors when projectName is missing", func() {
				baseCI.ProjectName = ""
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			It("errors when gitProvider is missing", func() {
				baseCI.GitProvider = ""
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("gitProvider"))
			})

			It("errors when goVersions is missing", func() {
				baseCI.GoVersion = nil
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("goVersions"))
			})

			It("errors when goVersions contains blank entries", func() {
				baseCI.GoVersion = []string{"1.25", "   "}
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("goVersions contains empty string"))
			})
		})

		Context("with invalid jobs", func() {
			It("errors on invalid job name", func() {
				baseCI.Jobs = map[string]*config.CIJob{
					"weird": {
						Commands: []string{"echo hi"},
					},
				}
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid job name"))
			})

			It("errors on missing job commands", func() {
				baseCI.Jobs["build"].Commands = nil
				err := baseCI.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("missing commands"))
			})
		})

		Context("with valid jobs", func() {
			It("derives job file names from job keys", func() {
				err := baseCI.Validate()
				Expect(err).NotTo(HaveOccurred())
				Expect(baseCI.Jobs["build"].File).To(Equal("build.yml"))
			})
		})

	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_ci.yml")
		})

		Context("with valid YAML file", func() {
			It("loads the configuration successfully", func() {
				yamlContent, err := loadTestFixture("config/base_ci/valid.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseCIConfig{}
				err = newConfig.ReadConfig(configPath, "", goboottypes.GitProviderGitLab)
				Expect(err).NotTo(HaveOccurred())

				Expect(newConfig.SourcePath).To(Equal("./templates/ci_base"))
				Expect(newConfig.Jobs).To(HaveKey("build"))
				Expect(newConfig.GitProvider).To(Equal(goboottypes.GitProviderGitLab))
				Expect(newConfig.GoVersion).To(ContainElement("1.25"))
			})
		})

		Context("with non-existent file", func() {
			It("returns an error", func() {
				err := baseCI.ReadConfig("/nonexistent/path.yml", "", goboottypes.GitProviderGitLab)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid YAML", func() {
			It("returns an error", func() {
				invalidYAML, err := loadTestFixture("config/base_ci/invalid.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, invalidYAML, 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseCIConfig{}
				err = newConfig.ReadConfig(configPath, "", goboottypes.GitProviderGitLab)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
