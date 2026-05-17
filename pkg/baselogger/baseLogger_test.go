package baselogger_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baselogger"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseLogger Service", func() {
	var (
		tempDir string
		service *baselogger.BaseLogger
		cfg     *config.BaseLoggerConfig
	)

	BeforeEach(func() {
		var err error

		tempDir, err = os.MkdirTemp("", "baselogger-test-*")
		Expect(err).NotTo(HaveOccurred())

		cfg = &config.BaseLoggerConfig{
			SourcePath:       tempDir,
			LoggerType:       goboottypes.LoggerTypeZerolog,
			ProjectName:      "testproject",
			LowerProjectName: "testproject",
			RepoPath:         "github.com/example/testproject",
		}

		service = baselogger.NewBaseLogger(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the expected service id", func() {
			Expect(service.ID()).To(Equal(goboottypes.ServiceNameBaseLogger))
		})
	})

	Describe("SetConfig", func() {
		It("accepts valid base_logger config", func() {
			Expect(service.SetConfig(cfg)).To(Succeed())
		})

		It("returns an error for invalid config type", func() {
			err := service.SetConfig(&config.BaseProjectConfig{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid config type"))
		})
	})

	Describe("Run", func() {
		It("does not render files directly", func() {
			cfg.LoggerType = goboottypes.LoggerTypeSlog

			Expect(service.SetConfig(cfg)).To(Succeed())
			Expect(service.Run()).To(Succeed())
			Expect(filepath.Join(tempDir, cfg.ProjectName)).NotTo(BeAnExistingFile())
		})
	})
})
