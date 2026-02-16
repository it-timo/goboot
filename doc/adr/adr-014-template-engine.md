# ADR-014: Template Engine, Structure, and Naming Rules

**Tags:** `templates`, `text/template`, `scaffolding`, `structure`

---

## Status

Accepted

---

## Context

Generation requires rendering both paths and file contents from project config.
The template system should remain predictable and easy to test.

---

## Decision

Use Go `text/template` for path and content rendering.

- Render template-derived paths and copy structure.
- Render file contents in a second pass.
- Keep templates under dedicated template directories.
- Use config fields as explicit template input.

---

## Advantages

- Uses standard library tooling.
- One template model for path and content rendering.
- Fits existing testing approach for deterministic output.

---

## Disadvantages

- Template expressions remain less ergonomic than feature-rich engines.
- Missing template fields fail at runtime unless covered by tests.

---

## Alternatives Considered

- **Suffix-driven mixed rendering behavior by file extension only:** rejected
  because one uniform rendering contract is simpler.
- **Third-party function-rich engines (for example, Sprig-heavy stacks):** rejected to limit dependency and complexity overhead.
