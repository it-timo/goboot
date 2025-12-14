# 📄 ADR-006: Strict Service Isolation

**Tags:** `architecture`, `modularity`, `dependencies`

---

## Status

✅ Accepted

---

## Context

In a modular system like `goboot`, services (e.g., `baseproject`, `baselint`) can easily become tightly coupled
if they import each other to share logic.
This leads to "spaghetti code," circular dependencies, and monolithic testing requirements
where changing one service breaks unrelated ones.

---

## Decision

We enforce **Strict Service Isolation**:

### 1. No Direct Service-to-Service Imports

- Services (packages in `pkg/<service_name>`) **must not** import other service packages.
- Example: `pkg/baselint` cannot import `pkg/baseproject`.

### 2. Communication via Shared Contracts

- Interaction between services occurs ONLY through:
    1. **Shared Config**: `pkg/config` (defines data structures).
    2. **Shared Types/Interfaces**: `pkg/goboottypes` (defines constants and interfaces like `Registrar`).
    3. **Orchestrator**: `pkg/goboot` (wires them together).

### 3. Dependency Injection

- Do not instantiate a dependency inside a service.
- If Service A needs to register scripts with Service B (e.g., `base_lint` -> `base_local`),
it must use the `Registrar` interface injected by the Orchestrator.

---

## Advantages

- **Decoupling**: Services can be developed, tested, and enabled/disabled independently.
- **Build Speed**: Changes in one leaf package do not invalidate the build cache of others.
- **Maintainability**: Clear boundaries prevent "god classes" from emerging across packages.

---

## Disadvantages

- **Indirection**: requires defining interfaces in `pkg/goboottypes` rather than just calling a public method.
