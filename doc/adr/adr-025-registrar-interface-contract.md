# 📄 ADR-025: Script Integration via `Registrar` Interface

**Tags:** `baselint`, `scripting`, `integration`, `interfaces`, `modularity`

---

## Status

✅ Accepted

---

## Context

Although the `baseLint` service generates and renders linting configuration files (e.g., `.golangci.yml`),
it must **not be responsible** for writing shell scripts, Makefiles, or task definitions directly.
That responsibility belongs to the `baseLocal` service, which handles local developer tooling and script file layout.

However, `baseLint` needs a way to **declare** which linting commands (e.g., `golangci-lint run`, `yamllint .`)
should be executed — so they can be included in Makefiles, Taskfiles, or shell wrappers.
This requires a **decoupled contract** between `baseLint` and `baseLocal`.

---

## Decision

The `types.Registrar` interface was introduced to provide a **loosely coupled scripting contract**.
It allows `baseLint` to declare script entries without needing to know their final format or destination.

```go
type Registrar interface {
    RegisterLines(name string, lines []string) error
    RegisterFile(name string, lines []string) error
}
```

The `baseLint` service optionally receives this interface during initialization.
If present, it registers relevant linter commands (from config) under its service name.

The `baseLocal` service implements this interface and stores all declared lines per script type.
During its own `Run()` execution, it renders the full script files using these declarations.

---

## Advantages

- Clean separation of responsibility: `baseLint` only declares, `baseLocal` renders.
- Future-proof: More services can declare scripts via `Registrar` without central coupling.
- Enables script file generation logic to evolve independently of linting logic.
- Avoid tight wiring or runtime reflection between services.

---

## Disadvantages

- Slight indirection — script files are not rendered at the time of declaration.
- Requires internal registry in `baseLocal` and logic to safely de-duplicate service entries.

---

## Alternatives Considered

- **Direct script rendering in `baseLint`:** Violates separation of concerns and would duplicate logic.
- **Global script registry object:** Harder to test, maintain,
  and reason about than explicit interface-based registration.
