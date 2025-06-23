# 📄 ADR-020: Service Execution Strategy with Config Matching

**Tags:** `execution`, `config-matching`, `service-manager`

---

## Status

✅ Accepted

---

### Context

Each registered service must only execute if its configuration has been loaded and validated.
This ensures services are never run in an unconfigured or invalid state.

The `serviceManager.runAll()` method performs this check by:

- Looking up the config via `cfgMgr.Get(id)`
- Skipping the service if no config is found
- Running the service with its config otherwise

---

## Decision

Implement config-aware execution:

```go
cfg, ok := sm.cfgMgr.Get(id)
if !ok {
    fmt.Printf("Service %q skipped (no configuration loaded)\n", id)
    continue
}
err := svc.Run(cfg)
```

---

## Advantages

- Prevents accidental service execution
- Ensures 1:1 match between service logic and config data
- Skippable by design — no hard failure for missing optional services

---

## Disadvantages

- All services require a matching config (even if minimal)
- Skipped services may lead to partial outputs if not monitored closely
