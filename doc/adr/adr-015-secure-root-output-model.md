# ADR-015: Scoped Filesystem Output via `os.Root`

**Tags:** `filesystem`, `security`, `sandboxing`, `go-1.23+`

---

## Status

Accepted

---

## Context

Project generation performs many filesystem writes.
Without path confinement, traversal paths and accidental overwrite outside the target directory are possible.

---

## Decision

Use `os.OpenRoot()` and perform generation writes through `*os.Root` scoped operations.
Avoid unbounded host-path writes in service rendering paths.

---

## Advantages

- Constrains writes to an explicit project root.
- Reduces traversal risk in template-driven path rendering.
- Aligns runtime behavior with secure-by-default filesystem handling.

---

## Disadvantages

- Requires newer Go runtime support and contributor familiarity.
- Existing `os`/`filepath` helper code may need adaptation.
- Integration tests must account for root-scoped behavior.

---

## Alternatives Considered

- **Raw host filesystem writes with manual sanitization:** rejected due to higher risk of incomplete path checks.
- **In-memory filesystem abstraction only:** rejected because real generated output still requires host filesystem behavior.
