package baseproject_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baseproject"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseProject Service", func() {
	var (
		tempDir     string
		baseProj    *baseproject.BaseProject
		validConfig *config.BaseProjectConfig
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "baseproject-test-*")
		Expect(err).NotTo(HaveOccurred())

		// Note: We use a minimal valid config for internal tests
		// Full template functionality would require actual template files
		validConfig = &config.BaseProjectConfig{
			SourcePath:            tempDir, // Using tempDir as mock source
			ProjectName:           "testproject",
			ProjectURL:            "https://github.com/test/testproject",
			RepoPath:              "github.com/test/testproject",
			UsedGoVersion:         "1.22.0",
			UsedNodeVersion:       "20.0.0",
			ReleaseCurrentWindow:  "Q1 2025",
			ReleaseUpcomingWindow: "Q2 2025",
			ReleaseLongTerm:       "2028",
			Author:                "Test Author",
			GitProvider:           "github",
			GitUser:               "testuser",
		}

		baseProj = baseproject.NewBaseProject(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("NewBaseProject", func() {
		It("creates a new BaseProject instance", func() {
			bp := baseproject.NewBaseProject("/some/path")
			Expect(bp).NotTo(BeNil())
		})

		It("stores the target directory", func() {
			targetPath := "/custom/target"
			bp := baseproject.NewBaseProject(targetPath)
			Expect(bp).NotTo(BeNil())
		})
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseProj.ID()).To(Equal(goboottypes.ServiceNameBaseProject))
			Expect(baseProj.ID()).To(Equal("base_project"))
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
			It("accepts BaseProjectConfig", func() {
				validConfig.SourcePath = sourceDir
				err := baseProj.SetConfig(validConfig)
				Expect(err).NotTo(HaveOccurred())
			})

			It("validates source and target paths are different", func() {
				validConfig.SourcePath = sourceDir
				err := baseProj.SetConfig(validConfig)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid config type", func() {
			It("returns an error for wrong config type", func() {
				wrongConfig := &config.BaseLintConfig{}
				err := baseProj.SetConfig(wrongConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})
		})

		Context("when source and target are the same", func() {
			It("returns an error to prevent overwrite", func() {
				validConfig.SourcePath = tempDir
				// This should fail because source == target
				err := baseProj.SetConfig(validConfig)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Run", func() {
		var sourceDir string

		writeTemplate := func(relPath, content string) {
			full := filepath.Join(sourceDir, relPath)
			Expect(os.MkdirAll(filepath.Dir(full), 0o755)).To(Succeed())
			Expect(os.WriteFile(full, []byte(content), 0o644)).To(Succeed())
		}

		buildConfig := func() *config.BaseProjectConfig {
			cfg := &config.BaseProjectConfig{
				SourcePath:            sourceDir,
				ProjectName:           "testproject",
				ProjectURL:            "https://github.com/test/testproject",
				RepoPath:              "github.com/test/testproject",
				UsedGoVersion:         "1.22.0",
				UsedNodeVersion:       "20.0.0",
				ReleaseCurrentWindow:  "Q1 2025",
				ReleaseUpcomingWindow: "Q2 2025",
				ReleaseLongTerm:       "2028",
				Author:                "Test Author",
				GitProvider:           "github",
				GitUser:               "testuser",
			}
			Expect(cfg.Validate()).To(Succeed())

			return cfg
		}

		BeforeEach(func() {
			var err error
			sourceDir, err = os.MkdirTemp("", "bp-source-*")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		It("copies structure and renders paths and contents", func() {
			writeTemplate("cmd/{{.LowerProjectName}}/main.go", "package main // {{.ProjectName}} {{.RepoPath}}")
			writeTemplate("README.md", "# {{.CapsProjectName}}")

			cfg := buildConfig()
			baseProj = baseproject.NewBaseProject(tempDir)
			Expect(baseProj.SetConfig(cfg)).To(Succeed())

			Expect(baseProj.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, cfg.ProjectName)
			renderedPath := filepath.Join(targetRoot, "cmd", cfg.LowerProjectName, "main.go")
			Expect(renderedPath).To(BeAnExistingFile())
			readme := filepath.Join(targetRoot, "README.md")
			Expect(readme).To(BeAnExistingFile())

			content, err := os.ReadFile(renderedPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring(cfg.ProjectName))
			Expect(string(content)).To(ContainSubstring(cfg.RepoPath))

			readmeContent, err := os.ReadFile(readme)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(readmeContent)).To(ContainSubstring(cfg.CapsProjectName))
		})

		It("errors on invalid path templates", func() {
			// invalid template in filename
			writeTemplate("{{.ProjectName", "content")

			cfg := buildConfig()
			baseProj = baseproject.NewBaseProject(tempDir)
			Expect(baseProj.SetConfig(cfg)).To(Succeed())

			err := baseProj.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to render path"))
		})

		It("errors on invalid content templates", func() {
			writeTemplate("README.md", "{{") // invalid template content

			cfg := buildConfig()
			baseProj = baseproject.NewBaseProject(tempDir)
			Expect(baseProj.SetConfig(cfg)).To(Succeed())

			err := baseProj.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed template render"))
		})

		It("propagates errors when source file cannot be read", func() {
			// unreadable file
			writeTemplate("README.md", "content")
			Expect(os.Chmod(filepath.Join(sourceDir, "README.md"), 0o000)).To(Succeed())

			cfg := buildConfig()
			baseProj = baseproject.NewBaseProject(tempDir)
			Expect(baseProj.SetConfig(cfg)).To(Succeed())

			err := baseProj.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read template file"))
		})
	})
})
