# 📄 ADR-019: Static Service Registration via `RegisterServices()`

**Tags:** `registration`, `explicit-control`, `top-level-mapping`

---

## Status

✅ Accepted

---

## Context

Many systems use maps, reflection, or dependency injection to register services.
This makes debugging harder and increases a cognitive load. `goboot` avoid all runtime registration patterns.

Instead, service registration is centralized in `GoBoot.RegisterServices()`, where each known service is:

- Checked against its ID
- Instantiated explicitly
- Registered via `serviceManager.register(...)`

Unknown IDs result in an error, ensuring only supported features are activated.

---

## Decision

Use a fixed switch-case dispatch to instantiate services based on ID:

```go
switch meta.ID {
case types.ServiceNameBaseProject:
    err = gb.ServiceMgr.register(baseProject.NewBaseProject(...))
default:
    return fmt.Errorf("unknown service ID: %s", meta.ID)
}
```

---

## Advantages

- Predictable and safe: no unexpected logic paths
- Enables static analysis and exhaustive testing
- Clear responsibility: registration happens only once, in one place

---

## Disadvantages

- Slightly more verbose
- Requires code changes to support new services
