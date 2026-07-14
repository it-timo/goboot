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
			yamlContent, err := loadTestFixture("config/goboot/read_from_path.yml")
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(configPath, yamlContent, 0644)
			Expect(err).NotTo(HaveOccurred())

			gb := config.NewGoBoot(configPath)
			err = gb.Init()
			Expect(err).NotTo(HaveOccurred())
			Expect(gb.ProjectName).To(Equal("FromCustomPath"))
			Expect(gb.TargetPath).To(Equal("/tmp/from-custom"))
		})
	})

	Describe("parallelism", func() {
		writeRootConfig := func(value string) {
			yamlContent := "projectName: \"ParallelProject\"\n" +
				"targetPath: \"/tmp/parallel-project\"\n" +
				value +
				"services: []\n"

			Expect(os.WriteFile(configPath, []byte(yamlContent), 0o644)).To(Succeed())
		}

		It("defaults to serial execution when omitted", func() {
			writeRootConfig("")

			goBoot = config.NewGoBoot(configPath)
			Expect(goBoot.Init()).To(Succeed())
			Expect(goBoot.Parallelism).To(Equal(config.DefaultParallelism))
		})

		It("accepts bounded parallel execution", func() {
			writeRootConfig("parallelism: 4\n")

			goBoot = config.NewGoBoot(configPath)
			Expect(goBoot.Init()).To(Succeed())
			Expect(goBoot.Parallelism).To(Equal(4))
		})

		DescribeTable("rejects unsafe parallelism values",
			func(value string) {
				writeRootConfig("parallelism: " + value + "\n")

				goBoot = config.NewGoBoot(configPath)
				err := goBoot.Init()
				Expect(err).To(MatchError(ContainSubstring("parallelism must be between")))
			},
			Entry("zero", "0"),
			Entry("negative", "-1"),
			Entry("above the maximum", "33"),
		)
	})

	Describe("Init", func() {
		Context("with valid configuration", func() {
			BeforeEach(func() {
				// Create a minimal valid config file
				yamlContent, err := loadTestFixtureWithVars("config/goboot/init_valid.yml", map[string]string{
					"BASE_PROJECT_PATH": filepath.Join(tempDir, "base_project.yml"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create the service config file
				serviceConfigContent, err := loadTestFixture("config/goboot/service_base_project.yml")
				Expect(err).NotTo(HaveOccurred())

				serviceConfigPath := filepath.Join(tempDir, "base_project.yml")
				err = os.WriteFile(serviceConfigPath, serviceConfigContent, 0644)
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
				Expect(goBoot.ProjectName).To(Equal(testProjectName))
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
				invalidYAML, err := loadTestFixture("config/goboot/invalid_goboot.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, invalidYAML, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
			})

			It("returns error for unknown root YAML fields", func() {
				yamlContent := []byte("projectName: \"TestProject\"\n" +
					"targetPath: \"/tmp/test\"\n" +
					"repoUrl: \"https://github.com/test/testproject\"\n" +
					"gitProvider: \"github\"\n" +
					"projectNmae: \"typo\"\n" +
					"services: []\n")

				err := os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to read goboot config"))
				Expect(err.Error()).To(ContainSubstring("field projectNmae not found"))
			})

			It("returns error when required base fields are missing", func() {
				yamlContent, err := loadTestFixture("config/goboot/missing_required_fields.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("projectName"))
				Expect(err.Error()).To(ContainSubstring("targetPath"))
			})

			It("returns error when enabled service has no confPath", func() {
				yamlContent, err := loadTestFixture("config/goboot/enabled_service_missing_confpath.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("services[base_project].confPath"))
			})

			DescribeTable("returns error when projectName is not a safe Go scaffold identifier",
				func(projectName string) {
					yamlContent := []byte("projectName: \"" + projectName + "\"\n" +
						"targetPath: \"/tmp/test\"\n" +
						"repoUrl: \"https://github.com/test/testproject\"\n" +
						"gitProvider: \"github\"\n" +
						"services: []\n")

					err := os.WriteFile(configPath, yamlContent, 0644)
					Expect(err).NotTo(HaveOccurred())

					goBoot = config.NewGoBoot(configPath)
					err = goBoot.Init()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("projectName"))
				},
				Entry("path traversal", "../outside"),
				Entry("hyphen", "cli-project"),
				Entry("space", "my project"),
				Entry("leading digit", "123app"),
				Entry("punctuation", "app!"),
				Entry("leading whitespace", " app"),
			)
		})

		Context("with disabled services", func() {
			BeforeEach(func() {
				yamlContent, err := loadTestFixture("config/goboot/disabled_service.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
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
				yamlContent, err := loadTestFixture("config/goboot/unknown_service.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
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
				yamlContent, err := loadTestFixtureWithVars("config/goboot/invalid_service_ref.yml", map[string]string{
					"INVALID_PATH": filepath.Join(tempDir, "invalid.yml"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create invalid service config (missing required fields)
				invalidServiceConfig, err := loadTestFixture("config/goboot/invalid_service_config.yml")
				Expect(err).NotTo(HaveOccurred())

				invalidPath := filepath.Join(tempDir, "invalid.yml")
				err = os.WriteFile(invalidPath, invalidServiceConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
			})

			It("returns error when service config validation fails", func() {
				err := goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to register config"))
			})

			It("returns error when service config file cannot be read", func() {
				yamlContent, err := loadTestFixtureWithVars("config/goboot/missing_service_ref.yml", map[string]string{
					"MISSING_PATH": filepath.Join(tempDir, "missing.yml"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)

				err = goBoot.Init()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to read config"))
			})
		})

		Context("with multiple services", func() {
			BeforeEach(func() {
				yamlContent, err := loadTestFixtureWithVars("config/goboot/multiple_services.yml", map[string]string{
					"BASE_PROJECT_PATH": filepath.Join(tempDir, "base_project.yml"),
					"BASE_LINT_PATH":    filepath.Join(tempDir, "base_lint.yml"),
					"BASE_LOCAL_PATH":   filepath.Join(tempDir, "base_local.yml"),
					"BASE_TEST_PATH":    filepath.Join(tempDir, "base_test.yml"),
					"BASE_LOGGER_PATH":  filepath.Join(tempDir, "base_logger.yml"),
					"BASE_CI_PATH":      filepath.Join(tempDir, "base_ci.yml"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_project config
				projectConfig, err := loadTestFixture("config/goboot/multiple_base_project.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_project.yml"), projectConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_lint config
				lintConfig, err := loadTestFixture("config/goboot/multiple_base_lint.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_lint.yml"), lintConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_local config
				localConfig, err := loadTestFixture("config/goboot/multiple_base_local.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_local.yml"), localConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_test config
				testConfig, err := loadTestFixture("config/goboot/multiple_base_test.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_test.yml"), testConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_ci config
				ciConfig, err := loadTestFixture("config/goboot/multiple_base_ci.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_ci.yml"), ciConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				// Create base_logger config
				loggerConfig, err := loadTestFixture("config/goboot/multiple_base_logger.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "base_logger.yml"), loggerConfig, 0644)
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

				// Check base_test (service)
				_, exist = goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseTest)
				Expect(exist).To(BeTrue())

				// Check base_ci (service)
				_, exist = goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseCI)
				Expect(exist).To(BeTrue())

				// Check base_logger (service)
				_, exist = goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseLogger)
				Expect(exist).To(BeTrue())
			})
		})

		Context("with base_test service", func() {
			It("registers the test config and fills defaults", func() {
				baseTestPath := filepath.Join(tempDir, "base_test.yml")
				yamlContent, err := loadTestFixtureWithVars("config/goboot/base_test.yml", map[string]string{
					"BASE_TEST_PATH": baseTestPath,
					"TARGET_PATH":    filepath.Join(tempDir, "out"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				baseTestConfig, err := loadTestFixture("config/goboot/base_test_config.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(baseTestPath, baseTestConfig, 0644)
				Expect(err).NotTo(HaveOccurred())

				goBoot = config.NewGoBoot(configPath)
				Expect(goBoot.Init()).To(Succeed())

				rawCfg, ok := goBoot.ConfManager.GetService(goboottypes.ServiceNameBaseTest)
				Expect(ok).To(BeTrue())

				testCfg, ok := rawCfg.(*config.BaseTestConfig)
				Expect(ok).To(BeTrue())
				Expect(testCfg.ProjectName).To(Equal(testProjectName))
				Expect(testCfg.RepoImportPath).To(Equal("github.com/user/testproject"))
				Expect(testCfg.TestCMD).To(Equal(goboottypes.DefaultGoTestCMD))
				Expect(testCfg.UseStyle).To(Equal(goboottypes.TestStyleGinkgo))
			})
		})
	})

	Describe("Real-world scenarios", func() {
		Context("when setting up a complete project", func() {
			It("handles a typical configuration", func() {
				yamlContent, err := loadTestFixtureWithVars("config/goboot/real_world.yml", map[string]string{
					"PROJECT_PATH": filepath.Join(tempDir, "project.yml"),
				})
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(configPath, yamlContent, 0644)
				Expect(err).NotTo(HaveOccurred())

				projectConfig, err := loadTestFixture("config/goboot/real_world_project.yml")
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(tempDir, "project.yml"), projectConfig, 0644)
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
