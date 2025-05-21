# 📚 Architecture Decision Records (ADR) — `goboot`

This directory contains formal design decisions made during the development of the `goboot` project.

Each ADR captures the **context**, **tradeoffs**, and **rationale** behind a structural 
or behavioral choice in the system — following the `ADR Template` outlined below.

> 📌 ADRs are ordered to reflect the **logical build-up of architectural foundations**, not by creation date.

---

## 🧱 ADR Template

> Use this format for all new ADRs in the `goboot` project.

```md
# 📄 ADR-XXX: [Title]

**Tags:** `tag1`, `tag2`, ...

## Status
✅ Accepted # or 🟡 Proposed / ❌ Rejected / 🔁 Superseded / 🧊 Deprecated

## Context
[Why this decision was needed — the technical or organizational situation]

## Decision
[What was decided — clearly and actionably]

## Advantages
- ...

## Disadvantages
- ...

## Alternatives Considered
(Optional — only when relevant or worth capturing)
````

---

## 📄 ADR Index

| ADR                                                    | Title                                                   | Tags                                  |
|--------------------------------------------------------|---------------------------------------------------------|---------------------------------------|
| [ADR-001](adr-001-minimalism.md)                       | Minimalism over Generalization                          | design-philosophy, idioms             |
| [ADR-002](adr-002-service-name-registry-in-types.md)   | Centralized Service Name Registry in `types/names.go`   | constants, config-binding             |
| [ADR-003](adr-003-template-engine.md)                  | Template Engine, Structure, and Naming Rules            | templates, text/template, scaffolding |
| [ADR-004](adr-004-config-structure.md)                 | Config System Structure and Philosophy                  | config, modularity, idiomatic-go      |
| [ADR-005](adr-005-no-reflection.md)                    | Explicit Separation of Concerns (No DI, Reflection)     | design-philosophy, maintainability    |
| [ADR-006](adr-006-secure-root-output-model.md)         | Scoped Filesystem Output via `os.Root`                  | filesystem, sandbox, security         |
| [ADR-007](adr-007-interface-scope.md)                  | Strict Interface Scope for Config Modules               | interfaces, modularity                |
| [ADR-008](adr-008-utils-as-pure-functional-toolbox.md) | Intentional Use of `pkg/utils` as Pure Functional Set   | utils, purity, no-side-effects        |
| [ADR-009](adr-009-config-safety.md)                    | Config Manager Behavior and Safety                      | validation, config-manager            |
| [ADR-010](adr-010-create-service-config.md)            | Centralized Config Dispatch via `createServiceConfig()` | config-factory, dispatch              |
| [ADR-011](adr-011-baseproject-scope.md)                | Service Responsibility & Scope – `baseProject`          | service, structure                    |
| [ADR-012](adr-012-rendering.md)                        | Template Rendering Strategy (Paths and Content)         | templates, text/template, generation  |
| [ADR-013](adr-013-service-model.md)                    | Modular Service Execution Model                         | services, architecture                |
| [ADR-014](adr-014-register-services.md)                | Static Service Registration via `RegisterServices()`    | wiring, registration                  |
| [ADR-015](adr-015-execution-matching.md)               | Service Execution Strategy with Config Matching         | config-aware, runtime-logic           |
| [ADR-016](adr-016-error-style.md)                      | Error Handling Style                                    | style, consistency                    |

> 💡 New ADRs must follow the `ADR Template` and be reviewed before merging.

---

## 📌 Status

All decisions marked ✅ have been **accepted and implemented** in the current codebase.
Changes to accepted ADRs require a superseding ADR or a major revision with rationale.
