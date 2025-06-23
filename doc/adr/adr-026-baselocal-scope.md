# 📄 ADR-026: `baseLocal` — Service Purpose and Script Boundaries

**Tags:** `service`, `scripts`, `responsibility`, `modularity`, `execution-boundaries`

---

## Status

✅ Accepted

---

## Context

In the `goboot` system, many services (like `baseLint`, etc.) define tooling
or features that may require **supporting scripts** — such as CLI wrappers, Make targets, or Taskfile steps.
However, these services should **not** write or manage those scripts themselves.

There must be a **single, authoritative service** that:

- Owns the structure and rules for `Makefile`, `Taskfile.yml`, and `scripts/`
- Accepts declarations from other services
- Applies rendering logic and enforces consistent layout

This prevents script logic from being fragmented across multiple services, and ensures all developer-side tooling
(e.g., `make lint`, `scripts/lint.sh`) is handled uniformly.

---

## Decision

The `baseLocal` service is responsible for **all developer-facing script integration**.

It defines the script boundaries for:

- `Makefile`
- `Taskfile.yml`
- `.pre-commit-config.yaml`
- `scripts/` directory (optional helper scripts)

All script-related files are registered centrally in `baseLocal`,
and other services use the `Registrar` interface to declare what commands they want included.

No other service writes script files directly.

---

## Advantages

- Single source of truth for all script-related files.
- Enables full user control over scripts via the `baseLocal` config.
- Simplifies downstream services: they only declare intent, not rendering.
- Allows future script types (e.g., OS-specific, PowerShell) without changing upstream services.

---

## Disadvantages

- Slight coupling — services must know they can optionally register via `Registrar`.
- Requires `baseLocal` to track internal maps of registered lines (per file type and service name).

---

## Alternatives Considered

- **Each service writes its own script lines/files:** Leads to duplication, formatting inconsistencies,
  and race conditions if multiple services write to the same file.
- **No script generation at all:** Would require every user to wire `make lint` or `task test` manually —
  a poor onboarding experience.
