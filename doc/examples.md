# Configuration Examples (Input -> Output)

These examples show what to configure and what output to expect.

## Example 1: Minimal base project only

### Input (`services` in `goboot.yml`)

```yaml
services:
  - id: "base_project"
    confPath: "./configs/base_project.yml"
    enabled: true
```

### Expected output

- base project structure and core files
- no lint/test/local/CI service artifacts

Typical files:

- `go.mod`
- `README.md`
- `cmd/<project>/main.go`

## Example 2: Developer workflow scaffold (lint + test + local)

### Input (`services` in `goboot.yml`)

```yaml
services:
  - id: "base_project"
    confPath: "./configs/base_project.yml"
    enabled: true
  - id: "base_lint"
    confPath: "./configs/base_lint.yml"
    enabled: true
  - id: "base_test"
    confPath: "./configs/base_test.yml"
    enabled: true
  - id: "base_local"
    confPath: "./configs/base_local.yml"
    enabled: true
```

### Expected output

- lint configs (`.golangci.yml`, `.yamllint.yml`, etc. depending on enabled linters)
- test scaffolding (Ginkgo or stdlib based on `useStyle`)
- local workflow files (`Makefile`, `Taskfile.yml`, `scripts/*.sh`) for enabled outputs

## Example 3: CI generation (`balanced` vs `strict`)

### Input (`base_ci.yml`)

```yaml
imagePolicy: "balanced"
goVersions: ["1.25", "1.26"]
autoBranches: ["main", "master"]
```

### Expected output (`balanced`)

- provider CI files generated from `templates/ci_base/<provider>/...`
- readable tag-based image references by default
- policy-aware includes/job files for enabled commands

### Input change (`base_ci.yml`)

```yaml
imagePolicy: "strict"
```

### Expected output (`strict`)

- digest-oriented references where provider template supports it
- explicit placeholders (for example `REPLACE_ME`) where user pinning is required

## Validation checklist after generation

- generated project has no unrendered template markers (`{{ ... }}`)
- expected service files exist for enabled services
- generated checks run (`make lint`, `make test`, CI syntax where applicable)
