package config_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseProjectConfig", func() {
	var (
		baseProject *config.BaseProjectConfig
		tempDir     string
		testPath    = "https://github.com/user/testproject"
	)

	BeforeEach(func() {
		baseProject = &config.BaseProjectConfig{
			SourcePath:            "./templates/project_base",
			ProjectURL:            testPath,
			RepoPath:              "github.com/user/testproject",
			ProjectName:           "testproject",
			UsedGoVersion:         "1.22.0",
			UsedNodeVersion:       "20.11.0",
			ReleaseCurrentWindow:  "Q1 2025",
			ReleaseUpcomingWindow: "Q3 2025",
			ReleaseLongTerm:       "2028",
			Author:                "Test Author",
			GitProvider:           "github",
			GitUser:               "testuser",
		}

		var err error
		tempDir, err = os.MkdirTemp("", "baseproject-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseProject.ID()).To(Equal(goboottypes.ServiceNameBaseProject))
			Expect(baseProject.ID()).To(Equal("base_project"))
		})
	})

	Describe("Validate", func() {
		Context("with valid config", func() {
			It("validates successfully", func() {
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("fills in derived fields", func() {
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.CapsProjectName).To(Equal("TESTPROJECT"))
				Expect(baseProject.LowerProjectName).To(Equal("testproject"))
			})

			It("autofills current year when not set", func() {
				baseProject.CurrentYear = 0
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.CurrentYear).To(Equal(time.Now().Year()))
			})

			It("preserves manually set current year", func() {
				baseProject.CurrentYear = 2023
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.CurrentYear).To(Equal(2023))
			})
		})

		Context("with missing required fields", func() {
			DescribeTable("errors when field is missing",
				func(setupFn func(*config.BaseProjectConfig), expectedField string) {
					setupFn(baseProject)
					err := baseProject.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedField))
				},
				Entry("sourcePath", func(bp *config.BaseProjectConfig) { bp.SourcePath = "" },
					"sourcePath"),
				Entry("projectUrl", func(bp *config.BaseProjectConfig) { bp.ProjectURL = "" },
					"projectUrl"),
				Entry("projectName", func(bp *config.BaseProjectConfig) { bp.ProjectName = "" },
					"projectName"),
				Entry("usedGoVersion", func(bp *config.BaseProjectConfig) { bp.UsedGoVersion = "" },
					"usedGoVersion"),
				Entry("usedNodeVersion", func(bp *config.BaseProjectConfig) { bp.UsedNodeVersion = "" },
					"usedNodeVersion"),
				Entry("releaseCurrentWindow", func(bp *config.BaseProjectConfig) { bp.ReleaseCurrentWindow = "" },
					"releaseCurrentWindow"),
				Entry("releaseUpcomingWindow", func(bp *config.BaseProjectConfig) { bp.ReleaseUpcomingWindow = "" },
					"releaseUpcomingWindow"),
				Entry("releaseLongTerm", func(bp *config.BaseProjectConfig) { bp.ReleaseLongTerm = "" },
					"releaseLongTerm"),
				Entry("author", func(bp *config.BaseProjectConfig) { bp.Author = "" }, "author"),
			)

			It("errors when multiple fields are missing", func() {
				baseProject.SourcePath = ""
				baseProject.ProjectName = ""
				baseProject.Author = ""
				err := baseProject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
				Expect(err.Error()).To(ContainSubstring("projectName"))
				Expect(err.Error()).To(ContainSubstring("author"))
			})
		})

		Context("with GitProvider validation", func() {
			BeforeEach(func() {
				baseProject.GitProvider = "github"
			})

			It("validates successfully with GitUser set", func() {
				baseProject.GitUser = "myuser"
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("errors when GitUser is missing", func() {
				baseProject.GitUser = ""
				err := baseProject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("gitUser"))
			})

			It("errors when GitUser is whitespace only", func() {
				baseProject.GitUser = "  \t\n"
				err := baseProject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("gitUser"))
			})
		})

		Context("with whitespace-only values", func() {
			It("treats whitespace-only sourcePath as missing", func() {
				baseProject.SourcePath = "     "
				err := baseProject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourcePath"))
			})

			It("treats whitespace-only projectName as missing", func() {
				baseProject.ProjectName = "\t\n  "
				err := baseProject.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
			})

			DescribeTable("treats whitespace-only fields as missing",
				func(setupFn func(*config.BaseProjectConfig), expectedField string) {
					setupFn(baseProject)
					err := baseProject.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedField))
				},
				Entry("usedGoVersion", func(bp *config.BaseProjectConfig) { bp.UsedGoVersion = "      " },
					"usedGoVersion"),
				Entry("usedNodeVersion", func(bp *config.BaseProjectConfig) { bp.UsedNodeVersion = "\n" },
					"usedNodeVersion"),
				Entry("releaseCurrentWindow", func(bp *config.BaseProjectConfig) { bp.ReleaseCurrentWindow = "\t" },
					"releaseCurrentWindow"),
				Entry("releaseUpcomingWindow", func(bp *config.BaseProjectConfig) { bp.ReleaseUpcomingWindow = "  " },
					"releaseUpcomingWindow"),
				Entry("releaseLongTerm", func(bp *config.BaseProjectConfig) { bp.ReleaseLongTerm = "\t " },
					"releaseLongTerm"),
				Entry("author", func(bp *config.BaseProjectConfig) { bp.Author = "\n" }, "author"),
			)
		})

		//nolint:dupl // Accept duplicated logic for clarity and extensibility.
		Context("when filling derived fields", func() {
			It("converts projectName to uppercase for CapsProjectName", func() {
				baseProject.ProjectName = "myAnotherProject"
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.CapsProjectName).To(Equal("MYANOTHERPROJECT"))
			})

			It("converts projectName to lowercase for LowerProjectName", func() {
				baseProject.ProjectName = "MyAnotherProject"
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.LowerProjectName).To(Equal("myanotherproject"))
			})

			It("handles already lowercase projectName", func() {
				baseProject.ProjectName = "alreadyanotherlower"
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.LowerProjectName).To(Equal("alreadyanotherlower"))
				Expect(baseProject.CapsProjectName).To(Equal("ALREADYANOTHERLOWER"))
			})

			It("handles already uppercase projectName", func() {
				baseProject.ProjectName = "ALREADYANOTHERUPPER"
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				Expect(baseProject.CapsProjectName).To(Equal("ALREADYANOTHERUPPER"))
				Expect(baseProject.LowerProjectName).To(Equal("alreadyanotherupper"))
			})
		})
	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_project.yml")
		})

		Context("with valid YAML file", func() {
			It("loads the configuration successfully", func() {
				yamlContent := `sourcePath: ./templates/project
usedGoVersion: "1.22.5"
usedNodeVersion: "20.12.0"
releaseCurrentWindow: Q2 2025
releaseUpcomingWindow: Q4 2025
releaseLongTerm: "2029"
author: John Doe
gitProvider: github
gitUser: johndoe
currentYear: 2024
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseProjectConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).NotTo(HaveOccurred())
				newConfig.ProjectName = "testproject"
				Expect(newConfig.Validate()).To(Succeed())

				Expect(newConfig.SourcePath).To(Equal("./templates/project"))
				Expect(newConfig.ProjectURL).To(Equal(testPath))
				Expect(newConfig.RepoPath).To(Equal("github.com/user/testproject"))
				Expect(newConfig.ProjectName).To(Equal("testproject"))
				Expect(newConfig.UsedGoVersion).To(Equal("1.22.5"))
				Expect(newConfig.UsedNodeVersion).To(Equal("20.12.0"))
				Expect(newConfig.ReleaseCurrentWindow).To(Equal("Q2 2025"))
				Expect(newConfig.ReleaseUpcomingWindow).To(Equal("Q4 2025"))
				Expect(newConfig.ReleaseLongTerm).To(Equal("2029"))
				Expect(newConfig.Author).To(Equal("John Doe"))
				Expect(newConfig.GitProvider).To(Equal("github"))
				Expect(newConfig.GitUser).To(Equal("johndoe"))
				Expect(newConfig.CurrentYear).To(Equal(2024))
			})

			It("loads Git configuration", func() {
				yamlContent := `sourcePath: ./templates
projectUrl: https://gitlab.com/group/project
repoPath: gitlab.com/group/project
projectName: project
usedGoVersion: "1.21.0"
usedNodeVersion: "18.0.0"
releaseCurrentWindow: Q1 2025
releaseUpcomingWindow: Q2 2025
releaseLongTerm: "2027"
author: Jane Smith
gitProvider: gitlab
gitUser: janesmith
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseProjectConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).NotTo(HaveOccurred())
				newConfig.ProjectName = "project"
				Expect(newConfig.Validate()).To(Succeed())

				Expect(newConfig.GitProvider).To(Equal("gitlab"))
				Expect(newConfig.GitUser).To(Equal("janesmith"))
			})
		})

		Context("with non-existent file", func() {
			It("returns an error", func() {
				err := baseProject.ReadConfig("/nonexistent/path.yml", testPath)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with invalid YAML", func() {
			It("returns an error", func() {
				invalidYAML := `sourcePath: ./templates
projectName: [invalid
usedGoVersion: not closed properly
`
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				newConfig := &config.BaseProjectConfig{}
				err = newConfig.ReadConfig(configPath, testPath)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Real-world scenarios", func() {
		Context("when setting up a new Go project", func() {
			It("provides all necessary fields for template rendering", func() {
				err := baseProject.Validate()
				Expect(err).NotTo(HaveOccurred())

				// Check that all template-required fields are set
				Expect(baseProject.ProjectName).NotTo(BeEmpty())
				Expect(baseProject.LowerProjectName).NotTo(BeEmpty())
				Expect(baseProject.CapsProjectName).NotTo(BeEmpty())
				Expect(baseProject.ProjectURL).NotTo(BeEmpty())
				Expect(baseProject.RepoPath).NotTo(BeEmpty())
				Expect(baseProject.UsedGoVersion).NotTo(BeEmpty())
				Expect(baseProject.CurrentYear).To(BeNumerically(">", 2020))
				Expect(baseProject.Author).NotTo(BeEmpty())
			})
		})
	})
})
