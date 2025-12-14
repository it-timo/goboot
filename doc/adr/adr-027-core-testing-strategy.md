# 📄 ADR-027: Core Testing Strategy & Coverage Baseline

**Tags:** `testing`, `bdd`, `coverage`, `filesystem`, `safety`

---

## Status

✅ Accepted

---

## Context

`goboot` moved from a largely untested scaffold to a generator that manipulates file systems,
template rendering, and service orchestration.
Without a disciplined test suite we risk silent regressions when changing path handling (`os.Root`),
service registration, or template rendering.
We also need deterministic coverage targets to keep contributors aligned as new services (lint/local/test) land.

---

## Decision

- Use **Ginkgo v2 + Gomega** for every package (including `cmd/`), keeping `_suite_test.go` bootstrap files per package.
- Target **80%+ coverage** per critical package and ~90% overall; enforce via the default `DefaultGoTestCMD`
(`go test -race -timeout=5m -coverprofile=coverage.txt`).
- Exercise **real filesystem flows** with `os.MkdirTemp` and `os.Root` instead of mocks to validate secure-root behavior,
template rendering, and path comparison.
- Prefer **table-driven specs via `DescribeTable`** for permutations and explicit `BeforeEach`/`AfterEach` for isolation.
- Keep **constants and registries under test** (`goboottypes`, `serviceManager`, config loaders) to guard contract drift.

---

## Advantages

- High confidence when refactoring templating, path safety, or service wiring.
- Regressions surface with readable BDD output instead of ad-hoc logging.
- Security-sensitive logic (`os.Root`, path normalization) is continuously exercised.

---

## Disadvantages

- Ginkgo/Gomega dependency in the dev toolchain.
- File-heavy specs increase test runtime compared to pure mocks.

---

## Alternatives Considered

- **Minimal smoke tests only:** Too weak for the generator’s filesystem-heavy surface area.
- **Mock-heavy unit tests:** Would miss integration issues around `os.Root`, template rendering, and script registration.
