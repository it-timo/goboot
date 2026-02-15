# ADR-017: Modular Service Execution Model

**Tags:** `services`, `modularity`, `run-logic`

---

## Status

Accepted

---

## Context

`goboot` runs multiple services that may or may not be enabled in one run.
Execution needs a common service contract with explicit orchestration.

---

## Decision

Adopt a service interface with explicit lifecycle methods:

```go
type Service interface {
    ID() string
    SetConfig(cfg config.ServiceConfig) error
    Run() error
}
```

Services are registered explicitly and executed only when matching config exists.
Wiring remains static in the orchestrator.

---

## Advantages

- Common lifecycle contract across services.
- Explicit assignment of typed config before execution.
- Execution ordering is controlled in one place.

---

## Disadvantages

- Each new service requires manual registration and ordering decisions.
- No runtime discovery for out-of-tree services.

---

## Alternatives Considered

- **Single monolithic generator without service boundaries:** rejected due to weak modularity and harder test isolation.
- **Runtime plugin/discovery model:** rejected due to added operational and debugging complexity.
