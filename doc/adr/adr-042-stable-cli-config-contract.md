# ADR-042: Stable CLI and Configuration Contract

**Tags:** `cli`, `config`, `schemas`, `compatibility`, `versioning`

---

## Status

Accepted — 2026-07-14

## Context

Goboot now supports safe repeated generation, but CI scripts and editors still
lack a stable validation, output, and schema contract. Error-string matching is
fragile, generation is unnecessarily expensive when only configuration needs
checking, and examples alone cannot provide editor completion.

## Decision

- Freeze the current v1 candidate flag and YAML field names.
- Add validation-only and version operations without introducing positional
  subcommands that would invalidate existing invocations.
- Define stable exit categories for usage, config, generation, and conflicts.
- Offer human output by default and exactly one JSON object on stdout when
  requested, with logs isolated on stderr.
- Publish strict draft 2020-12 schemas for the root and all built-in services,
  referenced by YAML language-server modelines.
- Require Linux, macOS, and Windows compatibility evidence in pull requests.
- Require documented deprecation and migration windows for future incompatible
  changes.

## Consequences

- CI can validate configuration without rendering or changing a target.
- Automation can use exit categories and JSON fields instead of message text.
- Editors provide completion and unknown-field diagnostics from versioned files.
- CLI and schema changes now carry a backward-compatibility obligation.
- Runtime checks remain necessary for relationships JSON Schema cannot express
  cleanly across separate files.

## Alternatives considered

### Introduce subcommands immediately

Changing `goboot --config ...` to `goboot generate --config ...` would create an
avoidable compatibility break. Independent operation flags preserve existing
calls and can remain stable through v1.

### Use logs as machine-readable output

Structured logs describe execution detail but are not a stable command result.
They remain on stderr while stdout carries one bounded result object.

### Publish only the root schema

The root file references separate service configs. Omitting their schemas would
leave most user-editable fields without completion or early diagnostics.
