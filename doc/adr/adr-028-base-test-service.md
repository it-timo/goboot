# 📄 ADR-028: Dedicated Test Scaffolding via `baseTest` Service

**Tags:** `service`, `testing`, `templates`, `scripts`, `scaffolding`

---

## Status

✅ Accepted

---

## Context

Earlier versions scaffolded projects with linting and local scripts but shipped **no test skeleton**.
Contributors had to manually wire suites, imports, and commands, leading to inconsistent quality and slow onboarding.
As testing became a first-class focus, we needed a modular service that could render project-aware test suites
and surface a runnable test command without coupling to other services.

---

## Decision

- Introduce a `baseTest` service (ID `base_test`) with a typed `BaseTestConfig` (fields: `sourcePath`, `useStyle`, `testCmd`,
derived project name/import variants).
- Run `baseTest` as a normal service after `base_project` and before `base_local`; it uses `gobootutils.CreateRootDir`
and `os.Root` to constrain all writes to the generated project root.
- Apply a **two-pass render**: walk `sourcePath` to render paths (templated directory/file names, remove `.tmpl`,
skip Ginkgo suite files when `useStyle` ≠ `ginkgo`), then walk the secure root to render file contents with `text/template`.
- Guard against accidental overwrite by comparing template and target paths (`ComparePaths`).
- When a `Registrar` is available (e.g., `base_local`), register the resolved `testCmd` into make/task targets
and `scripts/test.sh` via `RegisterLines`/`RegisterFile`.
- House templates under `templates/test_base/` to keep logic data-driven and allow future style variants.

---

## Advantages

- New projects start with runnable tests and imports already pointed at the module path.
- Keeps testing concerns isolated from `base_project` while still integrating with script generation when enabled.
- Template-driven approach eases evolution (additional styles, extra packages) without code changes.

---

## Disadvantages

- Adds another service to orchestrate and configure.
- Default Ginkgo templates pull in external dependencies unless `useStyle` is set to `go`.

---

## Alternatives Considered

- **Bundle tests into `base_project`:** Would bloat the base scaffold and blur responsibilities.
- **Static file copy without templating:** Would fail to inject import paths/project names and limit reuse across projects.
- **Hard-wire scripts inside `baseTest`:** Rejected to keep script rendering delegated to `base_local`.
