# ADR-036: Tag-Driven Release Automation

**Status:** Accepted  
**Milestone:** v0.2.1

## Context

The containerization milestone produces buildable CLI projects, but releases
still require manual cross-platform builds, archive naming, checksums, release
notes, and provider uploads.

## Decision

Add an optional `base_release` service. It owns GoReleaser configuration and
release instructions. Provider workflow files remain owned by `base_ci`; the
release service registers its job through the existing CI registrar.

Versions are sourced exclusively from annotated semantic-version Git tags.
The tool will not infer or commit a version bump. GoReleaser creates Linux,
macOS, and Windows AMD64/ARM64 binaries, archives, checksums, changelog content,
and provider releases.

## Consequences

- Release output is reproducible and reviewable before tagging.
- GitHub and GitLab use the same service contract.
- Publishing credentials remain provider-managed secrets.
- Maintainers retain explicit control over version selection.
- A bad published release is corrected by a new version, not a replaced tag.
