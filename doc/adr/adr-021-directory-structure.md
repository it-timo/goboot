# 📄 ADR-021: Service and Directory Naming Conventions

**Tags:** `filesystem`, `naming`, `oss-guidelines`

---

## Status

✅ Accepted

---

## Context

As `goboot` is designed to be modular and extensible, the need for consistent naming of internal
vs. external services becomes critical for:

- separation of concerns,
- clear ownership,
- and easier community collaboration.

The core `goboot` modules currently use the `base*` prefix (e.g. `baseproject`, `baselocal`, `baselint`)
to reflect foundational responsibilities.

However, third-party service modules are expected to follow a different convention to distinguish contributions
and prevent naming collisions.

---

## Decision

- **Internal services (first-party)** must be named using the prefix `base*`.
  - Example: `baseproject`, `baselint`, `baselocal`

- **External or user-contributed services** must use a **provider-scoped prefix**, such as:
  - `ghuser_linter`
  - `yourname_quality`
  - `corpteam_ci`

- **Package structure** must mirror the service ID:
  - A service ID of `baseproject` maps to `pkg/baseproject/`
  - A service ID of `johns_ci` maps to `pkg/johns_ci/`

This ensures both naming uniqueness and accountability within the OSS ecosystem.

---

## Advantages

- Prevents service name collisions
- Provides clear ownership (internal vs. external)
- Enables future service registries or discovery without ambiguity
- Makes it easy to spot trusted, core `goboot` logic at a glance

---

## Disadvantages

- Slightly more verbose for external contributors
- Requires communication of naming guidelines

---

## Alternatives Considered

- Single flat namespace for all services — **rejected** due to collision risk and ownership ambiguity
