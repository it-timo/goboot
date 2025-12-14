# 📄 ADR-024: Decoupled Script Coordination via `Registrar` Interface

**Tags:** `baseLocal`, `scripts`, `interface`, `coordination`, `registrar`, `extensibility`, `separation-of-concerns`

---

## Status

✅ Accepted

---

## Context

Multiple services in `goboot` (e.g. `baseLint`) may contribute shell script commands that are intended
to appear in centralized developer-facing files (e.g., `Makefile`, `Taskfile.yml`, or `scripts/lint.sh`).

Rather than requiring each service to:

- Know about `baseLocal`’s internal structure
- Directly manipulate script file outputs

…a generic registration mechanism is introduced via the `types.Registrar` interface:

```go
type Registrar interface {
    RegisterLines(name string, lines []string) error
    RegisterFile(name string, lines []string) error
}
```

This keeps services focused on **domain logic**, while `baseLocal` acts as a **script orchestrator**.

---

## Decision

The following design is adopted:

- The `baseLocal` service implements the `Registrar` interface.
- Each service that contributes script lines (like `baseLint`) receives a `SetRegistrar()` call during bootstrapping.
- During `Run()`, the contributing service calls `RegisterLines(...)` and/or `RegisterFile(...)` with its commands.
- `baseLocal` uses these entries during its own `copyFiles()` to render templates like:
  - `Makefile` → aggregated `make` lines per service
  - `Taskfile.yml` → collected `task` entries
  - `scripts/` → rendered shell script files (e.g., `lint.sh`, `test.sh`)

> ⚠️ **Lifecycle Note**: The registration hook `SetScriptReceiver` is called during the service bootstrap phase in `pkg/goboot/service.go`.

- This strict ordering ensures that `baseLocal` is ready to receive commands before other services are initialized.
- This preserves modularity and avoids hard coupling between services.

---

## Advantages

- Clean separation of responsibilities: Services don’t render scripts themselves
- Centralized control in `baseLocal` ensures consistent formatting
- New services can integrate easily by just calling `RegisterLines(...)`
- Supports different output formats (make, task, scripts) transparently

---

## Disadvantages

- Slight learning curve — contributors must understand the registrar pattern
- Registration order matters if line conflicts occur (rare in scoped services)

---

## Alternatives Considered

- **Direct file writing by each service**: Violates separation-of-concerns,
  leads to duplication and format inconsistency
- **Global file mutation helper**: Harder to validate, error-prone without struct-based guarantees
