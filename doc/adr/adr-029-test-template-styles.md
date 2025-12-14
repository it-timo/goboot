# 📄 ADR-029: Test Template Styles — Ginkgo by Default, Stdlib as Opt-Out

**Tags:** `testing`, `templates`, `ginkgo`, `stdlib`, `flexibility`

---

## Status

✅ Accepted

---

## Context

Teams using goboot have different testing preferences.
We want a **BDD-first** default to match the project’s internal tests,
but some consumers mandate pure stdlib tests to avoid extra dependencies.
The new `baseTest` service needs a deterministic way to select templates and commands without branching logic in user code.

---

## Decision

- Add a `useStyle` flag to `BaseTestConfig` with two supported values:
  - `ginkgo` (default): render BDD suites using Ginkgo/Gomega.
  - `go`: render plain `testing`-only files.
- Template tree `templates/test_base/` includes both variants. During the path render pass,
`suite_test.go` files are **skipped when `useStyle` = `go`** to avoid pulling Ginkgo.
- Content rendering injects `RepoImportPath`, `ProjectName`, and casing helpers into imports and comments for either style.
- `testCmd` defaults to a race-enabled, coverage-enabled invocation suitable for both styles
and is registered into scripts when `base_local` is present.
- Future styles (e.g., testify) must add new template variants rather than branching runtime code.

---

## Advantages

- Keeps BDD as the opinionated default while respecting teams that want zero third-party deps.
- Deterministic generation: style choice is explicit in config, not inferred or mixed.
- Simplifies future expansions by isolating style-specific files in templates.

---

## Disadvantages

- Two template variants increase maintenance overhead until we add tooling to validate them.
- Non-default styles still share one `testCmd`; some projects may need to override it manually.

---

## Alternatives Considered

- **Single Ginkgo-only path:** Blocks adoption in environments that disallow external test deps.
- **Runtime flag to toggle assertions inside one template:** Would create brittle conditional code
and harder-to-read generated tests.
