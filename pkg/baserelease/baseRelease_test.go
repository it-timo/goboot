package baserelease_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baserelease"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

type releaseRegistrar struct {
	lines map[string][]string
	files map[string][]string
}

func (r *releaseRegistrar) RegisterLines(name string, lines []string) error {
	r.lines[name] = lines

	return nil
}

func (r *releaseRegistrar) RegisterFile(name string, lines []string) error {
	r.files[name] = lines

	return nil
}

var _ = Describe("BaseRelease", func() {
	var (
		targetDir   string
		sourceDir   string
		service     *baserelease.BaseRelease
		releaseCfg  *config.BaseReleaseConfig
		ciRegistrar *releaseRegistrar
	)

	BeforeEach(func() {
		var err error

		targetDir, err = os.MkdirTemp("", "baserelease-target-*")
		Expect(err).NotTo(HaveOccurred())
		sourceDir, err = os.MkdirTemp("", "baserelease-source-*")
		Expect(err).NotTo(HaveOccurred())

		releaseCfg = &config.BaseReleaseConfig{
			SourcePath:  sourceDir,
			ProjectName: "IntroProject",
			GitProvider: goboottypes.GitProviderGitHub,
			BinaryName:  "introproject",
			MainPackage: "./cmd/introproject",
			Formats:     []string{"tar.gz", "zip"},
		}
		service = baserelease.NewBaseRelease(targetDir)
		ciRegistrar = &releaseRegistrar{lines: map[string][]string{}, files: map[string][]string{}}

		Expect(os.WriteFile(
			filepath.Join(sourceDir, ".goreleaser.yml"+goboottypes.TemplateSuffix),
			[]byte("project_name: {{.ProjectName}}\n"),
			0o644,
		)).To(Succeed())
		Expect(os.WriteFile(
			filepath.Join(sourceDir, "RELEASE.md"+goboottypes.TemplateSuffix),
			[]byte("# Release {{.ProjectName}}\n"),
			0o644,
		)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(targetDir)).To(Succeed())
		Expect(os.RemoveAll(sourceDir)).To(Succeed())
	})

	Describe("When everything works as expected", func() {
		It("renders files and registers the release workflow", func() {
			service.SetCIReceiver(ciRegistrar)
			Expect(service.SetConfig(releaseCfg)).To(Succeed())
			Expect(service.Run()).To(Succeed())

			projectDir := filepath.Join(targetDir, "IntroProject")
			content, err := os.ReadFile(filepath.Join(projectDir, ".goreleaser.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("IntroProject"))
			Expect(filepath.Join(projectDir, "RELEASE.md")).To(BeAnExistingFile())
			Expect(ciRegistrar.files).To(HaveKey(goboottypes.CIFileRelease))
		})

		It("returns its stable identifier", func() {
			Expect(service.ID()).To(Equal(goboottypes.ServiceNameBaseRelease))
		})
	})

	Describe("When the setup is incorrect", func() {
		It("rejects another service config type", func() {
			Expect(service.SetConfig(&config.BaseDockerConfig{})).NotTo(Succeed())
		})

		It("reports a missing release template", func() {
			Expect(os.Remove(filepath.Join(sourceDir, "RELEASE.md"+goboottypes.TemplateSuffix))).To(Succeed())
			Expect(service.SetConfig(releaseCfg)).To(Succeed())
			Expect(service.Run()).NotTo(Succeed())
		})
	})
})
