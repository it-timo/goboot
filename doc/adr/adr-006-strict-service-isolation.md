# ADR-006: Strict Service Isolation

**Tags:** `architecture`, `modularity`, `dependencies`

---

## Status

Accepted

---

## Context

Direct service-to-service imports can create tight coupling and make services hard to evolve independently.
`goboot` uses multiple services with distinct responsibilities, so dependency boundaries must be explicit.

---

## Decision

Enforce strict service isolation.

- Service packages must not import other service packages directly.
- Shared contracts belong in common packages (`config`, `goboottypes`).
- Cross-service cooperation is orchestrated by `pkg/goboot` through injected interfaces.

---

## Advantages

- Services remain independently testable.
- Dependency graph remains simpler to reason about.
- Responsibility boundaries are explicit.

---

## Disadvantages

- Additional interface and orchestration code is required.
- Cross-service features may need more upfront coordination.

---

## Alternatives Considered

- **Direct service imports:** rejected due to coupling and risk of cyclic dependencies.
- **Global mutable coordination state:** rejected because ownership and ordering become harder to reason about.
