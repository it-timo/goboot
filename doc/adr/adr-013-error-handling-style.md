# ADR-013: Error Handling Style: Explicit Early Returns

**Tags:** `errors`, `style`, `robustness`

---

## Status

Accepted

---

## Context

The codebase favors explicit control flow to support reviewability and debugging.
Mixed error styles make it harder to reason about lifecycle-heavy service logic.

---

## Decision

Use explicit early returns for errors as the default style.

- Prefer straightforward `err` assignment and `if err != nil` checks.
- Avoid panic-based control flow for expected runtime paths.
- Do not silently suppress errors without explicit rationale.

---

## Advantages

- Consistent error flow across services.
- Easier stack-level tracing during failures.
- Lower ambiguity in control-flow branches.

---

## Disadvantages

- Higher verbosity in some functions.
- Requires consistent review enforcement to avoid style drift.

---

## Alternatives Considered

- **Mixed style with frequent inline `if err := ...` blocks:** rejected as default to keep repository-wide consistency.
- **Panic for recoverable paths:** rejected because generation failures should surface as typed errors.
