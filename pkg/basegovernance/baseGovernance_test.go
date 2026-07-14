package basegovernance_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/basegovernance"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseGovernance", func() {
	var (
		targetDir     string
		sourceDir     string
		service       *basegovernance.BaseGovernance
		governanceCfg *config.BaseGovernanceConfig
	)

	BeforeEach(func() {
		var err error

		targetDir, err = os.MkdirTemp("", "basegovernance-target-*")
		Expect(err).NotTo(HaveOccurred())
		sourceDir, err = os.MkdirTemp("", "basegovernance-source-*")
		Expect(err).NotTo(HaveOccurred())

		writeGovernanceTemplates(sourceDir)

		governanceCfg = &config.BaseGovernanceConfig{
			SourcePath:    sourceDir,
			Maintainers:   []string{"@maintainer"},
			DefaultBranch: "main",
			ProjectName:   "IntroProject",
			ProjectURL:    "https://github.com/example/introproject",
			GitProvider:   goboottypes.GitProviderGitHub,
			Profile:       goboottypes.ProfileStandard,
		}
		service = basegovernance.NewBaseGovernance(targetDir)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(targetDir)).To(Succeed())
		Expect(os.RemoveAll(sourceDir)).To(Succeed())
	})

	It("renders common and standard GitHub governance files", func() {
		Expect(service.SetConfig(governanceCfg)).To(Succeed())
		Expect(service.Run()).To(Succeed())

		projectDir := filepath.Join(targetDir, "IntroProject")
		Expect(filepath.Join(projectDir, "CODEOWNERS")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, "CONTRIBUTING.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, "SECURITY.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".github", "PULL_REQUEST_TEMPLATE.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".github", "ISSUE_TEMPLATE", "bug_report.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".github", "ISSUE_TEMPLATE", "feature_request.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".github", "ISSUE_TEMPLATE", "documentation.yml")).NotTo(BeAnExistingFile())

		content, err := os.ReadFile(filepath.Join(projectDir, "CONTRIBUTING.md"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring("IntroProject"))
	})

	It("renders the enterprise GitLab change workflow", func() {
		governanceCfg.GitProvider = goboottypes.GitProviderGitLab
		governanceCfg.ProjectURL = "https://gitlab.com/example/introproject"
		governanceCfg.Profile = goboottypes.ProfileEnterprise

		Expect(service.SetConfig(governanceCfg)).To(Succeed())
		Expect(service.Run()).To(Succeed())

		projectDir := filepath.Join(targetDir, "IntroProject")
		Expect(filepath.Join(projectDir, ".gitlab", "merge_request_templates", "Default.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".gitlab", "issue_templates", "Bug.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".gitlab", "issue_templates", "Feature.md")).To(BeAnExistingFile())
		Expect(filepath.Join(projectDir, ".gitlab", "issue_templates", "Change.md")).To(BeAnExistingFile())
	})

	It("uses the reduced minimal issue baseline", func() {
		governanceCfg.Profile = goboottypes.ProfileMinimal
		Expect(service.SetConfig(governanceCfg)).To(Succeed())
		Expect(service.Run()).To(Succeed())

		issueDir := filepath.Join(targetDir, "IntroProject", ".github", "ISSUE_TEMPLATE")
		Expect(filepath.Join(issueDir, "bug_report.yml")).To(BeAnExistingFile())
		Expect(filepath.Join(issueDir, "feature_request.yml")).NotTo(BeAnExistingFile())
	})

	It("returns its stable identifier", func() {
		Expect(service.ID()).To(Equal(goboottypes.ServiceNameBaseGovernance))
	})

	It("rejects another service config type", func() {
		Expect(service.SetConfig(&config.BaseDockerConfig{})).NotTo(Succeed())
	})
})

func writeGovernanceTemplates(root string) {
	templates := map[string]string{
		"CODEOWNERS.tmpl":                         "* {{range .Maintainers}}{{.}} {{end}}\n",
		"CONTRIBUTING.md.tmpl":                    "# Contributing to {{.ProjectName}}\n",
		"SECURITY.md.tmpl":                        "# Security for {{.ProjectName}}\n",
		"github/PULL_REQUEST_TEMPLATE.md.tmpl":    "# Pull request\n",
		"github/bug_report.yml.tmpl":              "name: Bug for {{.ProjectName}}\n",
		"github/feature_request.yml.tmpl":          "name: Feature\n",
		"github/documentation.yml.tmpl":            "name: Documentation\n",
		"github/config.yml.tmpl":                   "blank_issues_enabled: false\n",
		"github/change_request.yml.tmpl":           "name: Change\n",
		"gitlab/Default.md.tmpl":                   "# Merge request\n",
		"gitlab/Bug.md.tmpl":                       "# Bug\n",
		"gitlab/Feature.md.tmpl":                   "# Feature\n",
		"gitlab/Documentation.md.tmpl":             "# Documentation\n",
		"gitlab/Change.md.tmpl":                    "# Change\n",
	}

	for relativePath, content := range templates {
		fullPath := filepath.Join(root, filepath.FromSlash(relativePath))
		Expect(os.MkdirAll(filepath.Dir(fullPath), 0o755)).To(Succeed())
		Expect(os.WriteFile(fullPath, []byte(content), 0o644)).To(Succeed())
	}
}
