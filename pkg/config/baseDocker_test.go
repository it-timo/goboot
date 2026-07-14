package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/config"
	"github.com/it-timo/goboot/pkg/goboottypes"
)

const (
	dockerFileEntry       = "dockerfile"
	composeFileEntry      = "compose"
	dockerignoreFileEntry = "dockerignore"
)

var _ = Describe("BaseDockerConfig", func() {
	var (
		baseDocker *config.BaseDockerConfig
		tempDir    string
	)

	BeforeEach(func() {
		baseDocker = &config.BaseDockerConfig{
			SourcePath:  "./templates/docker_base",
			ProjectName: testProjectName,
			FileList:    []string{dockerFileEntry, composeFileEntry, dockerignoreFileEntry},
		}

		var err error

		tempDir, err = os.MkdirTemp("", "basedocker-config-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ID", func() {
		It("returns the correct service identifier", func() {
			Expect(baseDocker.ID()).To(Equal(goboottypes.ServiceNameBaseDocker))
			Expect(baseDocker.ID()).To(Equal("base_docker"))
		})
	})

	Describe("Validate", func() {
		It("validates successfully and fills defaults", func() {
			err := baseDocker.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(baseDocker.GoVersion).To(Equal("1.26.5"))
			Expect(baseDocker.BinaryName).To(Equal("testproject"))
			Expect(baseDocker.MainPackage).To(Equal("./cmd/testproject"))
			Expect(baseDocker.ImageName).To(Equal("testproject:local"))
		})

		It("accepts explicit values", func() {
			baseDocker.GoVersion = "1.26.5"
			baseDocker.RuntimeImage = "alpine:3.22"
			baseDocker.BinaryName = "app"
			baseDocker.MainPackage = "./cmd/app"
			baseDocker.Ports = []string{"8080:8080"}

			err := baseDocker.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(baseDocker.RuntimeImage).To(Equal("alpine:3.22"))
		})

		It("errors when required fields are missing", func() {
			baseDocker.SourcePath = ""
			baseDocker.FileList = nil

			err := baseDocker.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("sourcePath"))
			Expect(err.Error()).To(ContainSubstring("fileList"))
		})

		It("errors on unsupported files", func() {
			baseDocker.FileList = []string{dockerFileEntry, composeFileEntry, "unknown"}

			err := baseDocker.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported entry"))
		})

		It("errors on duplicate file entries", func() {
			baseDocker.FileList = []string{dockerFileEntry, dockerFileEntry}

			err := baseDocker.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicates"))
		})

		It("errors on blank file entries", func() {
			baseDocker.FileList = []string{dockerFileEntry, " "}

			err := baseDocker.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("blank entries"))
		})

		It("errors on invalid port mappings", func() {
			baseDocker.Ports = []string{"8080:bad"}

			err := baseDocker.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("port"))
		})

		DescribeTable("errors on unsafe Docker template values",
			func(mutator func()) {
				mutator()

				err := baseDocker.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Or(
					ContainSubstring("goVersion"),
					ContainSubstring("runtimeImage"),
					ContainSubstring("binaryName"),
					ContainSubstring("mainPackage"),
				))
			},
			Entry("goVersion with command separator", func() {
				baseDocker.GoVersion = "1.26; echo bad"
			}),
			Entry("runtimeImage with newline", func() {
				baseDocker.RuntimeImage = "alpine:3.22\nRUN echo bad"
			}),
			Entry("runtimeImage with trailing Dockerfile instruction", func() {
				baseDocker.RuntimeImage = "alpine RUN echo bad"
			}),
			Entry("runtimeImage with command substitution", func() {
				baseDocker.RuntimeImage = "alpine:$(echo bad)"
			}),
			Entry("binaryName with path separator", func() {
				baseDocker.BinaryName = "bin/app"
			}),
			Entry("binaryName with quotes", func() {
				baseDocker.BinaryName = "app\"bad"
			}),
			Entry("mainPackage with semicolon", func() {
				baseDocker.MainPackage = "./cmd/app; echo bad"
			}),
			Entry("mainPackage with ampersands", func() {
				baseDocker.MainPackage = "./cmd/app && echo bad"
			}),
			Entry("mainPackage with backticks", func() {
				baseDocker.MainPackage = "./cmd/`echo bad`"
			}),
			Entry("mainPackage with command substitution", func() {
				baseDocker.MainPackage = "./cmd/$(echo bad)"
			}),
			Entry("mainPackage with parent traversal", func() {
				baseDocker.MainPackage = "./cmd/../bad"
			}),
		)
	})

	Describe("ReadConfig", func() {
		var configPath string

		BeforeEach(func() {
			configPath = filepath.Join(tempDir, "base_docker.yml")
		})

		It("loads the configuration successfully", func() {
			yamlContent, err := loadTestFixture("config/base_docker/valid.yml")
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(configPath, yamlContent, 0o644)
			Expect(err).NotTo(HaveOccurred())

			newConfig := &config.BaseDockerConfig{}
			err = newConfig.ReadConfig(configPath, "", "")
			Expect(err).NotTo(HaveOccurred())

			Expect(newConfig.SourcePath).To(Equal("./templates/docker_base"))
			Expect(newConfig.FileList).To(ContainElement("dockerfile"))
		})

		It("returns an error for invalid YAML", func() {
			invalidYAML, err := loadTestFixture("config/base_docker/invalid.yml")
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(configPath, invalidYAML, 0o644)
			Expect(err).NotTo(HaveOccurred())

			newConfig := &config.BaseDockerConfig{}
			err = newConfig.ReadConfig(configPath, "", "")
			Expect(err).To(HaveOccurred())
		})
	})
})
