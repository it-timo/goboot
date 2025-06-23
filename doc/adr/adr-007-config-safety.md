# 📄 ADR-007: Config Manager Behavior and Safety

**Tags:** `manager`, `validation`, `static-analysis`

---

## Status

✅ Accepted

---

## Context

The config manager is responsible for maintaining the set of enabled service configs during a goboot run.
Improper or invalid configuration could lead to incorrect project generation or crashes.

---

## Decision

The `Manager` only allows registration of configs that:

1. Pass `Validate()`
2. Have a non-empty `ID()`
3. Are registered under a unique key (by ID)

Duplicate or invalid entries are rejected with an error. Lookup is always explicit by ID.

---

## Advantages

- Prevents accidental misregistration or overlap
- Ensures all configs are validated before use
- Enables safe concurrent reading (future-safe)

---

## Disadvantages

- Requires each config to define its own strict `Validate()` logic
- It Does not allow multiple configs of the same type (by design)
