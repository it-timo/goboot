package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseLogger Configuration", func() {
	Describe("ID", func() {
		It("returns base_logger identifier", func() {
			cfg := &config.BaseLoggerConfig{}
			Expect(cfg.ID()).To(Equal(goboottypes.ServiceNameBaseLogger))
		})
	})

	Describe("ReadConfig and Validate", func() {
		var (
			tempDir    string
			configPath string
			baseLogger *config.BaseLoggerConfig
		)

		BeforeEach(func() {
			var err error

			tempDir, err = os.MkdirTemp("", "base-logger-config-*")
			Expect(err).NotTo(HaveOccurred())

			configPath = filepath.Join(tempDir, "base_logger.yml")
			baseLogger = &config.BaseLoggerConfig{
				ProjectName: "MyProject",
			}
		})

		AfterEach(func() {
			if tempDir != "" {
				Expect(os.RemoveAll(tempDir)).To(Succeed())
			}
		})

		It("accepts zerolog loggerType", func() {
			yamlContent, err := loadTestFixture("config/base_logger/valid_zerolog.yml")
			Expect(err).NotTo(HaveOccurred())
			Expect(os.WriteFile(configPath, yamlContent, 0o644)).To(Succeed())

			Expect(baseLogger.ReadConfig(configPath, "https://github.com/example/repo", "")).To(Succeed())
			Expect(baseLogger.Validate()).To(Succeed())
			Expect(baseLogger.LoggerType).To(Equal(goboottypes.LoggerTypeZerolog))
			Expect(baseLogger.RepoPath).To(Equal("github.com/example/repo"))
			Expect(baseLogger.LowerProjectName).To(Equal("myproject"))
		})

		It("accepts slog loggerType", func() {
			yamlContent, err := loadTestFixture("config/base_logger/valid_slog.yml")
			Expect(err).NotTo(HaveOccurred())
			Expect(os.WriteFile(configPath, yamlContent, 0o644)).To(Succeed())

			Expect(baseLogger.ReadConfig(configPath, "http://gitlab.com/example/repo", "")).To(Succeed())
			Expect(baseLogger.Validate()).To(Succeed())
			Expect(baseLogger.LoggerType).To(Equal(goboottypes.LoggerTypeSlog))
			Expect(baseLogger.RepoPath).To(Equal("gitlab.com/example/repo"))
		})

		It("returns error for invalid loggerType", func() {
			yamlContent, err := loadTestFixture("config/base_logger/invalid.yml")
			Expect(err).NotTo(HaveOccurred())
			Expect(os.WriteFile(configPath, yamlContent, 0o644)).To(Succeed())

			Expect(baseLogger.ReadConfig(configPath, "https://github.com/example/repo", "")).To(Succeed())
			err = baseLogger.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("loggerType must be"))
		})
	})
})
