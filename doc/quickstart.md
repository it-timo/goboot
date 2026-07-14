# Quickstart (No Code Reading Required)

This guide is for first-time users who want to run `goboot` and inspect generated output.

## 1. Pick a target directory

Example:

```bash
mkdir -p /tmp/goboot-demo
```

## 2. Use the provided configs as a baseline

Start from:

- `configs/goboot.yml`
- `configs/base_project.yml`
- `configs/base_lint.yml`
- `configs/base_test.yml`
- `configs/base_logger.yml`
- `configs/base_docker.yml`
- `configs/base_release.yml`
- `configs/base_governance.yml`
- `configs/base_local.yml`
- `configs/base_ci.yml`

Set in `configs/goboot.yml`:

- `projectName: "IntroProject"`
- `repoUrl: "https://github.com/projects"`
- `gitProvider: "gitlab"`
- `targetPath: "/tmp/goboot-demo"`
- `profile: "standard"`

## 3. Run goboot

```bash
go run ./cmd/goboot --config ./configs/goboot.yml
```

## 4. Check expected output

If `projectName` is `IntroProject` and `targetPath` is `/tmp/goboot-demo`, expect:

- `/tmp/goboot-demo/IntroProject/go.mod`
- `/tmp/goboot-demo/IntroProject/README.md`
- `/tmp/goboot-demo/IntroProject/PROFILE.md`
- `/tmp/goboot-demo/IntroProject/Makefile` (when `base_local` enabled)
- `/tmp/goboot-demo/IntroProject/scripts/lint.sh` (when `base_local` + lint registration enabled)
- `/tmp/goboot-demo/IntroProject/.golangci.yml` (when `base_lint` enabled)
- logger-aware CLI/service wiring for `slog` or `zerolog` (when `base_logger` enabled)
- `/tmp/goboot-demo/IntroProject/Dockerfile` and
  `/tmp/goboot-demo/IntroProject/docker-compose.yml` (when `base_docker` enabled)
- `/tmp/goboot-demo/IntroProject/.github/workflows/*.yml` or
  `/tmp/goboot-demo/IntroProject/.gitlab-ci.yml` (when `base_ci` enabled)
- `/tmp/goboot-demo/IntroProject/CODEOWNERS`, `CONTRIBUTING.md`, and
  provider-native contribution templates (when `base_governance` enabled)

`base_logger` does not overwrite project files directly. It provides validated logger
settings to `base_project`, which owns the generated runtime code.

`base_docker` packages the generated CLI-style application. It does not assume
the project is a network service, so compose port mappings are empty by default.

`base_governance` writes policy files only. Provider-side branch protection,
approval enforcement, and private vulnerability reporting remain administrator settings.

## 5. Iterate safely

Change one config at a time, regenerate, and diff output.

Suggested loop:

1. edit one config field
2. run `go run ./cmd/goboot --config ...`
3. inspect generated files
4. run generated project checks (`make lint`, `make test`, and
    `make container-check` when Docker is enabled)

## Notes

- `goboot` is pre-alpha and intentionally explicit.
- The goal is deterministic scaffolding, not one-click hidden behavior.
