# ADR-041: Manifest-Based Transactional Regeneration

**Tags:** `regeneration`, `ownership`, `filesystem`, `transactions`, `cli`

---

## Status

Accepted — 2026-07-14

## Context

A stable generator must support repeated runs after users edit a generated
repository. Directly writing templates into the target cannot reliably
distinguish goboot output from user-owned files, and a failed multi-service run
can otherwise leave a partially updated project.

Generic text merging is not safe across Go, YAML, Markdown, shell, Make, and
provider configuration formats. Timestamp-based ownership is also insufficient
because copies, checkouts, and formatting tools routinely alter timestamps.

## Decision

- Generate the complete requested project in a temporary staging tree.
- Record owned file paths, SHA-256 digests, and permission modes in
  `.goboot-manifest.yml` together with generator inputs.
- Default to a `managed` policy that updates only unchanged owned files and
  reports all other collisions.
- Provide explicit `replace` and `preserve` policies; do not provide a generic
  merge policy.
- Detect changed content, changed modes, stale owned files, unowned collisions,
  malformed manifests, path-type collisions, symbolic links, and special files.
- Add `--dry-run` to render and plan without changing the configured target.
- Build the final candidate beside the target, then use same-filesystem renames
  plus a backup to commit or roll back the whole project tree.
- Preserve unowned and policy-preserved files in successful candidates while
  excluding them from the next ownership manifest.

## Consequences

- Default repeated generation fails safely instead of silently overwriting user
  modifications.
- A committed manifest becomes part of the generated repository contract.
- Full staging uses additional temporary disk space approximately equal to the
  generated project plus the copied current project.
- Dry runs can still execute dependency resolution in staging so their bytes
  match a real run.
- Directory replacement uses two portable renames. There can be a short interval
  where the backup holds the old project path, but a failed commit rename restores
  it.
- Format-specific merge behavior may be introduced later behind an explicit
  service contract.

## Alternatives considered

### Always overwrite generated paths

This preserves the old implementation's simplicity but can destroy user work and
is not acceptable for a stable regeneration contract.

### Infer ownership from known template paths

Template paths change between versions and do not prove whether a generated file
was edited after creation.

### Merge every text file

Line-oriented merging cannot preserve the semantics of every generated format
and would create ambiguous conflict behavior.

### Modify the target and roll back individual files

Per-file journals are more complex and expose partial state during generation.
Building a complete candidate keeps target mutations at the final commit boundary.
