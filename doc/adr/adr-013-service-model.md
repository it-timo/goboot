## 📄 ADR-013: Modular Service Execution Model

**Tags:** `services`, `modularity`, `run-logic`

## Status
✅ Accepted

### Context

To support a modular CLI tool that can execute discrete logic like project generation, linting setup, or Docker config, 
`goboot` must support isolated and composable services. 
Each service needs to execute logic conditionally based on validated configuration.

Dynamic plugin loading or runtime registration is *not* aligned with the project’s goals of traceability, 
OSS quality, and Go idioms.

### Decision

Introduce a `Service` interface:

```go
type Service interface {
	ID() string
	Run(cfg config.ServiceConfig) error
}
```

Each service is:

* Instantiated explicitly
* Registered via a hardcoded switch in `RegisterServices()`
* Executed only if a matching config exists

Service wiring is static and controlled in `goboot.GoBoot`.

### Advantages

* Full separation between service logic and configuration
* Zero runtime hooks or dynamic dispatch
* Clear control flow, easy to debug and extend
* Follows Go best practices (no magic registration)

### Disadvantages

* Requires manual mapping of every new service
* No support for dynamic service loading or plugin discovery (by design)
