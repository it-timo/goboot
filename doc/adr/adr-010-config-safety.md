# ADR-010: Config Manager Behavior and Safety

**Tags:** `manager`, `validation`, `static-analysis`

---

## Status

Accepted

---

## Context

The config manager is the gate for service execution input.
Invalid or duplicate registrations can lead to ambiguous runtime behavior.

---

## Decision

Allow config registration only when:

1. `Validate()` succeeds,
2. `ID()` is non-empty, and
3. the ID is unique in the manager.

Reject invalid or duplicate entries with explicit errors.

---

## Advantages

- Prevents ambiguous service-to-config mapping.
- Moves failure to startup time instead of execution time.
- Keeps lookup behavior explicit by service ID.

---

## Disadvantages

- Validation quality depends on each config module implementation.
- Duplicate-by-design config patterns are not supported.

---

## Alternatives Considered

- **Allow duplicate IDs with merge behavior:** rejected due to ambiguity and harder debugging.
- **Lazy validation at service runtime:** rejected because failures would occur later and be harder to localize.
