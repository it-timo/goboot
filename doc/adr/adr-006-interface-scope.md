# 📄 ADR-006: Strict Interface Scope for Config Modules

**Tags:** `interfaces`, `validation`, `modularity`

---

## Status

✅ Accepted

---

## Context

In Go, interfaces should describe behavior needed by the caller — not be used for early abstraction.
The config system initially used a narrow `ServiceConfig` interface to ensure every config:

- Has a stable identifier (`ID`)
- Can load itself from YAML (`ReadConfig`)
- Can validate itself before registration (`Validate`)

These methods are the minimum required to treat configs uniformly within the manager.

---

## Decision

Use the following as the only required interface for config modules:

```go
type ServiceConfig interface {
    ID() string
    ReadConfig(confPath string) error
    Validate() error
}
```

All other logic remains in the concrete type.

---

## Advantages

- Minimally invasive interface: avoids bloated abstractions
- Encourages concrete-first thinking and discoverable logic
- All behavior remains transparent and inspectable

---

## Disadvantages

- Cannot invoke arbitrary behavior on configs unless cast to the concrete type
