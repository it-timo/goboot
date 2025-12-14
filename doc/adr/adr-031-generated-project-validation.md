# 📄 ADR-031: Validate Generated Projects with Lint & Test Runs

**Tags:** `templates`, `quality`, `ci`, `generated-project`, `linting`, `testing`

---

## Status

✅ Accepted

---

## Context

goboot’s own tests focus on the generator code and sample templates.
However, changes to templates, service wiring, or defaults can silently break the **generated** project
(e.g., invalid imports, missing scripts, stale commands).
We need a contributor habit that verifies the end product, not just the generator.

---

## Decision

- Any change touching templates, service logic, or config defaults requires contributors to:
  1. Generate a full project with the current `configs/goboot.yml` (or equivalent test fixture).
  2. Run the generated project’s lint and test commands (defaults: `make lint`, `make test` or their Taskfile equivalents).
- Validation must be performed on the rendered output (post-`.tmpl` stripping)
to ensure the scaffolds remain runnable and standards-compliant.
- Failures discovered in the generated project must block merges; fixes belong alongside the originating change.

---

## Advantages

- Catches template drift and broken script/command registrations early.
- Ensures defaults remain runnable for users without extra setup.
- Aligns generator evolution with real-world consumption, not just internal unit tests.

---

## Disadvantages

- Adds runtime to contributor workflows (generate + lint/test).
- Requires local tooling or containerized runners matching the generated commands.

---

## Alternatives Considered

- **Rely only on generator unit tests:** Misses integration gaps in rendered projects.
- **CI-only generation check:** Helpful but slower feedback; developers still need a local gate.
