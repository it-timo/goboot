# 📄 ADR-008: Centralized Config Dispatch via `createServiceConfig()`

**Tags:** `dispatch`, `registration`, `no-reflection`

---

## Status

✅ Accepted

---

## Context

Dynamic loading of unknown config types is error-prone, especially when reflection, registration maps,
or runtime hooks are involved.
Instead, the `goboot` tool makes config support explicit and central.

---

## Decision

Use a hardcoded factory method to map known config IDs to their concrete type:

```go
func createServiceConfig(id string) ServiceConfig {
    switch id {
    case types.ServiceNameBaseProject:
        return newBaseProjectConfig()
    default:
        return nil
    }
}
```

No runtime plugins. No reflection.

---

## Advantages

- Explicit control over what config types are allowed
- Impossible to inject unauthorized logic via YAML
- Easier to statically analyze and extend

---

## Disadvantages

- Requires a manual update for each new config module
- Not suitable for plugin ecosystems (which goboot does not aim to support)
