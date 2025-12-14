# 📄 ADR-005: Thin CLI Architecture

**Tags:** `architecture`, `cli`, `separation-of-concerns`

---

## Status

✅ Accepted

---

## Context

Command-line interface (CLI) entry points often accumulate business logic, making them hard to test and maintain.
Global state (like `flag.StringVar` variables at the package level) causes side effects
that make parallel testing impossible and hampers reusability.

---

## Decision

We enforce a **Thin CLI** architecture:

### 1. No Business Logic in `main`

- The `main` package is strictly an **entry point**.
- It allows ONLY:
    1. Argument parsing.
    2. Configuration loading.
    3. Service registration invocation.
    4. Exit code handling.

### 2. No Global Flag State

- **Prohibited**: `flag.Parse()` on the global `flag.CommandLine`.
- **Required**: Use `fs := flag.NewFlagSet(...)` and `fs.Parse(args)`.
- This allows `run(args []string)` to be tested with arbitrary inputs in parallel without race conditions.

### 3. Logic Delegation

- All actual work must happen in `pkg/`.
- `cmd/goboot` delegates immediately to `goboot.NewGoBoot` and its methods.

---

## Advantages

- **Testability**: The CLI flow can be end-to-end tested (`main_test.go`) without spawning subprocesses.
- **Reusability**: `pkg/goboot` can be embedded in other tools if needed.
- **Clarity**: It is immediately obvious where "wiring" ends and "logic" begins.

---

## Disadvantages

- **Boilerplate**: Passing arguments explicitly is slightly more verbose than using globals.
