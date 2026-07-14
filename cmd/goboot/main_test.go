package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

		It("prints the error and exits with non-zero status on failure", func() {
			defer withFakeGo()()

			buf := &bytes.Buffer{}
			outputWriter = buf

			var exitCode int

			exitFunc = func(code int) { exitCode = code }
			os.Args = []string{"goboot", argConfig, "/nonexistent/path.yml"}

			main()

			Expect(exitCode).To(Equal(1))
			Expect(buf.String()).To(ContainSubstring("failed to initialize configuration"))
		})
	})
})
