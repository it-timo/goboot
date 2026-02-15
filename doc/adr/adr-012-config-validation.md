# ADR-012: Typed Config Structure and Validation Strategy

**Tags:** `config`, `validation`, `typed-structure`

---

## Status

Accepted

---

## Context

Untyped configuration decoding increases runtime error risk and weakens tooling support.
`goboot` requires explicit config contracts per service.

---

## Decision

- Use dedicated config structs per service.
- Validate configs centrally during load.
- Pass only validated config instances into service lifecycle.

---

## Advantages

- Better compile-time and editor support.
- Earlier detection of invalid config data.
- Clear place for service-specific validation rules.

---

## Disadvantages

- More boilerplate for each new config model.
- Schema evolution requires synchronized struct and validation updates.

---

## Alternatives Considered

- **Unstructured config maps:** rejected due to reduced safety and maintainability.
- **Late validation inside each service run:** rejected due to delayed failure and inconsistent behavior.
