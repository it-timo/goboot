# ADR-030: Template Suffix `.tmpl` to Isolate Lint/Test Pipelines

**Tags:** `templates`, `linting`, `testing`, `tooling`, `scaffolding`

---

## Status

Accepted

---

## Context

Unrendered template files can be misinterpreted by repo-level lint/test tooling.
A consistent marker is needed to separate scaffold input from source files.

---

## Decision

Use `.tmpl` as the template suffix for scaffold inputs.
Services strip `.tmpl` only during generation output paths.

---

## Advantages

- Reduces false positives from unrendered template files.
- Keeps template assets in-repo with explicit identification.
- Aligns rendering behavior across services.

---

## Disadvantages

- Editors may provide weaker language tooling for suffixed files.
- Contributors must follow suffix conventions consistently.

---

## Alternatives Considered

- **Tool-by-tool ignore lists for template directories:** rejected due to maintenance drift.
- **Embed templates into Go source:** rejected due to poorer readability and iteration speed.
