package baselocal_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/baselocal"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

var _ = Describe("BaseLocal Service", func() {
	var (
		tempDir     string
		baseLocal   *baselocal.BaseLocal
		validConfig *config.BaseLocalConfig
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "baselocal-test-*")
		Expect(err).NotTo(HaveOccurred())

		validConfig = &config.BaseLocalConfig{
			SourcePath:  tempDir,
			ProjectName: "testproject",
			FileList: []string{
				goboottypes.ScriptNameMake,
				goboottypes.ScriptNameTask,
			},
		}

		baseLocal = baselocal.NewBaseLocal(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("NewBaseLocal", func() {
		It("creates a new BaseLocal instance", func() {
			bl := baselocal.NewBaseLocal("/some/path")
			Expect(bl).NotTo(BeNil())
		})
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseLocal.ID()).To(Equal(goboottypes.ServiceNameBaseLocal))
			Expect(baseLocal.ID()).To(Equal("base_local"))
		})
	})

	Describe("SetConfig", func() {
		var sourceDir string

		BeforeEach(func() {
			var err error
			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		Context("with valid config", func() {
			It("accepts BaseLocalConfig", func() {
				validConfig.SourcePath = sourceDir
				err := baseLocal.SetConfig(validConfig)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid config type", func() {
			It("returns an error for wrong config type", func() {
				wrongConfig := &config.BaseProjectConfig{}
				err := baseLocal.SetConfig(wrongConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})

			It("errors when source and target paths are identical", func() {
				validConfig.SourcePath = tempDir
				err := baseLocal.SetConfig(validConfig)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Registrar Interface", func() {
		var sourceDir string

		BeforeEach(func() {
			var err error
			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())
			validConfig.FileList = []string{
				goboottypes.ScriptNameMake,
				goboottypes.ScriptNameTask,
				goboottypes.ScriptNameScript,
			}
			validConfig.SourcePath = sourceDir
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		Describe("RegisterLines", func() {
			It("accepts script lines for registration", func() {
				lines := []string{"make lint", "make test"}
				err := baseLocal.RegisterLines("test_service", lines)
				Expect(err).NotTo(HaveOccurred())
				Expect(baseLocal.MakeScripts).To(HaveKey("test_service"))
				Expect(baseLocal.TaskScripts).To(HaveKey("test_service"))
			})

			It("handles empty lines", func() {
				err := baseLocal.RegisterLines("test_service", []string{})
				Expect(err).NotTo(HaveOccurred())
			})

			It("prevents duplicate registrations for the same service", func() {
				Expect(baseLocal.RegisterLines("dup_service", []string{"cmd"})).To(Succeed())
				err := baseLocal.RegisterLines("dup_service", []string{"cmd"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already registered"))
			})

			It("skips registrations for disabled script types", func() {
				validConfig.FileList = []string{goboottypes.ScriptNameMake}
				Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

				Expect(baseLocal.RegisterLines("svc", []string{"cmd"})).To(Succeed())
				Expect(baseLocal.TaskScripts).To(BeEmpty())
				Expect(baseLocal.MakeScripts).To(HaveKey("svc"))
			})
		})

		Describe("RegisterFile", func() {
			It("accepts file registration", func() {
				lines := []string{"#!/bin/bash", "echo test"}
				err := baseLocal.RegisterFile("lint.sh", lines)
				Expect(err).NotTo(HaveOccurred())
				Expect(baseLocal.ScriptFiles).To(HaveKey("lint.sh"))
			})

			It("prevents duplicate file registrations", func() {
				validConfig.FileList = []string{goboottypes.ScriptNameScript}
				Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

				Expect(baseLocal.RegisterFile("lint.sh", []string{"echo"})).To(Succeed())
				err := baseLocal.RegisterFile("lint.sh", []string{"echo again"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("already registered"))
			})

			It("ignores registrations when scripts are disabled", func() {
				validConfig.FileList = []string{goboottypes.ScriptNameMake}
				Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

				Expect(baseLocal.RegisterFile("lint.sh", []string{"echo"})).To(Succeed())
				Expect(baseLocal.ScriptFiles).To(BeEmpty())
			})
		})
	})

	Describe("Run", func() {
		var sourceDir string

		createSourceFile := func(relPath, content string) {
			fullPath := filepath.Join(sourceDir, relPath+goboottypes.TemplateSuffix)
			Expect(os.MkdirAll(filepath.Dir(fullPath), 0o755)).To(Succeed())
			Expect(os.WriteFile(fullPath, []byte(content), 0o644)).To(Succeed())
		}

		BeforeEach(func() {
			var err error
			sourceDir, err = os.MkdirTemp("", "source-*")
			Expect(err).NotTo(HaveOccurred())
			validConfig.SourcePath = sourceDir
			validConfig.FileList = []string{
				goboottypes.ScriptNameMake,
				goboottypes.ScriptNameTask,
				goboottypes.ScriptNameCommit,
				goboottypes.ScriptNameScript,
			}
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())
		})

		AfterEach(func() {
			if sourceDir != "" {
				Expect(os.RemoveAll(sourceDir)).To(Succeed())
			}
		})

		It("copies and renders enabled templates", func() {
			createSourceFile("Makefile", "make-{{.ProjectName}}-{{len (index .MakeScripts \"svc\")}}")
			createSourceFile("Taskfile.yml", "task-{{.ProjectName}}")
			createSourceFile(".pre-commit-config.yaml", "commit-{{.ProjectName}}")
			createSourceFile(filepath.Join(goboottypes.ScriptDirNameScript, "lint.sh"),
				"script-{{len (index .ScriptFiles \"lint.sh\")}}")

			Expect(baseLocal.RegisterLines("svc", []string{"cmd1", "cmd2"})).To(Succeed())
			Expect(baseLocal.RegisterFile("lint.sh", []string{"echo lint"})).To(Succeed())

			Expect(baseLocal.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			Expect(filepath.Join(targetRoot, "Makefile")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, "Taskfile.yml")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, ".pre-commit-config.yaml")).To(BeAnExistingFile())
			scriptFile := filepath.Join(targetRoot, goboottypes.ScriptDirNameScript, "lint.sh")
			Expect(scriptFile).To(BeAnExistingFile())

			content, err := os.ReadFile(filepath.Join(targetRoot, "Makefile"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring(validConfig.ProjectName))
			Expect(string(content)).To(ContainSubstring("2")) // len of registered cmds

			scriptContent, err := os.ReadFile(scriptFile)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(scriptContent)).To(ContainSubstring("1")) // one registered file entry
		})

		It("skips script directory when no scripts are registered", func() {
			createSourceFile("Makefile", "make")
			validConfig.FileList = []string{goboottypes.ScriptNameScript}
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

			Expect(baseLocal.Run()).To(Succeed())
			Expect(filepath.Join(tempDir, validConfig.ProjectName, goboottypes.ScriptDirNameScript)).NotTo(BeADirectory())
		})

		It("returns error when a template file is missing", func() {
			validConfig.FileList = []string{goboottypes.ScriptNameMake}
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

			err := baseLocal.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required template \"Makefile\""))
		})

		It("returns error when template rendering fails", func() {
			createSourceFile("Makefile", "{{") // invalid template

			validConfig.FileList = []string{goboottypes.ScriptNameMake}
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

			err := baseLocal.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed template render"))
		})

		It("propagates errors from scripts copy", func() {
			createSourceFile(filepath.Join(goboottypes.ScriptDirNameScript, "lint.sh"), "script")
			validConfig.FileList = []string{goboottypes.ScriptNameScript}
			Expect(baseLocal.SetConfig(validConfig)).To(Succeed())

			// simulate missing registration; copyFiles should skip when no ScriptFiles,
			// so force registration then remove template to fail.
			Expect(baseLocal.RegisterFile("lint.sh", []string{"echo"})).To(Succeed())
			Expect(os.Remove(filepath.Join(sourceDir, goboottypes.ScriptDirNameScript,
				"lint.sh"+goboottypes.TemplateSuffix))).To(Succeed())

			err := baseLocal.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required template \"lint.sh\""))
		})
	})
})
