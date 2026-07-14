package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/regeneration"
)

func withFakeGo() func() {
	return withFakeGoScript("#!/usr/bin/env bash\nexit 0\n")
}

func withFakeGoScript(script string) func() {
	fakeDir, err := os.MkdirTemp("", "fake-go-*")
	if err != nil {
		panic(err)
	}

	goPath := filepath.Join(fakeDir, "go")

	err = os.WriteFile(goPath, []byte(script), 0o755)
	if err != nil {
		panic(err)
	}

	origPath := os.Getenv("PATH")

	err = os.Setenv("PATH", fakeDir+string(os.PathListSeparator)+origPath)
	if err != nil {
		panic(err)
	}

	return func() {
		_ = os.Setenv("PATH", origPath)
		_ = os.RemoveAll(fakeDir)
	}
}

var _ = Describe("CLI entrypoint", func() {
	// No BeforeEach needed for flag cleanup anymore as strictly local FlagSets are used.
	It("runs end-to-end with a minimal valid config", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		targetDir := filepath.Join(tempDir, "out")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliProject",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		// Pass args explicitly, avoiding os.Args hacks
		err = run([]string{argConfig, configFile})
		Expect(err).To(Succeed())

		info, err := os.Stat(targetDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.IsDir()).To(BeTrue())
	})

	It("returns error for missing config file", func() {
		err := run([]string{argConfig, "/nonexistent/path.yml"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to initialize configuration"))
	})

	It("returns error for malformed flags", func() {
		err := run([]string{"--unknown-flag"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to parse flags"))
	})

	It("returns error for invalid log level", func() {
		err := run([]string{"--log-level", "verbose"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid --log-level value"))
	})

	It("returns error for malformed YAML", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		yamlContent, err := loadTestFixture("cmd_goboot/goboot/invalid.yml")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to initialize configuration"))
	})

	It("returns error when service registration fails", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		// no services declared -> RegisterServices fails
		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/missing_services.yml", map[string]string{
			fixtureTargetDir: filepath.Join(tempDir, "out"),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("service registration failed"))
		Expect(commandExitCode(err)).To(Equal(exitGeneration))
	})

	It("returns error when service execution fails", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		baseProjConfig := filepath.Join(tempDir, "base_project.yml")
		sourceDir := filepath.Join(tempDir, "templates")
		Expect(os.MkdirAll(sourceDir, 0o755)).To(Succeed())
		// invalid template filename to force render path failure
		Expect(os.WriteFile(filepath.Join(sourceDir, "{{.ProjectName"), []byte("content"), 0o644)).To(Succeed())

		baseProjContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/service_execution_base_project.yml",
			map[string]string{
				fixtureSourceDir: sourceDir,
			})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(baseProjConfig, baseProjContent, 0o644)).To(Succeed())

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/service_execution_goboot.yml", map[string]string{
			fixtureTargetDir:   filepath.Join(tempDir, "out"),
			fixtureBaseProject: baseProjConfig,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("service execution failed"))
	})

	It("accepts a valid log level", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		targetDir := filepath.Join(tempDir, "out")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliProject",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile, "--log-level", "debug"})
		Expect(err).To(Succeed())
	})

	It("skips host go mod tidy when requested", func() {
		defer withFakeGoScript("#!/usr/bin/env bash\nexit 42\n")()

		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		targetDir := filepath.Join(tempDir, "out")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliSkipTidy",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile, "--skip-go-mod-tidy"})
		Expect(err).To(Succeed())
	})

	It("prints a dry-run plan without changing the target", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		targetDir := filepath.Join(tempDir, "out")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliDryRun",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		originalWriter := outputWriter
		buffer := &bytes.Buffer{}

		outputWriter = buffer
		defer func() {
			outputWriter = originalWriter
		}()

		err = run([]string{argConfig, configFile, "--dry-run"})
		Expect(err).To(Succeed())
		Expect(targetDir).NotTo(BeAnExistingFile())
		Expect(buffer.String()).To(ContainSubstring("CREATE\t.goboot-manifest.yml"))
		Expect(buffer.String()).To(ContainSubstring("dry run completed"))
	})

	It("rejects an invalid regeneration-policy override", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliPolicy",
			fixtureTargetDir:   filepath.Join(tempDir, "out"),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		err = run([]string{argConfig, configFile, "--regeneration-policy", "merge"})
		Expect(err).To(MatchError(ContainSubstring("invalid regeneration policy")))
	})

	It("validates configuration without generating output", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")
		targetDir := filepath.Join(tempDir, "out")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliValidate",
			fixtureTargetDir:   targetDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		originalWriter := outputWriter
		buffer := &bytes.Buffer{}

		outputWriter = buffer
		defer func() {
			outputWriter = originalWriter
		}()

		err = run([]string{argConfig, configFile, "--validate"})
		Expect(err).NotTo(HaveOccurred())
		Expect(targetDir).NotTo(BeAnExistingFile())
		Expect(buffer.String()).To(Equal("configuration is valid.\n"))
	})

	It("emits one machine-readable validation result", func() {
		tempDir := GinkgoT().TempDir()
		configFile := filepath.Join(tempDir, "goboot.yml")

		yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
			fixtureProjectName: "CliJSON",
			fixtureTargetDir:   filepath.Join(tempDir, "out"),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

		originalWriter := outputWriter
		buffer := &bytes.Buffer{}

		outputWriter = buffer
		defer func() {
			outputWriter = originalWriter
		}()

		err = run([]string{argConfig, configFile, "--validate", "--output", "json"})
		Expect(err).NotTo(HaveOccurred())

		var result successResult
		Expect(json.Unmarshal(buffer.Bytes(), &result)).To(Succeed())
		Expect(result.Status).To(Equal("success"))
		Expect(result.Operation).To(Equal(operationValidate))
		Expect(result.Project).To(Equal("CliJSON"))
	})

	It("prints version without loading configuration", func() {
		originalWriter := outputWriter
		buffer := &bytes.Buffer{}

		outputWriter = buffer
		defer func() {
			outputWriter = originalWriter
		}()

		err := run([]string{"--version"})
		Expect(err).NotTo(HaveOccurred())
		Expect(buffer.String()).To(Equal("goboot " + version + "\n"))
	})

	It("prints help successfully without loading configuration", func() {
		originalWriter := outputWriter
		buffer := &bytes.Buffer{}

		outputWriter = buffer
		defer func() {
			outputWriter = originalWriter
		}()

		err := run([]string{"--help"})
		Expect(err).NotTo(HaveOccurred())
		Expect(buffer.String()).To(ContainSubstring("Usage of goboot:"))
		Expect(buffer.String()).To(ContainSubstring("-validate"))
	})

	DescribeTable("assigns stable exit codes",
		func(args []string, expectedCode int) {
			err := run(args)
			Expect(commandExitCode(err)).To(Equal(expectedCode))
		},
		Entry("usage errors", []string{"--output", "xml"}, exitUsage),
		Entry("configuration errors", []string{argConfig, "/nonexistent/path.yml"}, exitConfig),
	)

	It("assigns the dedicated conflict exit code and preserves the target", func() {
		tempDir := GinkgoT().TempDir()
		stagingRoot := filepath.Join(tempDir, "staging")
		targetRoot := filepath.Join(tempDir, "target")
		projectName := "ConflictProject"
		stagedProject := filepath.Join(stagingRoot, projectName)
		targetProject := filepath.Join(targetRoot, projectName)

		Expect(os.MkdirAll(stagedProject, 0o755)).To(Succeed())
		Expect(os.MkdirAll(targetProject, 0o755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(stagedProject, "README.md"), []byte("generated\n"), 0o644)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(targetProject, "README.md"), []byte("user\n"), 0o644)).To(Succeed())

		_, err := applyStagedProject(
			cliOptions{},
			&config.GoBoot{ProjectName: projectName},
			regeneration.PolicyManaged,
			stagingRoot,
			targetRoot,
		)
		Expect(err).To(HaveOccurred())
		Expect(commandExitCode(err)).To(Equal(exitConflict))

		content, readErr := os.ReadFile(filepath.Join(targetProject, "README.md"))
		Expect(readErr).NotTo(HaveOccurred())
		Expect(string(content)).To(Equal("user\n"))
	})

	Describe("main", func() {
		var (
			originalArgs   []string
			originalExit   func(int)
			originalWriter io.Writer
		)

		BeforeEach(func() {
			originalArgs = os.Args
			originalExit = exitFunc
			originalWriter = outputWriter
		})

		AfterEach(func() {
			os.Args = originalArgs
			exitFunc = originalExit
			outputWriter = originalWriter
		})

		It("runs without triggering exit on success", func() {
			defer withFakeGo()()

			tempDir := GinkgoT().TempDir()
			configFile := filepath.Join(tempDir, "goboot.yml")
			targetDir := filepath.Join(tempDir, "out")

			yamlContent, err := loadTestFixtureWithVars("cmd_goboot/goboot/minimal.yml", map[string]string{
				fixtureProjectName: "CliMain",
				fixtureTargetDir:   targetDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(os.WriteFile(configFile, yamlContent, 0o644)).To(Succeed())

			buf := &bytes.Buffer{}
			outputWriter = buf
			exitCalled := false
			exitFunc = func(code int) { exitCalled = true }
			os.Args = []string{"goboot", argConfig, configFile}

			main()

			Expect(exitCalled).To(BeFalse())
			Expect(buf.String()).To(ContainSubstring("goboot execution completed successfully."))

			info, err := os.Stat(targetDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
		})

		It("prints a JSON error and exits with the stable config status", func() {
			defer withFakeGo()()

			buf := &bytes.Buffer{}
			outputWriter = buf

			var exitCode int

			exitFunc = func(code int) { exitCode = code }
			os.Args = []string{"goboot", argConfig, "/nonexistent/path.yml", "--output", "json"}

			main()

			Expect(exitCode).To(Equal(exitConfig))

			var result errorResult
			Expect(json.Unmarshal(buf.Bytes(), &result)).To(Succeed())
			Expect(result.Category).To(Equal(categoryConfig))
			Expect(result.ExitCode).To(Equal(exitConfig))
		})
	})
})
