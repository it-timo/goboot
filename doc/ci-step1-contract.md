# CI Step 1 Contract (Config + Templates)

## Purpose

Define the CI config/template contract before changing runtime wiring in `base_ci`.

This step defines:

- Which config keys are official.
- Which template variables are official.
- How image pinning should work for generated repositories.

## Current State Summary

- `configs/base_ci.yml` defines `sourcePath`, `goVersions`, `imagePolicy`, `autoBranches`, and `jobs`.
- `pkg/config/baseCI.go` supports typed validation/defaults for `goVersions`,
  `autoBranches`, `imagePolicy`, and job-level `allowFailure`.
- `pkg/baseci/baseCI.go` renders provider-specific templates via
  `sourcePath/<provider>/...` and passes `GoVersions`, `AutoBranches`,
  `ImagePolicy`, and registered file/job command data into templates.
- GitLab and GitHub templates now consume dynamic command registration and policy-driven behavior.

Result: runtime wiring is aligned with the v1 contract and covered by provider/policy tests.

## Contract v1 (Current)

### Config Keys

Supported keys in `base_ci.yml`:

- `sourcePath` (required): root directory for CI templates.
- `goVersions` (required): list of Go versions for matrix jobs.
- `autoBranches` (optional): exact branch names that run CI automatically. Default: `["main","master"]`.
- `jobs` (optional): CI-only jobs with command lists.
- `imagePolicy` (optional): `balanced` (default), `strict`, or `simple`.

`autoBranches` in v1 is exact-match only (no regex/glob). Pattern support is deferred to a later version.

`dockerImage` should not be a top-level contract key in v1.
Reason: it is provider-specific and should be handled via generated provider files or policy presets.

### Job Model

`jobs` entries should support:

- `commands` (required): list of non-empty shell commands.
- `allowFailure` (optional, default `false`): marks generated job as non-blocking where provider supports it.

Any additional job-level metadata is out of scope for v1 and can be added in v1.1.

### Template Data Object

Official render fields passed to CI templates:

- `ProjectName`
- `CIDir`
- `JobScripts`
- `FileScripts`
- `EnabledJobFiles`
- `GoVersions`
- `AutoBranches`
- `ImagePolicy`

Optional (only when used by templates):

- `GoModulePath`

### Provider Path Rule

`sourcePath` points to CI template root, with provider subfolders:

- `sourcePath/gitlab/...`
- `sourcePath/github/...`

## Image Pinning Policy

### `balanced` (recommended default)

- Generated repos use readable tags by default.
- Optional generated companion file can provide digest-pinned variants for teams that want hard pinning.
- Best onboarding and maintenance for most users.

### `strict`

- Generated repos use digest-pinned images by default.
- Highest reproducibility/security, highest maintenance overhead.

### `simple`

- Generated repos use tags only and no digest artifacts.
- Lowest friction, lowest reproducibility guarantees.

## Recommendation

Use `balanced` as default for generated repositories.

Reasoning:

- Your maintainer repo can stay strict and digest-pinned.
- Generated repos are usually maintained by teams that optimize for clarity and speed first.
- This keeps security hardening available without forcing it on every bootstrap user.

## Naming Rules

- Config keys: lowerCamelCase (existing project style).
- Template render fields: exported Go field names.
- CI env vars: upper snake case.
- Version-derived CI variable names must normalize dots: `1.26` -> `1_26`.
  - Example: `GO_1_26_DIGEST`.
- Generated output ordering must be deterministic:
  - sort enabled job files
  - stable include order
  - stable matrix order from `goVersions`

## Template Intent Rules

### GitLab templates

- `.gitlab-ci.yml.tmpl` should include only enabled job files and required shared provider files.
- Job templates must consume `GoVersions` and `AutoBranches` from render data.
- Digest variable references must use normalized names.
- Provider shared files are contract-bound for GitLab:
  - `versions.yml` (image/version variables)
  - `commands.yml` (shared command helpers), when applicable by policy.

Policy behavior:

- `strict`: digest-based variables/files required.
- `balanced`: tag-based jobs by default, optional digest companion artifacts.
- `simple`: tag-based jobs only, no digest companion artifacts.

### GitHub templates

- Remove leftover invalid template fragments from lint workflow template.
- Test/build templates must iterate all registered commands, not first command only.
- Matrix Go versions should come from `GoVersions`.
- Keep pinned action SHAs for actions usage.
- In `strict` mode, lint container images must resolve from digest-based env vars in workflow scope.
  - Note: GitHub Actions runtime differs from GitLab; "strict" here means digest-pinned container image refs where supported.

## Acceptance Criteria for Step 1

- No config keys exist without a typed contract definition.
- No template references non-existent render fields.
- Provider path semantics are unambiguous (`sourcePath/<provider>/...`).
- Image pinning defaults are explicitly documented.
- CI variable naming rules are explicit and consistent.
- Policy behavior is explicit per provider (`strict`, `balanced`, `simple`).
- Output ordering is deterministic for stable diffs.

## Deferred (Post-v1)

- `autoBranches` regex/glob support.
- monorepo-specific template data (e.g., `IsMonorepo`).
- generation timestamp metadata (`GeneratedAt`).
- external config schema publication (JSON Schema/CUE) for IDE autocomplete.

This document records the current v1 CI contract. Operator-facing commands and
canary behavior live in [`ci.md`](./ci.md).
