# ADR-004: BDD Testing Standards and Tooling

**Tags:** `testing`, `bdd`, `standards`, `quality`

---

## Status

Accepted

---

## Context

`goboot` includes filesystem-heavy flows and configuration permutations.
Tests need readable structure, isolated setup, and explicit permutation coverage.

---

## Decision

Adopt Ginkgo/Gomega as the standard test style for this repository.

- Use `Describe`/`Context`/`It` for suite structure.
- Use `DescribeTable` for permutation coverage.
- Keep specs isolated with per-spec temporary directories.
- Avoid shared mutable global state across specs.

---

## Advantages

- Consistent test structure across packages.
- Better failure localization for matrix-like inputs.
- Clearer setup/teardown boundaries.

---

## Disadvantages

- Additional dependency in the test toolchain.
- Some contributors may prefer stdlib-only style.
- More verbose syntax for simple unit cases.

---

## Alternatives Considered

- **Stdlib-only tests (`testing` + table loops):** rejected as default due to
  lower readability in complex integration-style flows.
- **Mixed styles per package:** rejected because style variance raises review overhead and inconsistency.
