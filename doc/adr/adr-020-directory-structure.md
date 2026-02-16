# ADR-020: Service and Directory Naming Conventions

**Tags:** `filesystem`, `naming`, `oss-guidelines`

---

## Status

Accepted

---

## Context

Service IDs map to config sections, runtime registration, and package paths.
Inconsistent naming increases collision risk and lowers discoverability.

---

## Decision

Define naming conventions for internal and external services.

- Internal services use the `base*` prefix.
- External services use scoped prefixes (for example, organization/user prefixes).
- Package directory names mirror service IDs.

---

## Advantages

- Lower chance of service ID collisions.
- Clear distinction between first-party and external modules.
- Predictable package layout from service ID.

---

## Disadvantages

- External IDs become longer.
- Convention compliance requires documentation and review checks.

---

## Alternatives Considered

- **Flat global namespace:** rejected because ownership and collision handling are weaker.
- **Completely free-form naming:** rejected because mapping and tooling become less predictable.
