package basedocker_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/basedocker"
	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

const (
	dockerFileEntry       = "dockerfile"
	composeFileEntry      = "compose"
	dockerignoreFileEntry = "dockerignore"
	expectedBuildCommand  = "docker build -t introproject:local ."
	expectedComposeConfig = "docker compose config"
)

type recordingRegistrar struct {
	lines map[string][]string
	files map[string][]string
}

func newRecordingRegistrar() *recordingRegistrar {
	return &recordingRegistrar{
		lines: make(map[string][]string),
		files: make(map[string][]string),
	}
}

func (r *recordingRegistrar) RegisterLines(name string, lines []string) error {
	r.lines[name] = lines

	return nil
}

func (r *recordingRegistrar) RegisterFile(name string, lines []string) error {
	r.files[name] = lines

	return nil
}

var _ = Describe("BaseDocker Service", func() {
	var (
		tempDir     string
		sourceDir   string
		baseDocker  *basedocker.BaseDocker
		validConfig *config.BaseDockerConfig
	)

	BeforeEach(func() {
		var err error

		tempDir, err = os.MkdirTemp("", "basedocker-test-*")
		Expect(err).NotTo(HaveOccurred())

		sourceDir, err = os.MkdirTemp("", "basedocker-source-*")
		Expect(err).NotTo(HaveOccurred())

		validConfig = &config.BaseDockerConfig{
			SourcePath:       sourceDir,
			ProjectName:      "IntroProject",
			LowerProjectName: "introproject",
			GoVersion:        "1.26.3",
			RuntimeImage:     "gcr.io/distroless/static-debian12:nonroot",
			BinaryName:       "introproject",
			MainPackage:      "./cmd/introproject",
			ImageName:        "introproject:local",
			FileList:         []string{dockerFileEntry, composeFileEntry, dockerignoreFileEntry},
			Ports:            []string{"8080:8080"},
		}

		baseDocker = basedocker.NewBaseDocker(tempDir)
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}

		if sourceDir != "" {
			Expect(os.RemoveAll(sourceDir)).To(Succeed())
		}
	})

	createSourceFile := func(fileName, content string) {
		fullPath := filepath.Join(sourceDir, fileName+goboottypes.TemplateSuffix)
		Expect(os.WriteFile(fullPath, []byte(content), 0o644)).To(Succeed())
	}

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseDocker.ID()).To(Equal(goboottypes.ServiceNameBaseDocker))
		})
	})

	Describe("SetConfig", func() {
		It("accepts BaseDockerConfig", func() {
			err := baseDocker.SetConfig(validConfig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error for wrong config type", func() {
			err := baseDocker.SetConfig(&config.BaseProjectConfig{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid config type"))
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			createSourceFile("Dockerfile", "FROM golang:{{.GoVersion}}\n")
			createSourceFile("docker-compose.yml", "image: {{.ImageName}}\n{{ range .Ports }}port: {{ . }}\n{{ end }}")
			createSourceFile(".dockerignore", "outputs\n")

			Expect(baseDocker.SetConfig(validConfig)).To(Succeed())
		})

		It("copies and renders enabled templates", func() {
			Expect(baseDocker.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			Expect(filepath.Join(targetRoot, "Dockerfile")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, "docker-compose.yml")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, ".dockerignore")).To(BeAnExistingFile())

			content, err := os.ReadFile(filepath.Join(targetRoot, "docker-compose.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("introproject:local"))
			Expect(string(content)).To(ContainSubstring("8080:8080"))
		})

		It("respects fileList selection", func() {
			validConfig.FileList = []string{dockerFileEntry}

			Expect(baseDocker.Run()).To(Succeed())

			targetRoot := filepath.Join(tempDir, validConfig.ProjectName)
			Expect(filepath.Join(targetRoot, "Dockerfile")).To(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, "docker-compose.yml")).NotTo(BeAnExistingFile())
			Expect(filepath.Join(targetRoot, ".dockerignore")).NotTo(BeAnExistingFile())
		})

		It("returns a clear error when an enabled template is missing", func() {
			validConfig.FileList = []string{dockerFileEntry, composeFileEntry}

			Expect(os.Remove(filepath.Join(sourceDir, "docker-compose.yml"+goboottypes.TemplateSuffix))).To(Succeed())

			err := baseDocker.Run()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing required template"))
			Expect(err.Error()).To(ContainSubstring("docker-compose.yml"))
		})

		It("registers local and CI commands", func() {
			scripts := newRecordingRegistrar()
			ciRegistrar := newRecordingRegistrar()

			baseDocker.SetScriptReceiver(scripts)
			baseDocker.SetCIReceiver(ciRegistrar)

			Expect(baseDocker.Run()).To(Succeed())

			expectedScriptLines := []string{
				expectedBuildCommand,
				"docker compose up --build",
				expectedComposeConfig + " && " + expectedBuildCommand + " && docker run --rm introproject:local -h",
			}
			expectedCILines := []string{
				expectedComposeConfig,
				expectedBuildCommand,
				"docker run --rm introproject:local -h",
			}

			Expect(scripts.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedScriptLines))
			Expect(scripts.files).To(HaveKeyWithValue(goboottypes.ScriptFileDocker, expectedScriptLines))
			Expect(ciRegistrar.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedCILines))
			Expect(ciRegistrar.files).To(HaveKeyWithValue(goboottypes.CIFileContainer, expectedCILines))
		})

		It("omits compose commands when compose output is disabled", func() {
			validConfig.FileList = []string{dockerFileEntry}
			scripts := newRecordingRegistrar()
			ciRegistrar := newRecordingRegistrar()

			baseDocker.SetScriptReceiver(scripts)
			baseDocker.SetCIReceiver(ciRegistrar)

			Expect(baseDocker.Run()).To(Succeed())

			expectedScriptLines := []string{
				expectedBuildCommand,
				"echo \"No compose output enabled.\"",
				expectedBuildCommand + " && docker run --rm introproject:local -h",
			}
			expectedCILines := []string{
				expectedBuildCommand,
				"docker run --rm introproject:local -h",
			}

			Expect(scripts.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedScriptLines))
			Expect(scripts.files).To(HaveKeyWithValue(goboottypes.ScriptFileDocker, expectedScriptLines))
			Expect(ciRegistrar.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedCILines))
			Expect(ciRegistrar.files).To(HaveKeyWithValue(goboottypes.CIFileContainer, expectedCILines))
		})

		It("omits Docker build commands when Dockerfile output is disabled", func() {
			validConfig.FileList = []string{composeFileEntry}
			scripts := newRecordingRegistrar()
			ciRegistrar := newRecordingRegistrar()

			baseDocker.SetScriptReceiver(scripts)
			baseDocker.SetCIReceiver(ciRegistrar)

			Expect(baseDocker.Run()).To(Succeed())

			expectedScriptLines := []string{
				"echo \"No Dockerfile output enabled.\"",
				"docker compose up --build",
				expectedComposeConfig,
			}
			expectedCILines := []string{
				expectedComposeConfig,
			}

			Expect(scripts.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedScriptLines))
			Expect(scripts.files).To(HaveKeyWithValue(goboottypes.ScriptFileDocker, expectedScriptLines))
			Expect(ciRegistrar.lines).To(HaveKeyWithValue(goboottypes.ServiceNameBaseDocker, expectedCILines))
			Expect(ciRegistrar.files).To(HaveKeyWithValue(goboottypes.CIFileContainer, expectedCILines))
		})
	})
})
