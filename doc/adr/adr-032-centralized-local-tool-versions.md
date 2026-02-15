# ADR-032: Centralized Local Tool Versions via `versions.env`

**Tags:** `tooling`, `linting`, `dev-experience`, `versions`

---

## Status

Accepted

---

## Context

Version values used in local scripts and task runners can drift when duplicated.
A shared source is required to keep local workflows aligned.

---

## Decision

Use `versions.env` as the single source of local tool version values and load it from local entry points.
Local values remain tag-based for readability and maintenance ease.

---

## Advantages

- One update point for local tool version changes.
- Consistent versions across Make/Task/scripts/pre-commit.
- Improved readability for contributors.

---

## Disadvantages

- Requires discipline to keep all local entry points sourcing the same file.
- Tag-based references have weaker immutability than digest pins.

---

## Alternatives Considered

- **Digest pins for all local tools:** rejected due to readability and frequent update overhead.
- **Duplicate versions per file:** rejected because drift risk is high.
