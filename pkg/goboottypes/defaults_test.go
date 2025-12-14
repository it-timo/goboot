package goboottypes_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("Types Package - Defaults and Constants", func() {
	Describe("Default Linter Commands", func() {
		DescribeTable("linter commands are non-empty and valid",
			func(cmd string) {
				Expect(cmd).NotTo(BeEmpty())
				Expect(len(cmd)).To(BeNumerically(">", 20), "Command should be a meaningful string (docker run ...)")
				Expect(cmd).To(ContainSubstring("{{DOCKER_RUN}}"))
				Expect(cmd).To(ContainSubstring(":"))
			},

			Entry("Go linter command", goboottypes.DefaultGoLintCmd),
			Entry("YAML linter command", goboottypes.DefaultYMLLintCmd),
			Entry("Make linter command", goboottypes.DefaultMakeLintCmd),
			Entry("Markdown linter command", goboottypes.DefaultMDLintCmd),
			Entry("Shell linter command", goboottypes.DefaultShellLintCmd),
			Entry("SHFMT linter command", goboottypes.DefaultSHFMTCmd),
		)

		Context("when inspecting specific default commands", func() {
			It("matches full golangci-lint command", func() {
				Expect(goboottypes.DefaultGoLintCmd).To(ContainSubstring("golangci/golangci-lint:v2.7.1"))
				Expect(goboottypes.DefaultGoLintCmd).To(ContainSubstring("golangci-lint run ./..."))
			})

			It("matches full yamllint command", func() {
				Expect(goboottypes.DefaultYMLLintCmd).To(ContainSubstring("pipelinecomponents/yamllint:0.35.9"))
				Expect(goboottypes.DefaultYMLLintCmd).To(ContainSubstring("yamllint ."))
			})

			It("matches full checkmake command", func() {
				Expect(goboottypes.DefaultMakeLintCmd).To(ContainSubstring("cytopia/checkmake:latest-0.5"))
				Expect(goboottypes.DefaultMakeLintCmd).To(ContainSubstring("Makefile"))
			})

			It("pins markdownlint docker image, mount, and pattern", func() {
				Expect(goboottypes.DefaultMDLintCmd).To(ContainSubstring("ghcr.io/igorshubovych/markdownlint-cli:v0.46.0"))
				Expect(goboottypes.DefaultMDLintCmd).To(ContainSubstring("markdownlint \"**/*.md\""))
			})

			It("pins shellcheck docker image, mount, and pattern", func() {
				Expect(goboottypes.DefaultShellLintCmd).To(ContainSubstring("cytopia/shellcheck:latest-0.8.0"))
				Expect(goboottypes.DefaultShellLintCmd).To(ContainSubstring("shellcheck {{SH_FILES}}"))
			})

			It("pins shfmt docker image, mount, and pattern", func() {
				Expect(goboottypes.DefaultSHFMTCmd).To(ContainSubstring("cytopia/shfmt:latest-1.10"))
				Expect(goboottypes.DefaultSHFMTCmd).To(ContainSubstring("shfmt -d {{SH_FILES}}"))
			})
		})
	})

	Describe("Default Linter Identifiers", func() {
		It("matches exact identifiers", func() {
			Expect(goboottypes.LinterGo).To(Equal("golang"))
			Expect(goboottypes.LinterYAML).To(Equal("yaml"))
			Expect(goboottypes.LinterMake).To(Equal("make"))
			Expect(goboottypes.LinterMD).To(Equal("markdown"))
			Expect(goboottypes.LinterShell).To(Equal("shellcheck"))
			Expect(goboottypes.LinterSHFMT).To(Equal("shfmt"))
		})
	})

	Describe("Default Test Commands", func() {
		DescribeTable("test commands are non-empty and valid",
			func(cmd string) {
				Expect(cmd).NotTo(BeEmpty())
				Expect(len(cmd)).To(BeNumerically(">", 20), "Command should be a meaningful string (test ...)")
				Expect(cmd).To(ContainSubstring("race"))
			},

			Entry("Go test command", goboottypes.DefaultGoTestCMD),
		)

		Context("when inspecting specific default commands", func() {
			It("matches full golangci-lint command", func() {
				Expect(goboottypes.DefaultGoTestCMD).To(ContainSubstring("go test -race -timeout=5m"))
			})
		})
	})

	Describe("Default Test Identifiers", func() {
		It("matches exact identifiers", func() {
			Expect(goboottypes.TestStyleGinkgo).To(Equal("ginkgo"))
			Expect(goboottypes.TestStyleGo).To(Equal("go"))
		})
	})

	Describe("Default Local Script Names", func() {
		It("matches exact script name constants", func() {
			Expect(goboottypes.ScriptNameMake).To(Equal("make"))
			Expect(goboottypes.ScriptNameTask).To(Equal("task"))
			Expect(goboottypes.ScriptNameScript).To(Equal("script"))
			Expect(goboottypes.ScriptNameCommit).To(Equal("commit"))
		})
	})

	Describe("Default Local File Names", func() {
		It("matches exact directory names", func() {
			Expect(goboottypes.ScriptDirNameScript).To(Equal("scripts"))
		})
	})

	Describe("Default Local Script File Names", func() {
		It("matches exact script file names", func() {
			Expect(goboottypes.ScriptFileLint).To(Equal("lint.sh"))
			Expect(goboottypes.ScriptFileLint).To(HaveSuffix(".sh"))
		})
	})

	Describe("Service Names", func() {
		It("matches exact service names", func() {
			Expect(goboottypes.ServiceNameBaseProject).To(Equal("base_project"))
			Expect(goboottypes.ServiceNameBaseLint).To(Equal("base_lint"))
			Expect(goboottypes.ServiceNameBaseLocal).To(Equal("base_local"))
			Expect(goboottypes.ServiceNameBaseTest).To(Equal("base_test"))
		})
	})

	Describe("Directory Permissions", func() {
		Context("when checking default permissions", func() {
			It("defines directory permissions as 0755", func() {
				Expect(os.FileMode(goboottypes.DirPerm)).To(Equal(os.FileMode(0o755)))
			})
		})

		Context("when validating permission values", func() {
			It("has sensible directory permissions (owner rwx, group rx, others rx)", func() {
				// 0755 = owner: rwx (7), group: r-x (5), others: r-x (5)
				perm := os.FileMode(goboottypes.DirPerm)
				Expect(perm&0o700).To(Equal(os.FileMode(0o700)), "Owner should have full permissions")
				Expect(perm&0o050).To(Equal(os.FileMode(0o050)), "Group should have read+execute")
				Expect(perm&0o005).To(Equal(os.FileMode(0o005)), "Others should have read+execute")
			})
		})
	})

	Describe("File Permissions", func() {
		Context("when checking default permissions", func() {
			It("defines file permissions as 0644", func() {
				Expect(os.FileMode(goboottypes.FilePerm)).To(Equal(os.FileMode(0o644)))
			})
		})

		Context("when validating permission values", func() {
			It("has sensible file permissions (owner rw, group rx, others rx)", func() {
				// 0644 = owner: rw (6), group: r-x (4), others: r-x (4)
				perm := os.FileMode(goboottypes.FilePerm)
				Expect(perm&0o600).To(Equal(os.FileMode(0o600)), "Owner should have read+write permissions")
				Expect(perm&0o040).To(Equal(os.FileMode(0o040)), "Group should have read")
				Expect(perm&0o004).To(Equal(os.FileMode(0o004)), "Others should have read")
			})
		})
	})

	Describe("Script Permissions", func() {
		Context("when checking default permissions", func() {
			It("defines script permissions as 0755", func() {
				Expect(os.FileMode(goboottypes.ScriptPerm)).To(Equal(os.FileMode(0o755)))
			})
		})

		Context("when validating permission values", func() {
			It("has sensible script permissions (owner rwx, group rx, others rx)", func() {
				// 0755 = owner: rwx (7), group: r-x (5), others: r-x (5)
				perm := os.FileMode(goboottypes.ScriptPerm)
				Expect(perm&0o700).To(Equal(os.FileMode(0o700)), "Owner should have full permissions")
				Expect(perm&0o050).To(Equal(os.FileMode(0o050)), "Group should have read+execute")
				Expect(perm&0o005).To(Equal(os.FileMode(0o005)), "Others should have read+execute")
			})
		})
	})
})
