package gobootutils_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/gobootutils"
)

var _ = Describe("Template helpers (rendering)", func() {
	var (
		tempDir string
		root    *os.Root
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "goboot-template-*")
		Expect(err).NotTo(HaveOccurred())

		root, err = os.OpenRoot(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if root != nil {
			Expect(root.Close()).To(Succeed())
		}

		if tempDir != "" {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		}
	})

	Describe("ExecuteTemplateText", func() {
		Context("with valid template and data", func() {
			It("renders the template correctly", func() {
				template := "Hello {{.Name}}!"
				data := struct{ Name string }{Name: "World"}

				result, err := gobootutils.ExecuteTemplateText("test", template, data)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("Hello World!"))
			})
		})

		Context("with complex template data", func() {
			It("handles nested fields", func() {
				template := "Project: {{.ProjectName}}, Version: {{.Version}}"
				data := struct {
					ProjectName string
					Version     string
				}{
					ProjectName: "goboot",
					Version:     "1.0.0",
				}

				result, err := gobootutils.ExecuteTemplateText("complex", template, data)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("Project: goboot, Version: 1.0.0"))
			})
		})

		Context("with invalid template syntax", func() {
			It("returns an error", func() {
				template := "Hello {{.Name" // Missing closing braces
				data := struct{ Name string }{Name: "World"}

				_, err := gobootutils.ExecuteTemplateText("invalid", template, data)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed template parse"))
			})
		})

		Context("with missing template data", func() {
			It("returns an error when accessing undefined fields", func() {
				template := "Hello {{.MissingField}}!"
				data := struct{ Name string }{Name: "World"}

				_, err := gobootutils.ExecuteTemplateText("missing", template, data)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed template execution"))
			})
		})

		Context("with indent helper", func() {
			It("indents all lines in the provided string", func() {
				template := "entry:\n{{ indent 2 .Command }}"
				data := struct{ Command string }{
					Command: "line1\nline2\n\nline4",
				}

				result, err := gobootutils.ExecuteTemplateText("indent", template, data)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal("entry:\n  line1\n  line2\n\n  line4"))
			})
		})
	})

	Describe("RenderTemplateToFile", func() {
		Context("when rendering a template file", func() {
			It("renders and overwrites the file correctly", func() {
				templateContent := "Name: {{.Name}}\nValue: {{.Value}}"
				file, err := root.Create("template.txt")
				Expect(err).NotTo(HaveOccurred())
				_, err = file.WriteString(templateContent)
				Expect(err).NotTo(HaveOccurred())
				Expect(file.Close()).To(Succeed())

				data := struct {
					Name  string
					Value int
				}{
					Name:  "TestItem",
					Value: 42,
				}

				err = gobootutils.RenderTemplateToFile("test-render", root, "template.txt", data)
				Expect(err).NotTo(HaveOccurred())

				content, err := os.ReadFile(filepath.Join(tempDir, "template.txt"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(Equal("Name: TestItem\nValue: 42"))
			})
		})

		Context("when file does not exist", func() {
			It("returns an error", func() {
				data := struct{ Name string }{Name: "Test"}
				err := gobootutils.RenderTemplateToFile("missing", root, "nonexistent.txt", data)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to open file"))
			})
		})

		Context("when file cannot be opened", func() {
			It("returns an error", func() {
				file, err := root.Create("locked.txt")
				Expect(err).NotTo(HaveOccurred())
				Expect(file.Close()).To(Succeed())
				Expect(os.Chmod(filepath.Join(tempDir, "locked.txt"), 0o000)).To(Succeed())
				defer func() {
					Expect(os.Chmod(filepath.Join(tempDir, "locked.txt"), 0o644)).To(Succeed())
				}()

				err = gobootutils.RenderTemplateToFile("locked", root, "locked.txt", struct{}{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to open file"))
			})
		})

		Context("when template content is invalid", func() {
			It("returns a rendering error", func() {
				file, err := root.Create("broken.txt")
				Expect(err).NotTo(HaveOccurred())
				_, err = file.WriteString("{{")
				Expect(err).NotTo(HaveOccurred())
				Expect(file.Close()).To(Succeed())

				err = gobootutils.RenderTemplateToFile("broken", root, "broken.txt", struct{}{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed template render"))
			})
		})
	})

	Describe("oneLine", func() {
		It("converts newlines to spaces and trims surrounding whitespace", func() {
			input := "  first line\r\nsecond line\nthird line  "
			result, err := gobootutils.ExecuteTemplateText("oneline", "{{ oneLine .Val }}", map[string]string{
				"Val": input,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("first line second line third line"))
		})

		It("preserves intentional spacing between content segments", func() {
			input := "alpha\n\nbeta\tgamma"
			result, err := gobootutils.ExecuteTemplateText("oneline-spaces", "{{ oneLine .Val }}", map[string]string{
				"Val": input,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("alpha  beta\tgamma"))
		})
	})
})
