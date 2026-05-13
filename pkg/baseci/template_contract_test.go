package baseci_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func baseCIRepoRoot(t GinkgoTInterface) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to resolve current file location")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

var _ = Describe("CI template render-field contract", func() {
	It("uses only allowed render fields in ci_base templates", func() {
		allowedFields := map[string]struct{}{
			"ProjectName":      {},
			"CIDir":            {},
			"JobScripts":       {},
			"FileScripts":      {},
			"EnabledJobFiles":  {},
			"GoVersions":       {},
			"AutoBranches":     {},
			"ImagePolicy":      {},
			"GoModulePath":     {},
			"FileAllowFailure": {},
		}

		root := baseCIRepoRoot(GinkgoT())
		templateRoot := filepath.Join(root, "templates", "ci_base")
		fieldPattern := regexp.MustCompile(`\.[A-Z][A-Za-z0-9_]*`)

		var violations []string

		err := filepath.WalkDir(templateRoot, func(path string, dirEntry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if dirEntry.IsDir() || filepath.Ext(path) != ".tmpl" {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template %q: %w", path, err)
			}

			matches := fieldPattern.FindAllString(string(content), -1)
			for _, match := range matches {
				field := strings.TrimPrefix(match, ".")
				if _, ok := allowedFields[field]; ok {
					continue
				}

				relPath, relErr := filepath.Rel(templateRoot, path)
				if relErr != nil {
					relPath = path
				}

				violations = append(violations, relPath+": "+field)
			}

			return nil
		})

		Expect(err).NotTo(HaveOccurred())

		sort.Strings(violations)
		Expect(violations).To(BeEmpty(), "unexpected template render fields:\n%s", strings.Join(violations, "\n"))
	})
})
