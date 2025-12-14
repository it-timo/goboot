package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("GoBoot Configuration Orchestrator", func() {
	var (
		tempDir    string
		configPath string
		goBoot     *config.GoBoot
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "goboot-config-test-*")
		Expect(err).NotTo(HaveOccurred())

		configPath = filepath.Join(tempDir, "goboot.yml")
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("NewGoBoot", func() {
		It("creates a new GoBoot instance", func() {
			gb := config.NewGoBoot(configPath)
			Expect(gb).NotTo(BeNil())
		})

		It("initializes the ConfManager", func() {
			gb := config.NewGoBoot(configPath)
			Expect(gb.ConfManager).NotTo(BeNil())
		})

		It("reads from the provided config path during init", func() {
			yamlContent := `projectName: from-custom-path
targetPath: /tmp/from-custom
services: []
`
			err := os.WriteFile(configPath, []byte(yamlContent), 0644)
			Expect(err).NotTo(HaveOccurred())

			gb := config.NewGoBoot(configPath)
			err = gb.Init()
			Expect(err).NotTo(HaveOccurred())
			Expect(gb.ProjectName).To(Equal("from-custom-path"))
			Expect(gb.TargetPath).To(Equal("/tmp/from-custom"))
		})
	})

	Describe("Init", func() {
		Context("with valid configuration", func() {
			BeforeEach(func() {
				// Create a minimal valid config file
				yamlContent := `projectName: testproject
repoUrl: https://github.com/user/testproject
targetPath: /tmp/test
services:
  - id: base_project
    confPath: ` + filepath.Join(tempDir, "base_project.yml") + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create the service config file
				serviceConfigContent := `sourcePath: /tmp/templates
projectName: testproject
usedGoVersion: "1.22.0"
usedNodeVersion: "20.0.0"
releaseCurrentWindow: Q1 2025
releaseUpcomingWindow: Q2 2025
releaseLongTerm: "2028"
author: Test Author
gitProvider: github
gitUser: testuser
`
				serviceConfigPath := filepath.Join(tempDir, "base_project.yml")
				err = os.WriteFile(serviceConfigPath, []byte(serviceConfigContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("successfully initializes and loads config", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())
			})

			It("populates project name", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())
				Expect(goBoot.ProjectName).To(Equal("testproject"))
			})

			It("populates target path", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())
				Expect(goBoot.TargetPath).To(Equal("/tmp/test"))
			})

			It("loads and registers service configs", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())

				// Verify service was registered
				cfg, ok := goBoot.ConfManager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeTrue())
				Expect(cfg).NotTo(BeNil())
			})
		})

		Context("with invalid configuration file", func() {
			It("returns error for non-existent config", func() {
				goBoot = config.NewGoBoot("/nonexistent/config.yml")
				err := goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to read goboot config"))
			})

			It("returns error for malformed YAML", func() {
				invalidYAML := `projectName: [invalid yaml
services: not properly formatted
`
				err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
			})

			It("returns error when required base fields are missing", func() {
				yamlContent := `projectName: ""
targetPath: ""
services: []
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
				Expect(err.Error()).To(ContainSubstring("targetPath"))
			})

			It("returns error when enabled service has no confPath", func() {
				yamlContent := `projectName: testproject
targetPath: /tmp/test
services:
  - id: base_project
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("services[base_project].confPath"))
			})
		})

		Context("with disabled services", func() {
			BeforeEach(func() {
				yamlContent := `projectName: testproject
targetPath: /tmp/test
services:
  - id: base_project
    confPath: /tmp/base_project.yml
    enabled: false
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("skips loading disabled services", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())

				// Disabled service should not be registered
				_, ok := goBoot.ConfManager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(ok).To(BeFalse())
			})
		})

		Context("with unknown service ID", func() {
			BeforeEach(func() {
				yamlContent := `projectName: testproject
repoUrl: https://github.com/user/testproject
targetPath: /tmp/test
services:
  - id: unknown_service_xyz
    confPath: /tmp/config.yml
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("returns error for unknown service", func() {
				err := goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid or nil config"))
			})
		})

		Context("with invalid service config", func() {
			BeforeEach(func() {
				yamlContent := `projectName: testproject
repoUrl: https://github.com/user/testproject
targetPath: /tmp/test
services:
  - id: base_project
    confPath: ` + filepath.Join(tempDir, "invalid.yml") + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create invalid service config (missing required fields)
				invalidServiceConfig := `projectName: test
# Missing all required fields
`
				invalidPath := filepath.Join(tempDir, "invalid.yml")
				err = os.WriteFile(invalidPath, []byte(invalidServiceConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("returns error when service config validation fails", func() {
				err := goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to register config"))
			})

			It("returns error when service config file cannot be read", func() {
				yamlContent := `projectName: testproject
repoUrl: https://github.com/user/testproject
targetPath: /tmp/test
services:
  - id: base_project
    confPath: ` + filepath.Join(tempDir, "missing.yml") + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)

				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to read config"))
			})
		})

		Context("with multiple services", func() {
			BeforeEach(func() {
				yamlContent := `projectName: testproject
repoUrl: https://github.com/user/testproject
targetPath: /tmp/test
services:
  - id: base_project
    confPath: ` + filepath.Join(tempDir, "base_project.yml") + `
    enabled: true
  - id: base_lint
    confPath: ` + filepath.Join(tempDir, "base_lint.yml") + `
    enabled: true
  - id: base_local
    confPath: ` + filepath.Join(tempDir, "base_local.yml") + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_project config
				projectConfig := `sourcePath: /tmp/templates
projectName: testproject
usedGoVersion: "1.22.0"
usedNodeVersion: "20.0.0"
releaseCurrentWindow: Q1 2025
releaseUpcomingWindow: Q2 2025
releaseLongTerm: "2028"
author: Test Author
gitProvider: github
gitUser: testuser
`
				err = os.WriteFile(filepath.Join(tempDir, "base_project.yml"), []byte(projectConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_lint config
				lintConfig := `sourcePath: /tmp/lint
linters:
  golang:
    enabled: true
`
				err = os.WriteFile(filepath.Join(tempDir, "base_lint.yml"), []byte(lintConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_local config
				localConfig := `sourcePath: /tmp/local
fileList:
  - Makefile
  - Taskfile.yml
`
				err = os.WriteFile(filepath.Join(tempDir, "base_local.yml"), []byte(localConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("loads and registers all services", func() {
				err := goBoot.Init()
				Expect(err).NotTo(HaveOccurred())

				// Check base_project (registrar)
				_, exist := goBoot.ConfManager.GetRegistrar(goboottypes.ServiceNameBaseProject)
				Expect(exist).To(BeTrue())

				// Check base_lint (service)
				_, exist = goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseLint)
				Expect(exist).To(BeTrue())

				// Check base_local (service)
				_, exist = goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseLocal)
				Expect(exist).To(BeTrue())
			})
		})

		Context("with base_test service", func() {
			It("registers the test config and fills defaults", func() {
				baseTestPath := filepath.Join(tempDir, "base_test.yml")
				yamlContent := `projectName: testproject
repoUrl: github.com/user/testproject
targetPath: ` + filepath.Join(tempDir, "out") + `
services:
  - id: base_test
    confPath: ` + baseTestPath + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				baseTestConfig := `sourcePath: ./templates/test_base
useStyle: ginkgo
`
				err = os.WriteFile(baseTestPath, []byte(baseTestConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				Expect(goBoot.Init()).To(Succeed())

				rawCfg, ok := goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseTest)
				Expect(ok).To(BeTrue())

				testCfg, ok := rawCfg.(*config.BaseTestConfig)
				Expect(ok).To(BeTrue())
				Expect(testCfg.ProjectName).To(Equal("testproject"))
				Expect(testCfg.RepoImportPath).To(Equal("github.com/user/testproject"))
				Expect(testCfg.TestCMD).To(Equal(goboottypes.DefaultGoTestCMD))
				Expect(testCfg.UseStyle).To(Equal(goboottypes.TestStyleGinkgo))
			})
		})
	})

	Describe("Real-world scenarios", func() {
		Context("when setting up a complete project", func() {
			It("handles a typical configuration", func() {
				yamlContent := `projectName: myproject
repoUrl: https://github.com/user/myproject
targetPath: /tmp/myproject
services:
  - id: base_project
    confPath: ` + filepath.Join(tempDir, "project.yml") + `
    enabled: true
`
				err := os.WriteFile(configPath, []byte(yamlContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				projectConfig := `sourcePath: /tmp/templates
usedGoVersion: "1.22.5"
usedNodeVersion: "20.12.0"
releaseCurrentWindow: Q2 2025
releaseUpcomingWindow: Q4 2025
releaseLongTerm: "2029"
author: John Doe
gitProvider: github
gitUser: johndoe
`
				err = os.WriteFile(filepath.Join(tempDir, "project.yml"), []byte(projectConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				testGoBoot := config.NewGoBoot(configPath)
				err = testGoBoot.Init()
				Expect(err).NotTo(HaveOccurred())

				Expect(testGoBoot.ProjectName).To(Equal("myproject"))
				Expect(testGoBoot.Services).To(HaveLen(1))
				Expect(testGoBoot.Services[0].ID).To(Equal(goboottypes.ServiceNameBaseProject))
				Expect(testGoBoot.Services[0].IsEnabled()).To(BeTrue())
			})
		})
	})
})
