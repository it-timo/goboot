# Stable CLI and Configuration Contract

Version 0.8 freezes the v1 candidate command-line flags, root YAML fields, and
built-in service YAML fields. The freeze makes automation safe to build before
v1 while still permitting backward-compatible additions.

## Commands and flags

Goboot remains a flag-oriented command with these public operations:

| Operation | Flags | Target changes |
| --- | --- | --- |
| Generate | `--config PATH` | Transactional project update |
| Dry run | `--config PATH --dry-run` | None |
| Validate | `--config PATH --validate` | None |
| Version | `--version` | None |

Common flags are `--log-level` and `--output`. Generation additionally accepts
`--regeneration-policy` and `--skip-go-mod-tidy`. Unsupported positional
arguments and incompatible operation flags are usage errors.

`--validate` loads the strict root YAML file, resolves every enabled service
config, rejects unknown fields, applies defaults, and runs all configuration
validators. It does not render templates, invoke `go mod tidy`, create staging
directories, or inspect the target project.

## Output contract

`--output human` is the default. `--output json` writes exactly one JSON object
to stdout. Structured logs remain on stderr.

Successful JSON results contain:

- `status`: `success`
- `operation`: `generate`, `dry-run`, `validate`, or `version`
- `version`: the goboot build version
- `message`: a stable human-readable summary
- optional `config`, `project`, and dry-run `plan` fields

Error JSON results contain:

- `status`: `error`
- `operation` and `version`
- numeric `exit_code`
- `category`: `internal`, `usage`, `config`, `generation`, or `conflict`
- `message`
- an optional conflict `plan`

Consumers should branch on `status`, `category`, and `exit_code`, not message
text. New optional JSON fields may be added without a breaking change.

## Exit codes

| Code | Category | Meaning |
| --- | --- | --- |
| 0 | success | Requested operation completed |
| 1 | internal | Unclassified internal failure |
| 2 | usage | Invalid flags, combinations, or positional arguments |
| 3 | config | Root or enabled-service configuration is invalid |
| 4 | generation | Rendering, dependency, filesystem, or apply failure |
| 5 | conflict | Safe regeneration detected user-owned changes |

## Schemas and editor completion

[`../schemas/`](../schemas/) contains JSON Schema draft 2020-12 documents for
the root config and every built-in service. Schemas reject unknown properties
and publish enums and basic constraints. Runtime validation remains authoritative
for cross-file and service-dependent rules.

Every example under [`../configs/`](../configs/) begins with a
`yaml-language-server` schema modeline, enabling completion and diagnostics in
compatible editors without committing editor-specific workspace settings.

## Compatibility matrix

The compatibility workflow builds the CLI and tests portable config and
regeneration contracts on Ubuntu 24.04, macOS 15, and Windows 2025 with the
repository's pinned Go version. Platform-specific packaging remains covered by
GoReleaser.

## Deprecation and migration policy

- Existing flags, exit meanings, JSON field meanings, and YAML fields will not
  be removed or repurposed before v1.
- Backward-compatible optional fields and enum values may be added in minor
  releases.
- A post-v1 incompatible change requires a deprecation notice in the changelog
  and documentation for at least two minor releases before removal.
- Deprecation notices go to stderr so JSON stdout remains parseable.
- Each removal must include an old-to-new example and an automated migration or
  deterministic manual steps.
- Unknown YAML fields continue to fail immediately; goboot never silently drops
  obsolete configuration.

The schemas, examples, changelog, and runtime decoder must change together in
the same pull request.
