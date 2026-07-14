package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const schemaFileCount = 11

var _ = Describe("Published configuration schemas", func() {
	var repositoryRoot string

	BeforeEach(func() {
		_, currentFile, _, found := runtime.Caller(0)
		Expect(found).To(BeTrue())

		repositoryRoot = filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	})

	It("publishes valid JSON Schema documents", func() {
		schemaDirectory := filepath.Join(repositoryRoot, "schemas")
		entries, err := os.ReadDir(schemaDirectory)
		Expect(err).NotTo(HaveOccurred())

		jsonFiles := 0

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
				continue
			}

			jsonFiles++

			content, err := os.ReadFile(filepath.Join(schemaDirectory, entry.Name()))
			Expect(err).NotTo(HaveOccurred())

			var schema map[string]any

			Expect(json.Unmarshal(content, &schema)).To(Succeed())
			Expect(schema).To(HaveKeyWithValue("$schema", "https://json-schema.org/draft/2020-12/schema"))
			Expect(schema).To(HaveKey("$id"))
			Expect(schema).To(HaveKeyWithValue("additionalProperties", false))
		}

		Expect(jsonFiles).To(Equal(schemaFileCount))
	})

	It("connects every example config to its matching schema", func() {
		configDirectory := filepath.Join(repositoryRoot, "configs")
		entries, err := os.ReadDir(configDirectory)
		Expect(err).NotTo(HaveOccurred())

		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".yml" {
				continue
			}

			content, err := os.ReadFile(filepath.Join(configDirectory, entry.Name()))
			Expect(err).NotTo(HaveOccurred())

			schemaName := strings.TrimSuffix(entry.Name(), ".yml") + ".schema.json"
			modeline := "# yaml-language-server: $schema=../schemas/" + schemaName

			Expect(string(content)).To(HavePrefix(modeline + "\n"))
			Expect(filepath.Join(repositoryRoot, "schemas", schemaName)).To(BeAnExistingFile())
		}
	})
})
