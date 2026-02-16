# ADR-009: Strict Interface Scope for Config Modules

**Tags:** `interfaces`, `validation`, `modularity`

---

## Status

Accepted

---

## Context

The config manager needs a minimal common contract to treat service configs uniformly.
Overly broad interfaces would hide concrete behavior and increase abstraction cost.

---

## Decision

Keep `ServiceConfig` minimal and focused on required lifecycle operations:

```go
type ServiceConfig interface {
    ID() string
    ReadConfig(confPath string) error
    Validate() error
}
```

Other behavior stays on concrete config types.

---

## Advantages

- Small interface surface with clear purpose.
- Concrete behavior remains visible.
- Lower risk of interface bloat.

---

## Disadvantages

- Advanced config-specific operations require concrete type access.
- Some callers may need explicit type assertions.

---

## Alternatives Considered

- **Larger shared config interface:** rejected because many methods would be unused by the manager.
- **No shared interface:** rejected because manager logic would require type switches per config module.
