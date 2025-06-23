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

---

## Status

✅ Accepted # or 🟡 Proposed / ❌ Rejected / 🔁 Superseded / 🧊 Deprecated

---

## Context

[Why this decision was needed — the technical or organizational situation]

---

## Decision

[What was decided — clearly and actionably]

---

## Advantages

- ...

---

## Disadvantages

- ...

---

## Alternatives Considered
(Optional — only when relevant or worth capturing)
````

---

## 📄 ADR Index

| ADR                                                    | Title                                                         | Tags                                                                                          |
|--------------------------------------------------------|---------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| [ADR-001](adr-001-minimalism.md)                       | Minimalism over Generalization                                | design-philosophy, minimalism, abstraction                                                    |
| [ADR-002](adr-002-no-reflection.md)                    | Explicit Separation of Concerns — No Runtime Reflection or DI | philosophy, idioms, anti-patterns                                                             |
| [ADR-003](adr-003-utils-as-pure-functional-toolbox.md) | Intentional Use of `pkg/utils` as Pure Functional Set         | utils, hygiene, modularity, stateless                                                         |
| [ADR-004](adr-004-service-name-registry-in-types.md)   | Centralized Service Name Registry in `types/names.go`         | constants, service-names, structure, decoupling                                               |
| [ADR-005](adr-005-config-structure.md)                 | Config System Structure and Philosophy                        | config, modular-design, idiomatic-go                                                          |
| [ADR-006](adr-006-interface-scope.md)                  | Strict Interface Scope for Config Modules                     | interfaces, validation, modularity                                                            |
| [ADR-007](adr-007-config-safety.md)                    | Config Manager Behavior and Safety                            | manager, validation, static-analysis                                                          |
| [ADR-008](adr-008-create-service-config.md)            | Centralized Config Dispatch via `createServiceConfig()`       | dispatch, registration, no-reflection                                                         |
| [ADR-009](adr-009-config-validation.md)                | Typed Config Structure and Validation Strategy                | config, validation, typed-structure                                                           |
| [ADR-010](adr-010-error-style.md)                      | Error Handling Style                                          | errors, style, readability                                                                    |
| [ADR-011](adr-011-error-handling-style.md)             | Error Handling Style: Explicit Early Returns                  | errors, style, robustness                                                                     |
| [ADR-012](adr-012-template-engine.md)                  | Template Engine, Structure, and Naming Rules                  | templates, text/template, scaffolding, structure                                              |
| [ADR-013](adr-013-rendering.md)                        | Template Rendering Strategy (Paths and Content)               | templates, rendering, text/template                                                           |
| [ADR-014](adr-014-template-rendering.md)               | Template Rendering Strategy using `text/template`             | templates, text/template, scaffolding                                                         |
| [ADR-015](adr-015-secure-root-output-model.md)         | Scoped Filesystem Output via `os.Root`                        | filesystem, security, sandboxing, go-1.23+                                                    |
| [ADR-016](adr-016-baseproject-scope.md)                | Service Responsibility & Scope – `baseProject`                | service, responsibility, structure                                                            |
| [ADR-017](adr-017-service-model.md)                    | Modular Service Execution Model                               | services, modularity, run-logic                                                               |
| [ADR-018](adr-018-service-registration.md)             | Static Service Registration and Orchestration                 | services, registration, explicit-architecture                                                 |
| [ADR-019](adr-019-register-services.md)                | Static Service Registration via `RegisterServices()`          | registration, explicit-control, top-level-mapping                                             |
| [ADR-020](adr-020-execution-matching.md)               | Service Execution Strategy with Config Matching               | execution, config-matching, service-manager                                                   |
| [ADR-021](adr-021-directory-structure.md)              | Service and Directory Naming Conventions                      | filesystem, naming, oss-guidelines                                                            |
| [ADR-022](adr-022-future-extensions.md)                | Extensibility Strategy for New Services and Features          | extensibility, oss, architecture                                                              |
| [ADR-023](adr-023-baselint-scope.md)                   | Dedicated Linting via `baseLint` Service                      | service, linting, quality, separation-of-concerns                                             |
| [ADR-024](adr-024-base-linter-strategy.md)             | Linter Configuration Rendering Strategy                       | baselint, linting, templates, rendering, golangci                                             |
| [ADR-025](adr-025-registrar-interface-contract.md)     | Script Integration via `Registrar` Interface                  | baselint, scripting, integration, interfaces, modularity                                      |
| [ADR-026](adr-026-baselocal-scope.md)                  | `baseLocal` — Service Purpose and Script Boundaries           | service, scripts, responsibility, modularity, execution-boundaries                            |
| [ADR-027](adr-027-baselocal-abstractions.md)           | Script Type Abstractions and FileList Control for `baseLocal` | baseLocal, scripts, filelist, conditional-rendering, modularity, extensibility                |
| [ADR-028](adr-028-registrar-interface.md)              | Decoupled Script Coordination via `Registrar` Interface       | baseLocal, scripts, interface, coordination, registrar, extensibility, separation-of-concerns |

> 💡 New ADRs must follow the `ADR Template` and be reviewed before merging.

---

## 📌 Status

All decisions marked ✅ have been **accepted and implemented** in the current codebase.
Changes to accepted ADRs require a superseding ADR or a major revision with rationale.
