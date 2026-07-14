package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

const governanceMaintainerTest = "@maintainer"

var _ = Describe("BaseGovernanceConfig", func() {
	var governanceConfig *config.BaseGovernanceConfig

	BeforeEach(func() {
		governanceConfig = &config.BaseGovernanceConfig{
			SourcePath:  "./templates/governance_base",
			Maintainers: []string{governanceMaintainerTest, "@example/team"},
			ProjectName: testProjectName,
			ProjectURL:  "https://github.com/example/testproject",
			GitProvider: goboottypes.GitProviderGitHub,
			Profile:     goboottypes.ProfileStandard,
		}
	})

	It("returns its stable identifier", func() {
		Expect(governanceConfig.ID()).To(Equal(goboottypes.ServiceNameBaseGovernance))
	})

	It("defaults the target branch to main", func() {
		Expect(governanceConfig.Validate()).To(Succeed())
		Expect(governanceConfig.DefaultBranch).To(Equal("main"))
	})

	It("accepts GitLab and an explicit target branch", func() {
		governanceConfig.GitProvider = goboottypes.GitProviderGitLab
		governanceConfig.DefaultBranch = "release/v1"
		Expect(governanceConfig.Validate()).To(Succeed())
	})

	It("rejects missing, malformed, or duplicate maintainers", func() {
		governanceConfig.Maintainers = nil
		Expect(governanceConfig.Validate()).NotTo(Succeed())

		governanceConfig.Maintainers = []string{"maintainer"}
		Expect(governanceConfig.Validate()).NotTo(Succeed())

		governanceConfig.Maintainers = []string{governanceMaintainerTest, governanceMaintainerTest}
		Expect(governanceConfig.Validate()).NotTo(Succeed())
	})

	It("rejects an unsupported provider or unsafe branch", func() {
		governanceConfig.GitProvider = "other"
		Expect(governanceConfig.Validate()).NotTo(Succeed())

		governanceConfig.GitProvider = goboottypes.GitProviderGitHub
		governanceConfig.DefaultBranch = "../main"
		Expect(governanceConfig.Validate()).NotTo(Succeed())
	})
})
