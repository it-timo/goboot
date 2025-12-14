# 📄 ADR-030: Template Suffix `.tmpl` to Isolate Lint/Test Pipelines

**Tags:** `templates`, `linting`, `testing`, `tooling`, `scaffolding`

---

## Status

✅ Accepted

---

## Context

Template sources now include runnable Go files, scripts, and config snippets for linting/testing services.
Running repo-level `go test` or lint commands against these raw templates caused false positives
(invalid imports, placeholder paths, or mixed styles).
We needed a deterministic way to keep template logic in-tree without polluting the project’s own lint/test results.

---

## Decision

- All scaffold inputs stored under `templates/*` use a dedicated filename suffix `TemplateSuffix = ".tmpl"`.
- Generation services (`baseProject`, `baseLint`, `baseLocal`, `baseTest`) **strip `.tmpl` only after rendering** paths
and contents inside the secure `os.Root`, keeping templates invisible to host tooling.
- Lint/test tools ignore `.tmpl` files by default (wrong language/suffix),
so CI naturally focuses on goboot’s source, not template payloads.
- Template files must remain syntactically valid **after rendering**
but may contain placeholders that would otherwise fail lint/test when unrendered.

---

## Advantages

- Eliminates false lint/test failures from placeholder code.
- Keeps template assets co-located with source without special CI filters per file type.
- Simplifies service logic: a single suffix convention signals “render then strip”.

---

## Disadvantages

- Editors/viewers may not apply language tooling to `.tmpl` files automatically.
- Contributors must remember to strip the suffix when adding new render steps.

---

## Alternatives Considered

- **Dedicated `templates/` exclusion lists per linter/test:** Fragile across tools
and would hide accidental non-template files.
- **Embedding templates in Go code:** Reduces readability and increases binary size; harder to iterate on template content.
