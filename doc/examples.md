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

## Example 3: Logger-aware project scaffold

### Input (`services` in `goboot.yml`)

```yaml
services:
  - id: "base_project"
    confPath: "./configs/base_project.yml"
    enabled: true
  - id: "base_logger"
    confPath: "./configs/base_logger.yml"
    enabled: true
```

### Input (`base_logger.yml`)

```yaml
loggerType: "slog"
```

Supported values are:

- `slog`
- `zerolog`

### Expected output

- CLI logging setup in `cmd/<project>/main.go`
- service logger injection in `pkg/<project>/<project>.go`
- `--log-level` support for generated CLIs
- `zerolog` dependency in `go.mod` only when `loggerType: "zerolog"`

`base_logger` is a settings provider. Runtime project files remain owned by
`base_project`, so logger selection does not require a second service to overwrite
files it does not own.

## Example 4: CI generation (`balanced` vs `strict`)

### Input (`base_ci.yml`)

```yaml
imagePolicy: "balanced"
goVersions: ["1.26.3"]
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

## Example 5: CLI container packaging

### Input (`services` in `goboot.yml`)

```yaml
services:
  - id: "base_project"
    confPath: "./configs/base_project.yml"
    enabled: true
  - id: "base_docker"
    confPath: "./configs/base_docker.yml"
    enabled: true
  - id: "base_local"
    confPath: "./configs/base_local.yml"
    enabled: true
  - id: "base_ci"
    confPath: "./configs/base_ci.yml"
    enabled: true
```

### Input (`base_docker.yml`)

```yaml
sourcePath: "./templates/docker_base"
goVersion: "1.26.3"
runtimeImage: "gcr.io/distroless/static-debian12:nonroot"
fileList:
  - dockerfile
  - compose
  - dockerignore
ports: []
```

### Expected output

- `Dockerfile` with a Go build stage and non-root runtime stage
- `docker-compose.yml` for local CLI-container execution
- `.dockerignore`
- Docker helper targets in generated `Makefile` and `Taskfile.yml`
- `scripts/docker.sh` when script generation is enabled
- provider CI container job when `base_ci` is enabled

The scaffolded app is currently CLI-style. Empty `ports` is intentional and
does not indicate a missing service port.

## Validation checklist after generation

- generated project has no unrendered template markers (`{{ ... }}`)
- expected service files exist for enabled services
- generated checks run (`make lint`, `make test`, `make container-check` when
  Docker is enabled, and CI syntax where applicable)
