# ADR-027: Core Testing Strategy & Coverage Baseline

**Tags:** `testing`, `bdd`, `coverage`, `filesystem`, `safety`

---

## Status

Accepted

---

## Context

Generator behavior spans filesystem, templates, and orchestration.
Coverage and test style need clear repository-wide expectations.

---

## Decision

- Use Ginkgo/Gomega across packages.
- Maintain package-level coverage expectations for critical paths.
- Enforce the root overall coverage gate in canonical test entry points.
- Exercise real filesystem behavior for root/path/template flows.
- Prefer table-style specs for permutation-heavy logic.

---

## Advantages

- Broad regression protection for high-change generator surfaces.
- Clear test style consistency across packages.
- Better confidence for path and template changes.

---

## Disadvantages

- Test runtime is higher than mock-only approaches.
- Toolchain dependency on Ginkgo/Gomega remains required.

---

## Alternatives Considered

- **Smoke-tests only:** rejected because integration regressions would be missed.
- **Mostly mocks:** rejected because filesystem and render integration issues are critical to validate end-to-end.
