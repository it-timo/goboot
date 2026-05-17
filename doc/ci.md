# CI Guide

This document explains CI behavior for:

- the `goboot` maintainer repository
- generated projects produced by `base_ci`

## Two CI Contexts

### Maintainer CI (`goboot` repository)

- Hardened and curated by maintainers.
- Can use pinned images and project-specific tooling conventions.
- Optimized for repository governance and internal quality gates.

### Generated CI (output projects)

- Produced by the `base_ci` service.
- Designed for reusable bootstrap behavior across projects.
- Policy-driven via `imagePolicy`.

## `base_ci` Config Keys

`configs/base_ci.yml` supports:

- `sourcePath` (required)
- `goVersions` (required)
- `autoBranches` (optional; exact names, default `main/master`)
- `jobs` (optional; supports `commands` and `allowFailure`)
- `imagePolicy` (optional; `balanced` default, `strict`, `simple`)

## Image Policies

### `balanced` (default)

- Uses readable tags in generated jobs where practical.
- Intended as the default usability profile.

### `strict`

- Uses digest-oriented references where supported by provider/runtime model.
- Generates `REPLACE_ME` digest placeholders in strict variable blocks.
- Manual replacement is expected for generated repos.

### `simple`

- Uses tags and keeps generation minimal.
- Lowest operational overhead, lower reproducibility guarantees.

## Provider Behavior

### GitLab

- Generates `.gitlab-ci.yml` plus `.gitlab/ci/*` files.
- Includes policy-aware `versions.yml` and shared `commands.yml`.
- Strict mode uses digest-oriented variables including Go runtime and lint images.

### GitHub

- Generates `.github/workflows/*`.
- Keeps pinned action SHAs for action usage.
- Strict mode provides digest env placeholders for container lint image references where supported.

## `REPLACE_ME` Placeholders

`REPLACE_ME` is intentional:

- avoids hidden runtime fetch logic ("no magic"),
- avoids hardcoded, fast-stale digest data,
- keeps strict mode explicit and auditable.

Users choosing `strict` are expected to fill digest values before production use.

## Determinism Guarantees

- Enabled CI job files are sorted before rendering.
- Include lists are stable.
- `goVersions` order controls matrix order.

## Related Docs

- `doc/ci-step1-contract.md`
- `doc/adr/adr-033-ci-policy-and-provider-layout.md`
- `configs/base_ci.yml`

## Canary Verification (Generated CI)

Before tagging a release, run local CI simulation for both root and generated output:

```bash
make verify_ci_canary
```

Direct script usage:

```bash
./scripts/verify_ci_canary.sh
```

Skip regeneration if `outputs/IntroProject` is already present:

```bash
./scripts/verify_ci_canary.sh --skip-generate
```

Select provider mode explicitly:

```bash
./scripts/verify_ci_canary.sh --provider=github
./scripts/verify_ci_canary.sh --provider=gitlab
./scripts/verify_ci_canary.sh --provider=both
```

Notes:

- The flow runs the same `act` and `gitlab-ci-local` command sequence in:
  - repository root
  - `outputs/IntroProject`
- By default (`--provider=config`), provider selection is derived from `configs/goboot.yml` (`gitProvider`).
- This verifies provider workflow behavior without pushing to remote canary repositories.
- This flow is not fully dockerized: it requires host-installed `act` and `gitlab-ci-local`, plus Docker daemon/socket access.
- Running inside heavily containerized or restricted environments can fail due to socket/privilege limits.
- `act` runs in offline mode and with `--use-gitignore=false` to avoid missing dependency files such as `go.sum`.
- `act` and GitLab canary image pre-pulls run by default; use `--skip-prepull` to disable
them and `--refresh-images` to force fresh pulls.
