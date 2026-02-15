# ADR-019: Service Execution Strategy with Config Matching

**Tags:** `execution`, `config-matching`, `service-manager`

---

## Status

Accepted

---

## Context

Registered services should execute only when corresponding validated configuration exists.
Missing config should not crash optional service flows.

---

## Decision

Use config-aware execution in `serviceManager`:

- assign config to services before run (`SetConfig(...)`),
- skip services that have no loaded config, and
- run configured services via `Run()`.

---

## Advantages

- Prevents execution with undefined input.
- Supports optional service enablement.
- Keeps skip behavior explicit in logs.

---

## Disadvantages

- Partial output is possible when users expect disabled services to run.
- Misconfigured service IDs may fail silently if only skip logs are observed.

---

## Alternatives Considered

- **Fail hard when any registered service lacks config:** rejected because optional services are part of the design.
- **Run with implicit zero-value config:** rejected because hidden defaults can mask configuration mistakes.
