# ADR-025: `baseLocal` - Service Purpose and Script Boundaries

**Tags:** `service`, `scripts`, `responsibility`, `modularity`, `execution-boundaries`

---

## Status

Accepted

---

## Context

Generated local developer tooling files are shared outputs across services.
A single owner is needed to avoid fragmented write responsibilities.

---

## Decision

`baseLocal` owns developer-facing script/file rendering.

Owned outputs include:

- `Makefile`
- `Taskfile.yml`
- `.pre-commit-config.yaml`
- optional `scripts/` files

Other services register intent through shared interfaces; they do not write these files directly.

---

## Advantages

- Centralized rendering rules for local tooling files.
- Better consistency for generated developer workflows.
- Other services remain focused on domain configuration.

---

## Disadvantages

- `baseLocal` becomes a coordination point that must remain stable.
- Contributor mistakes in registration names can affect output composition.

---

## Alternatives Considered

- **Per-service direct script writes:** rejected due to duplication and conflict risk.
- **No script generation support:** rejected because it shifts repeated setup burden to users.
