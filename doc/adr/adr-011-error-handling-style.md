# 📄 ADR-011: Error Handling Style: Explicit Early Returns

**Tags:** `errors`, `style`, `robustness`

---

## Status

✅ Accepted

---

## Context

Go allows inline error handling using `if err := ...; err != nil`, but this can obscure variable lifetime
and hurt readability in larger, modular functions.

Given `goboot`’s architectural aim of readability, clarity, and testability,
we enforce a consistent error style throughout the codebase.

---

## Decision

- Always declare errors using `err = ...` and follow with `if err != nil { return ... }`
- Avoid `panic()` unless truly unrecoverable (e.g., internal invariant breach)
- No silent fallback behavior — every error path must be handled or explicitly ignored with a rationale

---

## Advantages

- Uniform readability across files and services
- Easier debugging and logging
- Predictable control flow, especially in services and orchestrators

---

## Disadvantages

- More verbose
- Requires code review discipline

---

## Alternatives Considered

- Inline error handling (`if err := ...`) — rejected for consistency and traceability
- `panic` for expected paths — rejected due to robustness and OSS expectations
