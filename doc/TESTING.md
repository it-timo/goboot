# Testing Guide for goboot

This guide describes testing approach, structure, and quality gates for goboot.

goboot uses **Ginkgo v2** and **Gomega** for BDD-style tests across packages.

## Coverage Targets

The project enforces strict quality gates to ensure long-term maintainability:

- **Overall project coverage**: >80%
- **Critical services** (`baseproject`, `config`, `goboot`): >90%
- **Helpers & Utils**: >85%
- **All tests**: Must pass with `-race` enabled

## Running Tests

### Quick Start

```bash
# Run all tests (race + coverage, filtered to real packages)
make test

# Using go-task (if installed)
task test
```

`make test` and `task test` use `go list ./...` and intentionally exclude
`/test/noauto` and `/templates`.

### Detailed Commands

```bash
# Run all package tests directly (no filters)
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage profile
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out
```

### Running Specific Package Tests

```bash
# Test gobootutils package
go test -v ./pkg/gobootutils

# Test goboottypes package
go test -v ./pkg/goboottypes

# Test config package
go test -v ./pkg/config
```

### Updating Golden Digests (E2E)

`cmd/goboot/main_e2e_test.go` contains deterministic digest assertions for
selected generated outputs.

Update them only when output changes are intentional.

```bash
# Run the relevant E2E specs and read the failing digest from test output
go test ./cmd/goboot -ginkgo.focus "scaffolds a full project|supports go-style tests|generates deterministic output"

# After updating expected digest constants, verify everything
go test ./...
make lint
```

Rules:

- Do not update digest constants to "make tests pass" without reviewing file diffs.
- Confirm the changed generated files match the intended contract/policy changes.
- Keep digest inputs stable (the tests intentionally hash a fixed high-signal file set).

## Test Organization

### Suite Structure

Each package follows this structure:

``` bash
pkg/
└── packagename/
    ├── packagename.go
    ├── packagename_suite_test.go  # Ginkgo suite bootstrap
    └── packagename_test.go        # Test specs
```

### Test Naming Conventions

We follow BDD-style naming with Ginkgo's `Describe`, `Context`, and `It` blocks:

```go
Describe("FunctionName", func() {
    Context("when specific condition", func() {
        It("does expected behavior", func() {
            // Test implementation
        })
    })
})
```

## Package Test Coverage

### pkg/gobootutils

Tests filesystem operations, templates, and path utilities:

- ✅ `EnsureDir` - Directory creation with path escape prevention
- ✅ `ComparePaths` - Path comparison and normalization
- ✅ `CreateRootDir` - Root directory creation
- ✅ `CloseFileWithErr` - Safe file closing
- ✅ `ExecuteTemplateText` - Template rendering
- ✅ `RenderTemplateToFile` - File template rendering

**Security**: Includes tests for path traversal prevention (`../` attacks)

### pkg/goboottypes

Tests constants and type definitions:

- ✅ Default linter commands validation
- ✅ Linter identifiers naming
- ✅ Script names and conventions
- ✅ Service name conventions
- ✅ File permission constants

### pkg/config

Tests configuration management:

- ✅ Manager initialization
- ✅ Service registration and retrieval
- ✅ Registrar management
- ✅ Unregistration operations
- ✅ ServiceConfigMeta operations
- ✅ Loader and validator logic (full coverage)

## Writing New Tests

### Best Practices

1. **BDD Style**: Use descriptive `Describe`, `Context`, and `It` blocks
2. **Table-Driven Tests**: Use `DescribeTable` for testing multiple scenarios
3. **Setup/Teardown**: Use `BeforeEach`/`AfterEach` for test isolation
4. **Assertions**: Use Gomega matchers for readable assertions
5. **Mock Usage**: Create simple mock implementations for interfaces

### Example Test

```go
var _ = Describe("MyFunction", func() {
    var (
        input  string
        result string
    )

    BeforeEach(func() {
        input = "test"
    })

    Context("when input is valid", func() {
        It("returns expected output", func() {
            result = MyFunction(input)
            Expect(result).To(Equal("expected"))
        })
    })
})
```

### Table-Driven Tests

```go
DescribeTable("validation tests",
    func(input string, expected bool) {
        result := Validate(input)
        Expect(result).To(Equal(expected))
    },
    Entry("valid input", "valid", true),
    Entry("invalid input", "", false),
)
```

## Test Helpers and Utilities

### Temporary Directories

For filesystem tests, use temporary directories:

```go
BeforeEach(func() {
    var err error
    tempDir, err = os.MkdirTemp("", "test-*")
    Expect(err).NotTo(HaveOccurred())
})

AfterEach(func() {
    os.RemoveAll(tempDir)
})
```

### Mock Implementations

Create simple mocks for testing:

```go
type mockService struct {
    id          string
    shouldError bool
}

func (m *mockService) ID() string {
    return m.id
}

func (m *mockService) Run() error {
    if m.shouldError {
        return errors.New("mock error")
    }
    return nil
}
```

## Coverage Goals

- **Overall Project**: 80%+
- **Critical Packages** (utils, config): 90%+
- **Type Packages**: 85%+
- **Service Packages**: 80%+

## Continuous Integration

Tests can be run in CI using the same entrypoints:

```yaml
- name: Run linters
  run: make lint

- name: Run tests
  run: make test
```

## Troubleshooting

### Common Issues

#### Tests fail with "no test files"

- Ensure test files end with `_test.go`
- Ensure package name is `packagename_test`

#### Import cycle errors

- Use `packagename_test` package name
- Import the package under test explicitly

#### Coverage not showing

- Ensure you're using `-cover` flag
- Check that test files are in the same directory as source files

## Additional Resources

- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Gomega Matchers](https://onsi.github.io/gomega/)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
