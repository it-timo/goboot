package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

const releaseZipFormat = "zip"

var _ = Describe("BaseReleaseConfig", func() {
	var baseRelease *config.BaseReleaseConfig

	BeforeEach(func() {
		baseRelease = &config.BaseReleaseConfig{
			SourcePath:  "./templates/release_base",
			ProjectName: testProjectName,
			GitProvider: goboottypes.GitProviderGitHub,
		}
	})

	Describe("When everything works as expected", func() {
		It("returns its stable identifier", func() {
			Expect(baseRelease.ID()).To(Equal(goboottypes.ServiceNameBaseRelease))
		})

		It("fills deterministic defaults", func() {
			Expect(baseRelease.Validate()).To(Succeed())
			Expect(baseRelease.BinaryName).To(Equal("testproject"))
			Expect(baseRelease.MainPackage).To(Equal("./cmd/testproject"))
			Expect(baseRelease.Formats).To(Equal([]string{"tar.gz", releaseZipFormat}))
		})

		It("accepts GitLab and explicit release settings", func() {
			baseRelease.GitProvider = goboottypes.GitProviderGitLab
			baseRelease.BinaryName = "custom-cli"
			baseRelease.MainPackage = "./cmd/custom"
			baseRelease.Formats = []string{releaseZipFormat}
			Expect(baseRelease.Validate()).To(Succeed())
		})
	})

	Describe("When the setup is incorrect", func() {
		It("rejects an unsupported provider", func() {
			baseRelease.GitProvider = unsupportedProvider
			Expect(baseRelease.Validate()).NotTo(Succeed())
		})

		It("rejects unsafe package input", func() {
			baseRelease.MainPackage = "./cmd/app; echo bad"
			Expect(baseRelease.Validate()).NotTo(Succeed())
		})

		It("rejects unsupported or duplicate formats", func() {
			baseRelease.Formats = []string{releaseZipFormat, releaseZipFormat}
			Expect(baseRelease.Validate()).NotTo(Succeed())
			baseRelease.Formats = []string{"deb"}
			Expect(baseRelease.Validate()).NotTo(Succeed())
		})
	})
})
