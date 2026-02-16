# ADR-007: Centralized Service Name Registry in `goboottypes/names.go`

**Tags:** `constants`, `service-names`, `structure`, `decoupling`

---

## Status

Accepted

---

## Context

Service IDs are used in configuration, registration, and execution paths.
Repeated string literals increase typo risk and make refactors harder.

---

## Decision

Define service IDs as constants in a single shared location (`pkg/goboottypes/names.go`).

---

## Advantages

- Reduces mismatch risk between config keys and runtime IDs.
- Improves discoverability of supported services.
- Makes rename/refactor operations safer.

---

## Disadvantages

- Introduces shared dependency on the constants package.
- Requires central-file updates when adding services.

---

## Alternatives Considered

- **Inline string IDs per package:** rejected due to duplication and typo risk.
- **Generated IDs at runtime:** rejected because explicit identifiers are easier to audit and test.
