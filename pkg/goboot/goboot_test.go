package goboot_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboot"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

func writeGoMod(dir string) {
	content := "module example.com/temp\n\ngo 1.21\n"
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o644)
}

var _ = Describe("GoBoot Core Orchestration", func() {
	var (
		tempDir string
		goBoot  *goboot.GoBoot
		cfg     *config.GoBoot
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "goboot-test-*")
		Expect(err).NotTo(HaveOccurred())

		// Create a basic config
		cfg = &config.GoBoot{
			TargetPath:  tempDir,
			ConfManager: config.NewConfigManager(),
			Services: []config.ServiceConfigMeta{
				{
					ID:       goboottypes.ServiceNameBaseProject,
					ConfPath: "/path/to/config.yml",
					Enabled:  true,
				},
			},
		}
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("NewGoBoot", func() {
		It("creates a new GoBoot instance", func() {
			goBoot = goboot.NewGoBoot(cfg)
			Expect(goBoot).NotTo(BeNil())
		})

		It("wires the service manager correctly", func() {
			goBoot = goboot.NewGoBoot(cfg)
			Expect(goBoot.ServiceMgr).NotTo(BeNil())
		})
	})

	Describe("RegisterServices", func() {
		BeforeEach(func() {
			goBoot = goboot.NewGoBoot(cfg)
		})

		Context("with valid configuration", func() {
			It("creates target directory if it doesn't exist", func() {
				targetDir := filepath.Join(tempDir, "new-target")
				cfg.TargetPath = targetDir

				Expect(goBoot.RegisterServices()).To(Succeed())
				// May error due to service registration, but directory should be created
				info, statErr := os.Stat(targetDir)
				Expect(statErr).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})

			It("succeeds when target directory already exists", func() {
				// tempDir already exists from BeforeEach
				cfg.TargetPath = tempDir
				Expect(goBoot.RegisterServices()).To(Succeed())
				// Verify directory still exists
				info, err := os.Stat(tempDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})

		Context("with no services declared", func() {
			It("returns an error", func() {
				cfg.Services = nil
				err := goBoot.RegisterServices()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no services declared"))
			})
		})

		Context("with empty services list", func() {
			It("succeeds but registers nothing", func() {
				cfg.Services = []config.ServiceConfigMeta{}
				err := goBoot.RegisterServices()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when target directory creation fails", func() {
			It("returns an error for invalid path", func() {
				// Use an invalid path that will fail
				cfg.TargetPath = "/dev/null/cannot/create/here"
				err := goBoot.RegisterServices()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to create target directory"))
			})
		})

		Context("with unknown service ID", func() {
			It("returns an error for unknown service", func() {
				cfg.Services = []config.ServiceConfigMeta{
					{
						ID:      "unknown_service_xyz",
						Enabled: true,
					},
				}
				err := goBoot.RegisterServices()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown service ID"))
			})
		})

		Context("with disabled services", func() {
			It("skips disabled services", func() {
				cfg.Services = []config.ServiceConfigMeta{
					{
						ID:      goboottypes.ServiceNameBaseProject,
						Enabled: false,
					},
					{
						ID:      goboottypes.ServiceNameBaseLint,
						Enabled: false,
					},
				}
				err := goBoot.RegisterServices()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with mixed enabled and disabled services", func() {
			It("registers only enabled services", func() {
				cfg.Services = []config.ServiceConfigMeta{
					{
						ID:      goboottypes.ServiceNameBaseProject,
						Enabled: true,
					},
					{
						ID:      goboottypes.ServiceNameBaseLint,
						Enabled: false,
					},
				}
				Expect(goBoot.RegisterServices()).To(Succeed())
			})
		})

		Context("with all built-in services including tests", func() {
			It("registers base_test alongside other services", func() {
				cfg.Services = []config.ServiceConfigMeta{
					{
						ID:      goboottypes.ServiceNameBaseLocal,
						Enabled: true,
					},
					{
						ID:      goboottypes.ServiceNameBaseProject,
						Enabled: true,
					},
					{
						ID:      goboottypes.ServiceNameBaseLint,
						Enabled: true,
					},
					{
						ID:      goboottypes.ServiceNameBaseTest,
						Enabled: true,
					},
				}

				goBoot = goboot.NewGoBoot(cfg)
				Expect(goBoot.RegisterServices()).To(Succeed())
				writeGoMod(tempDir)
				Expect(goBoot.RunServices()).To(Succeed())
			})
		})
	})

	Describe("RunServices", func() {
		BeforeEach(func() {
			goBoot = goboot.NewGoBoot(cfg)
		})

		Context("with no services registered", func() {
			It("completes without error", func() {
				cfg.Services = []config.ServiceConfigMeta{}
				err := goBoot.RegisterServices()
				Expect(err).NotTo(HaveOccurred())

				writeGoMod(tempDir)
				err = goBoot.RunServices()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("RunGoModTidy", func() {
			var (
				origPath string
				goBinDir string
			)

			writeFakeGo := func(content string) {
				goScript := filepath.Join(goBinDir, "go")
				Expect(os.WriteFile(goScript, []byte(content), 0o755)).To(Succeed())
			}

			BeforeEach(func() {
				var err error
				goBinDir, err = os.MkdirTemp("", "fake-go-bin-*")
				Expect(err).NotTo(HaveOccurred())
				origPath = os.Getenv("PATH")
				Expect(os.Setenv("PATH", goBinDir+string(os.PathListSeparator)+origPath)).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Setenv("PATH", origPath)).To(Succeed())
				if goBinDir != "" {
					Expect(os.RemoveAll(goBinDir)).To(Succeed())
				}
			})

			It("skips when execute is false", func() {
				Expect(goBoot.RunGoModTidy(false)).To(Succeed())
			})

			It("skips when go.mod does not exist", func() {
				Expect(goBoot.RunGoModTidy(true)).To(Succeed())
			})

			It("runs go mod tidy when go.mod exists", func() {
				targetDir := filepath.Join(tempDir, "tidy-success")
				Expect(os.MkdirAll(targetDir, 0o755)).To(Succeed())
				cfg.ProjectName = "tidyproj"
				projectRoot := filepath.Join(targetDir, cfg.ProjectName)
				Expect(os.MkdirAll(projectRoot, 0o755)).To(Succeed())
				writeGoMod(projectRoot)

				marker := filepath.Join(tempDir, "go-invoked")
				writeFakeGo("#!/usr/bin/env bash\npwd > \"" + marker + "\"\nprintf \"%s\" \"$@\" >> \"" + marker + "\"\n")

				cfg.TargetPath = targetDir
				goBoot = goboot.NewGoBoot(cfg)

				Expect(goBoot.RunGoModTidy(true)).To(Succeed())

				data, err := os.ReadFile(marker)
				Expect(err).NotTo(HaveOccurred())
				output := string(data)
				Expect(strings.HasPrefix(output, projectRoot)).To(BeTrue())
				Expect(output).To(ContainSubstring("modtidy"))
			})

			It("propagates go command failures", func() {
				targetDir := filepath.Join(tempDir, "tidy-fail")
				Expect(os.MkdirAll(targetDir, 0o755)).To(Succeed())
				cfg.ProjectName = "tidyproj"
				projectRoot := filepath.Join(targetDir, cfg.ProjectName)
				Expect(os.MkdirAll(projectRoot, 0o755)).To(Succeed())
				writeGoMod(projectRoot)

				writeFakeGo("#!/usr/bin/env bash\nexit 1\n")
				cfg.TargetPath = targetDir
				goBoot = goboot.NewGoBoot(cfg)

				err := goBoot.RunGoModTidy(true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to run go mod tidy"))
			})
		})
	})

	Describe("Service Registration Flow", func() {
		It("follows the correct order: pre-services, main services", func() {
			cfg.Services = []config.ServiceConfigMeta{
				{
					ID:      goboottypes.ServiceNameBaseLocal,
					Enabled: true,
				},
				{
					ID:      goboottypes.ServiceNameBaseProject,
					Enabled: true,
				},
				{
					ID:      goboottypes.ServiceNameBaseLint,
					Enabled: true,
				},
			}
			goBoot = goboot.NewGoBoot(cfg)
			Expect(goBoot.RegisterServices()).To(Succeed())
		})
	})

	Describe("Real-world scenarios", func() {
		Context("when setting up a new project", func() {
			It("handles typical service configuration", func() {
				cfg.Services = []config.ServiceConfigMeta{
					{
						ID:      goboottypes.ServiceNameBaseProject,
						Enabled: true,
					},
				}
				goBoot = goboot.NewGoBoot(cfg)

				Expect(goBoot.RegisterServices()).To(Succeed())

				// Target directory should exist
				info, err := os.Stat(tempDir)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})
	})

	Describe("Service name validation", func() {
		Context("with valid service names", func() {
			DescribeTable("accepts known service IDs",
				func(serviceID string) {
					cfg.Services = []config.ServiceConfigMeta{
						{
							ID:      serviceID,
							Enabled: true,
						},
					}
					goBoot = goboot.NewGoBoot(cfg)
					Expect(goBoot.RegisterServices()).To(Succeed())
				},
				Entry("base_project", goboottypes.ServiceNameBaseProject),
				Entry("base_lint", goboottypes.ServiceNameBaseLint),
				Entry("base_local", goboottypes.ServiceNameBaseLocal),
			)
		})
	})
})
