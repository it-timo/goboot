# ADR-011: Centralized Config Dispatch via `createServiceConfig()`

**Tags:** `dispatch`, `registration`, `no-reflection`

---

## Status

Accepted

---

## Context

The manager needs a deterministic mapping from service IDs to config types.
Reflection-based dispatch increases runtime complexity and obscures allowed config set.

---

## Decision

Use a hardcoded factory switch (`createServiceConfig`) that maps known service IDs to concrete config structs.
Unknown IDs return `nil`.

---

## Advantages

- Supported config types are explicit.
- Startup behavior is deterministic.
- Easier static analysis and review of allowed modules.

---

## Disadvantages

- Every new service requires a factory update.
- Not suitable for runtime plugin ecosystems.

---

## Alternatives Considered

- **Reflection-based constructor lookup:** rejected due to lower transparency and higher runtime risk.
- **Map registry populated at init time:** rejected because registration side
  effects are less explicit than one central switch.
