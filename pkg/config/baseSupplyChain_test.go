package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

const supplyChainMITLicense = "MIT"

var _ = Describe("BaseSupplyChainConfig", func() {
	var supplyChain *config.BaseSupplyChainConfig

	BeforeEach(func() {
		supplyChain = &config.BaseSupplyChainConfig{
			SourcePath:  "./templates/supplychain_base",
			ProjectName: testProjectName,
			GitProvider: goboottypes.GitProviderGitHub,
		}
	})

	It("fills deterministic defaults", func() {
		Expect(supplyChain.Validate()).To(Succeed())
		Expect(supplyChain.ID()).To(Equal(goboottypes.ServiceNameBaseSupplyChain))
		Expect(supplyChain.GovulncheckVersion).To(Equal("v1.1.4"))
		Expect(supplyChain.GoLicensesVersion).To(Equal("v2.0.1"))
		Expect(supplyChain.AllowedLicenses).To(ContainElement("Apache-2.0"))
	})

	It("accepts GitLab and explicit policy", func() {
		supplyChain.GitProvider = goboottypes.GitProviderGitLab
		supplyChain.GovulncheckVersion = "v1.2.3"
		supplyChain.GoLicensesVersion = "v2.3.4"
		supplyChain.AllowedLicenses = []string{supplyChainMITLicense}
		Expect(supplyChain.Validate()).To(Succeed())
	})

	It("rejects unsafe versions and license identifiers", func() {
		supplyChain.GovulncheckVersion = "latest"
		Expect(supplyChain.Validate()).NotTo(Succeed())

		supplyChain.GovulncheckVersion = "v1.1.4"
		supplyChain.AllowedLicenses = []string{"MIT; echo bad"}
		Expect(supplyChain.Validate()).NotTo(Succeed())
	})

	It("rejects duplicate licenses and unsupported providers", func() {
		supplyChain.AllowedLicenses = []string{supplyChainMITLicense, supplyChainMITLicense}
		Expect(supplyChain.Validate()).NotTo(Succeed())

		supplyChain.AllowedLicenses = []string{supplyChainMITLicense}
		supplyChain.GitProvider = unsupportedProvider
		Expect(supplyChain.Validate()).NotTo(Succeed())
	})
})
