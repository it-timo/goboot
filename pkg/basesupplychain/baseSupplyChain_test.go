package basesupplychain_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/basesupplychain"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

type supplyChainRegistrar struct {
	lines map[string][]string
	files map[string][]string
}

func (r *supplyChainRegistrar) RegisterLines(name string, lines []string) error {
	r.lines[name] = lines

	return nil
}

func (r *supplyChainRegistrar) RegisterFile(name string, lines []string) error {
	r.files[name] = lines

	return nil
}

var _ = Describe("BaseSupplyChain", func() {
	var (
		targetDir string
		sourceDir string
		service   *basesupplychain.BaseSupplyChain
		policyCfg *config.BaseSupplyChainConfig
		registrar *supplyChainRegistrar
	)

	BeforeEach(func() {
		var err error

		targetDir, err = os.MkdirTemp("", "basesupplychain-target-*")
		Expect(err).NotTo(HaveOccurred())
		sourceDir, err = os.MkdirTemp("", "basesupplychain-source-*")
		Expect(err).NotTo(HaveOccurred())

		policyCfg = &config.BaseSupplyChainConfig{
			SourcePath:         sourceDir,
			ProjectName:        "IntroProject",
			GitProvider:        goboottypes.GitProviderGitHub,
			GovulncheckVersion: "v1.1.4",
			GoLicensesVersion:  "v2.0.1",
			AllowedLicenses:    []string{"Apache-2.0", "MIT"},
		}
		service = basesupplychain.NewBaseSupplyChain(targetDir)
		registrar = &supplyChainRegistrar{lines: map[string][]string{}, files: map[string][]string{}}

		Expect(os.WriteFile(
			filepath.Join(sourceDir, "SUPPLY_CHAIN.md"+goboottypes.TemplateSuffix),
			[]byte("# {{.ProjectName}} supply chain\n"),
			0o644,
		)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(targetDir)).To(Succeed())
		Expect(os.RemoveAll(sourceDir)).To(Succeed())
	})

	It("renders policy and registers pinned security checks", func() {
		service.SetCIReceiver(registrar)
		Expect(service.SetConfig(policyCfg)).To(Succeed())
		Expect(service.Run()).To(Succeed())

		policyPath := filepath.Join(targetDir, "IntroProject", "SUPPLY_CHAIN.md")
		content, err := os.ReadFile(policyPath)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring("IntroProject supply chain"))
		Expect(registrar.files).To(HaveKey(goboottypes.CIFileSecurity))
		Expect(registrar.files[goboottypes.CIFileSecurity]).To(ContainElement(
			"go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...",
		))
	})

	It("returns its stable identifier", func() {
		Expect(service.ID()).To(Equal(goboottypes.ServiceNameBaseSupplyChain))
	})

	It("rejects another service config type", func() {
		Expect(service.SetConfig(&config.BaseDockerConfig{})).NotTo(Succeed())
	})
})
