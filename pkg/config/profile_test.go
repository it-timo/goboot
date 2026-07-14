package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
)

var (
	_ goboottypes.ProfileReceiver = (*config.BaseProjectConfig)(nil)
	_ goboottypes.ProfileReceiver = (*config.BaseLintConfig)(nil)
	_ goboottypes.ProfileReceiver = (*config.BaseTestConfig)(nil)
)

var _ = Describe("Template profiles", func() {
	Describe("When everything works as expected", func() {
		It("applies profile-specific test defaults", func() {
			cases := []struct {
				profile string
				style   string
				command string
			}{
				{goboottypes.ProfileMinimal, goboottypes.TestStyleGo, "go test ./..."},
				{goboottypes.ProfileStandard, goboottypes.TestStyleGinkgo, "-race"},
				{goboottypes.ProfileEnterprise, goboottypes.TestStyleGinkgo, "-shuffle=on"},
				{goboottypes.ProfileOSS, goboottypes.TestStyleGinkgo, "-covermode=atomic"},
			}

			for _, curCase := range cases {
				testConfig := &config.BaseTestConfig{
					SourcePath:     "templates/test_base",
					ProjectName:    testProjectName,
					RepoImportPath: "github.com/test/testproject",
				}
				testConfig.SetProfile(curCase.profile)

				Expect(testConfig.Validate()).To(Succeed())
				Expect(testConfig.UseStyle).To(Equal(curCase.style))
				Expect(testConfig.TestCMD).To(ContainSubstring(curCase.command))
			}
		})

		It("applies profile-specific complexity thresholds", func() {
			expectedThresholds := map[string]int{
				goboottypes.ProfileMinimal:    20,
				goboottypes.ProfileStandard:   10,
				goboottypes.ProfileEnterprise: 8,
				goboottypes.ProfileOSS:        10,
			}

			for profile, expected := range expectedThresholds {
				lintConfig := &config.BaseLintConfig{
					SourcePath:     "templates/lint_base",
					ProjectName:    testProjectName,
					RepoImportPath: "github.com/test/testproject",
					Linters: map[string]*config.Linter{
						goboottypes.LinterGo: {Enabled: true},
					},
				}
				lintConfig.SetProfile(profile)

				Expect(lintConfig.Validate()).To(Succeed())
				Expect(lintConfig.CyclomaticComplexity).To(Equal(expected))
			}
		})

		It("preserves explicit test settings", func() {
			testConfig := &config.BaseTestConfig{
				SourcePath:     "templates/test_base",
				ProjectName:    testProjectName,
				RepoImportPath: "github.com/test/testproject",
				UseStyle:       goboottypes.TestStyleGo,
				TestCMD:        "go test -count=1 ./...",
			}
			testConfig.SetProfile(goboottypes.ProfileEnterprise)

			Expect(testConfig.Validate()).To(Succeed())
			Expect(testConfig.UseStyle).To(Equal(goboottypes.TestStyleGo))
			Expect(testConfig.TestCMD).To(Equal("go test -count=1 ./..."))
		})

		It("renders a smaller Go linter set for the minimal profile", func() {
			templatePath := filepath.Join("..", "..", "templates", "lint_base", ".golangci.yml.tmpl")
			templateContent, err := os.ReadFile(templatePath)
			Expect(err).NotTo(HaveOccurred())

			lintConfig := &config.BaseLintConfig{
				SourcePath:     "templates/lint_base",
				ProjectName:    testProjectName,
				RepoImportPath: "github.com/test/testproject",
				Linters: map[string]*config.Linter{
					goboottypes.LinterGo: {Enabled: true},
				},
			}
			lintConfig.SetProfile(goboottypes.ProfileMinimal)
			Expect(lintConfig.Validate()).To(Succeed())

			rendered, err := gobootutils.ExecuteTemplateText("minimal_lint", string(templateContent), lintConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(rendered).To(ContainSubstring("max-complexity: 20"))
			Expect(rendered).To(ContainSubstring("- gosec"))
			Expect(rendered).NotTo(ContainSubstring("- ginkgolinter"))
		})

		It("defaults an omitted root profile to standard", func() {
			tempDir := GinkgoT().TempDir()
			configPath := filepath.Join(tempDir, "goboot.yml")
			content := "projectName: TestProject\n" +
				"targetPath: /tmp/test\n" +
				"services: []\n"

			Expect(os.WriteFile(configPath, []byte(content), 0o644)).To(Succeed())
			goBoot := config.NewGoBoot(configPath)
			Expect(goBoot.Init()).To(Succeed())
			Expect(goBoot.Profile).To(Equal(goboottypes.ProfileStandard))
		})
	})

	Describe("When the setup is incorrect", func() {
		It("rejects unsupported root profiles", func() {
			tempDir := GinkgoT().TempDir()
			configPath := filepath.Join(tempDir, "goboot.yml")
			content := "projectName: TestProject\n" +
				"targetPath: /tmp/test\n" +
				"profile: impossible\n" +
				"services: []\n"

			Expect(os.WriteFile(configPath, []byte(content), 0o644)).To(Succeed())
			goBoot := config.NewGoBoot(configPath)
			err := goBoot.Init()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported profile"))
		})
	})
})
