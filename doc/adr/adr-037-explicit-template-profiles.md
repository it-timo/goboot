# ADR-037: Explicit Template Profiles

**Status:** Accepted  
**Milestone:** v0.3.0

## Context

Projects need different initial quality baselines, but silently changing enabled
services or overwriting explicit YAML would violate goboot's deterministic,
auditable configuration model.

## Decision

Add a validated root `profile` field with four values: `minimal`, `standard`,
`enterprise`, and `oss`. Omission selects `standard`.

Profile-aware configs receive the selected value before validation. Profiles may
fill missing values and alter owned template baselines, but explicit service
settings take precedence. Profiles never add, remove, enable, or disable
services.

The initial profile surface covers:

- Go linter selection and cyclomatic-complexity thresholds
- default test style and command
- generated profile documentation

## Consequences

- Teams can select a meaningful baseline with one explicit root value.
- Generated output records the selected profile.
- Existing configurations remain compatible through the `standard` default.
- Service composition remains explicit.
- Adding future profile behavior requires an owned, profile-aware config rather
  than global mutation.
