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
- `configs/base_local.yml`
- `configs/base_ci.yml`

Set in `configs/goboot.yml`:

- `projectName`
- `repoUrl`
- `gitProvider`
- `targetPath`

## 3. Run goboot

```bash
go run ./cmd/goboot --config ./configs/goboot.yml
```

## 4. Check expected output

If `projectName` is `IntroProject` and `targetPath` is `outputs`, expect:

- `outputs/IntroProject/go.mod`
- `outputs/IntroProject/README.md`
- `outputs/IntroProject/Makefile` (when `base_local` enabled)
- `outputs/IntroProject/scripts/lint.sh` (when `base_local` + lint registration enabled)
- `outputs/IntroProject/.golangci.yml` (when `base_lint` enabled)
- `outputs/IntroProject/.github/workflows/*.yml` or `outputs/IntroProject/.gitlab-ci.yml` (when `base_ci` enabled)

## 5. Iterate safely

Change one config at a time, regenerate, and diff output.

Suggested loop:

1. edit one config field
2. run `go run ./cmd/goboot --config ...`
3. inspect generated files
4. run generated project checks (`make lint`, `make test` if present)

## Notes

- `goboot` is pre-alpha and intentionally explicit.
- The goal is deterministic scaffolding, not one-click hidden behavior.
