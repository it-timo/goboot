# ADR-005: Thin CLI Architecture

**Tags:** `architecture`, `cli`, `separation-of-concerns`

---

## Status

Accepted

---

## Context

CLI entry points often become a second application layer with business logic and global state.
That increases coupling and makes CLI behavior harder to test deterministically.

---

## Decision

Keep `cmd/goboot` as a thin entry point.

- `main` handles argument parsing, startup wiring, and process exit behavior.
- Service logic remains in `pkg/`.
- Use local `flag.FlagSet` instances instead of global `flag.CommandLine` state.

---

## Advantages

- CLI logic is easier to test with injected arguments.
- Service behavior remains reusable outside the CLI entry point.
- Ownership boundary between wiring and business logic stays clear.

---

## Disadvantages

- Slightly more explicit wiring code in `cmd/`.
- New contributors must follow layering constraints when adding features.

---

## Alternatives Considered

- **Business logic directly in `main`:** rejected due to low testability and mixed responsibilities.
- **Global flag variables:** rejected due to hidden shared state and reduced test isolation.
