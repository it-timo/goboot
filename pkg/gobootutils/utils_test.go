package gobootutils_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/it-timo/goboot/pkg/goboottypes"
	"github.com/it-timo/goboot/pkg/gobootutils"
)

var _ = Describe("Utils Package", func() {
	var (
		tempDir string
		root    *os.Root
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "goboot-test-*")
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

	Describe("EnsureDir", func() {
		Context("when creating a single-level directory", func() {
			It("creates the directory successfully", func() {
				err := gobootutils.EnsureDir("testdir", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())

				// Verify directory exists
				info, err := root.Stat("testdir")
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})

		Context("when creating nested directories", func() {
			It("creates all nested directories successfully", func() {
				err := gobootutils.EnsureDir("a/b/c/d", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())

				// Verify nested structure exists
				info, err := root.Stat("a/b/c/d")
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})

		Context("when directory already exists", func() {
			It("is idempotent and does not error", func() {
				// Create directory first time
				err := gobootutils.EnsureDir("existing", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())

				// Create same directory again
				err = gobootutils.EnsureDir("existing", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())

				// Verify it still exists
				info, err := root.Stat("existing")
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})

		Context("when path is empty or current directory", func() {
			It("handles empty string without error", func() {
				err := gobootutils.EnsureDir("", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())
			})

			It("handles current directory marker without error", func() {
				err := gobootutils.EnsureDir(".", root, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when path attempts to escape root (security tests)", func() {
			DescribeTable("rejects path escape attempts",
				func(maliciousPath string) {
					err := gobootutils.EnsureDir(maliciousPath, root, goboottypes.DirPerm)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("path escapes root"))
				},
				Entry("parent directory reference", "../escape"),
				Entry("absolute parent at start", ".."),
				Entry("parent in middle of path", "a/../../../escape"),
				Entry("hidden parent reference", "a/b/../../../escape"),
			)
		})
	})

	Describe("ComparePaths", func() {
		var (
			testFile1 string
			testFile2 string
			sameFile  string
		)

		BeforeEach(func() {
			testFile1 = filepath.Join(tempDir, "file1.txt")
			testFile2 = filepath.Join(tempDir, "file2.txt")
			sameFile = filepath.Join(tempDir, "same.txt")

			// Create test files
			Expect(os.WriteFile(testFile1, []byte("test1"), 0o644)).To(Succeed())
			Expect(os.WriteFile(testFile2, []byte("test2"), 0o644)).To(Succeed())
			Expect(os.WriteFile(sameFile, []byte("same"), 0o644)).To(Succeed())
		})

		Context("when forceDiffer is true", func() {
			It("succeeds when paths are different", func() {
				err := gobootutils.ComparePaths(testFile1, testFile2, true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("errors when paths are the same", func() {
				err := gobootutils.ComparePaths(sameFile, sameFile, true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("must be different"))
			})

			It("errors when paths resolve to same location", func() {
				// Create another reference to the same file
				samePath := filepath.Join(tempDir, ".", "same.txt")
				err := gobootutils.ComparePaths(sameFile, samePath, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when forceDiffer is false", func() {
			It("succeeds when paths are the same", func() {
				err := gobootutils.ComparePaths(sameFile, sameFile, false)
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds when paths resolve to same location", func() {
				samePath := filepath.Join(tempDir, ".", "same.txt")
				err := gobootutils.ComparePaths(sameFile, samePath, false)
				Expect(err).NotTo(HaveOccurred())
			})

			It("errors when paths are different", func() {
				err := gobootutils.ComparePaths(testFile1, testFile2, false)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("must be the same"))
			})
		})

		Context("with path normalization", func() {
			It("handles relative paths correctly", func() {
				// This will compare after converting to absolute paths
				cwd, err := os.Getwd()
				Expect(err).NotTo(HaveOccurred())

				err = gobootutils.ComparePaths(cwd, ".", false)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("CreateRootDir", func() {
		var targetDir string

		BeforeEach(func() {
			var err error
			targetDir, err = os.MkdirTemp("", "goboot-rootdir-test-*")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if targetDir != "" {
				Expect(os.RemoveAll(targetDir)).To(Succeed())
			}
		})

		Context("when creating a new root directory", func() {
			It("creates the directory and returns a Root", func() {
				projectName := "testproject"
				root, err := gobootutils.CreateRootDir(targetDir, projectName)
				Expect(err).NotTo(HaveOccurred())
				Expect(root).NotTo(BeNil())

				// Verify directory was created
				expectedPath := filepath.Join(targetDir, projectName)
				info, err := os.Stat(expectedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})

			It("sets correct permissions", func() {
				projectName := "permtest"
				root, err := gobootutils.CreateRootDir(targetDir, projectName)
				Expect(err).NotTo(HaveOccurred())
				Expect(root).NotTo(BeNil())

				// Verify permissions
				expectedPath := filepath.Join(targetDir, projectName)
				info, err := os.Stat(expectedPath)
				Expect(err).NotTo(HaveOccurred())
				// Check that permissions include at least user rwx
				Expect(info.Mode().Perm() & 0o700).To(Equal(os.FileMode(0o700)))
			})
		})

		Context("when directory already exists", func() {
			It("opens the existing directory without error", func() {
				projectName := "existing"
				expectedPath := filepath.Join(targetDir, projectName)

				// Pre-create the directory
				err := os.MkdirAll(expectedPath, goboottypes.DirPerm)
				Expect(err).NotTo(HaveOccurred())

				// CreateRootDir should still succeed
				root, err := gobootutils.CreateRootDir(targetDir, projectName)
				Expect(err).NotTo(HaveOccurred())
				Expect(root).NotTo(BeNil())
			})
		})

		Context("with nested project names", func() {
			It("handles project names with path separators", func() {
				projectName := "org/team/project"
				root, err := gobootutils.CreateRootDir(targetDir, projectName)
				Expect(err).NotTo(HaveOccurred())
				Expect(root).NotTo(BeNil())

				// Verify nested structure was created
				expectedPath := filepath.Join(targetDir, projectName)
				info, err := os.Stat(expectedPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info.IsDir()).To(BeTrue())
			})
		})

		Context("when target directory is not writable", func() {
			It("returns an error", func() {
				Expect(os.Chmod(targetDir, 0o500)).To(Succeed())
				defer func() {
					Expect(os.Chmod(targetDir, 0o755)).To(Succeed())
				}()

				_, err := gobootutils.CreateRootDir(targetDir, "nowrite")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("CloseFileWithErr", func() {
		Context("when closing a valid file", func() {
			It("closes without panicking", func() {
				// Create a temporary file
				file, err := os.CreateTemp(tempDir, "closefile-test-*")
				Expect(err).NotTo(HaveOccurred())

				// This should not panic
				Expect(func() {
					gobootutils.CloseFileWithErr(file)
				}).NotTo(Panic())
			})
		})

		Context("when closing an already closed file", func() {
			It("ignores the error but does not panic", func() {
				file, err := os.CreateTemp(tempDir, "closefile-test-*")
				Expect(err).NotTo(HaveOccurred())

				// Close it once
				err = file.Close()
				Expect(err).NotTo(HaveOccurred())

				// Closing again should not panic
				Expect(func() {
					gobootutils.CloseFileWithErr(file)
				}).NotTo(Panic())
			})
		})
	})

})
