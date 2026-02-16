# Architecture Decision Records (ADR) - `goboot`

This directory contains architecture decisions for `goboot`.

Each ADR captures:

- context (why the decision was needed),
- decision (what was chosen),
- consequences (advantages and disadvantages), and
- alternatives considered.

ADRs are ordered by logical architecture progression, not by creation date.

---

## ADR Template

Use this structure for all new ADRs.

```md
# ADR-XXX: [Title]

**Tags:** `tag1`, `tag2`, ...

---

## Status

Accepted # or Proposed / Rejected / Superseded / Deprecated

---

## Context

[Technical and organizational background]

---

## Decision

[Concrete decision statement]

---

## Advantages

- ...

---

## Disadvantages

- ...

---

## Alternatives Considered

- [Alternative 1 and why it was not selected]
- [Alternative 2 and why it was not selected]
```

---

## ADR Index

| ADR                                                    | Title                                                         | Tags                                                                           |
|--------------------------------------------------------|---------------------------------------------------------------|--------------------------------------------------------------------------------|
| [ADR-001](adr-001-minimalism.md)                       | Minimalism over Generalization                                | design-philosophy, minimalism, abstraction                                     |
| [ADR-002](adr-002-no-reflection.md)                    | Explicit Separation of Concerns - No Runtime Reflection or DI | philosophy, idioms, anti-patterns                                              |
| [ADR-003](adr-003-utils-as-pure-functional-toolbox.md) | Intentional Use of `pkg/gobootutils` as Pure Functional Set   | utils, hygiene, modularity, stateless                                          |
| [ADR-004](adr-004-bdd-testing-standards.md)            | BDD Testing Standards and Tooling                             | testing, bdd, standards, quality                                               |
| [ADR-005](adr-005-thin-cli-architecture.md)            | Thin CLI Architecture                                         | architecture, cli, separation-of-concerns                                      |
| [ADR-006](adr-006-strict-service-isolation.md)         | Strict Service Isolation                                      | architecture, modularity, dependencies                                         |
| [ADR-007](adr-007-service-name-registry-in-types.md)   | Centralized Name Registry in `goboottypes/names.go`           | constants, service-names, structure, decoupling                                |
| [ADR-008](adr-008-config-structure.md)                 | Config System Structure and Philosophy                        | config, modular-design, idiomatic-go                                           |
| [ADR-009](adr-009-interface-scope.md)                  | Strict Interface Scope for Config Modules                     | interfaces, validation, modularity                                             |
| [ADR-010](adr-010-config-safety.md)                    | Config Manager Behavior and Safety                            | manager, validation, static-analysis                                           |
| [ADR-011](adr-011-create-service-config.md)            | Centralized Config Dispatch via `createServiceConfig()`       | dispatch, registration, no-reflection                                          |
| [ADR-012](adr-012-config-validation.md)                | Typed Config Structure and Validation Strategy                | config, validation, typed-structure                                            |
| [ADR-013](adr-013-error-handling-style.md)             | Error Handling Style: Explicit Early Returns                  | errors, style, robustness                                                      |
| [ADR-014](adr-014-template-engine.md)                  | Template Engine & Rendering Strategy                          | templates, text/template, scaffolding, structure                               |
| [ADR-015](adr-015-secure-root-output-model.md)         | Scoped Filesystem Output via `os.Root`                        | filesystem, security, sandboxing, go-1.23+                                     |
| [ADR-016](adr-016-baseproject-scope.md)                | Service Responsibility & Scope - `baseProject`                | service, responsibility, structure                                             |
| [ADR-017](adr-017-service-model.md)                    | Modular Service Execution Model                               | services, modularity, run-logic                                                |
| [ADR-018](adr-018-service-registration.md)             | Static Service Registration and Orchestration                 | services, registration, explicit-architecture                                  |
| [ADR-019](adr-019-execution-matching.md)               | Service Execution Strategy with Config Matching               | execution, config-matching, service-manager                                    |
| [ADR-020](adr-020-directory-structure.md)              | Service and Directory Naming Conventions                      | filesystem, naming, oss-guidelines                                             |
| [ADR-021](adr-021-future-extensions.md)                | Extensibility Strategy for New Services and Features          | extensibility, oss, architecture                                               |
| [ADR-022](adr-022-baselint-scope.md)                   | Dedicated Linting via `baseLint` Service                      | service, linting, quality, separation-of-concerns                              |
| [ADR-023](adr-023-base-linter-strategy.md)             | Linter Configuration Rendering Strategy                       | baselint, linting, templates, rendering, golangci                              |
| [ADR-024](adr-024-registrar-interface.md)              | Registrar Interface & Script Lifecycle                        | baselint, scripting, integration, interfaces, modularity                       |
| [ADR-025](adr-025-baselocal-scope.md)                  | `baseLocal` - Service Purpose and Script Boundaries           | service, scripts, responsibility, modularity, execution-boundaries             |
| [ADR-026](adr-026-baselocal-abstractions.md)           | Script Type Abstractions and FileList Control for `baseLocal` | baseLocal, scripts, filelist, conditional-rendering, modularity, extensibility |
| [ADR-027](adr-027-core-testing-strategy.md)            | Core Testing Strategy & Coverage Baseline                     | testing, bdd, coverage, filesystem, safety                                     |
| [ADR-028](adr-028-base-test-service.md)                | Dedicated Test Scaffolding via `baseTest` Service             | service, testing, templates, scripts, scaffolding                              |
| [ADR-029](adr-029-test-template-styles.md)             | Test Template Styles - Ginkgo by Default, Stdlib as Opt-Out   | testing, templates, ginkgo, stdlib, flexibility                                |
| [ADR-030](adr-030-template-suffix-policy.md)           | Template Suffix `.tmpl` to Isolate Lint/Test Pipelines        | templates, linting, testing, tooling, scaffolding                              |
| [ADR-031](adr-031-generated-project-validation.md)     | Validate Generated Projects with Lint & Test Runs             | templates, quality, ci, generated-project, linting, testing                    |
| [ADR-032](adr-032-centralized-local-tool-versions.md)  | Centralized Local Tool Versions via `versions.env`            | tooling, linting, dev-experience, versions                                     |
| [ADR-033](adr-033-ci-policy-and-provider-layout.md)    | CI Policy Modes and Provider-Scoped Template Layout           | ci, templates, config, security, provider-layout                               |

---

## Status Policy

Decisions marked `Accepted` are treated as current architecture.
Changes to accepted ADRs should be made by either:

- creating a superseding ADR, or
- revising the ADR with explicit rationale in the same change.
