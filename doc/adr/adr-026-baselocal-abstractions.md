# ADR-026: Script Type Abstractions and FileList Control (`baseLocal`)

**Tags:** `baseLocal`, `scripts`, `filelist`, `conditional-rendering`, `modularity`, `extensibility`

---

## Status

Accepted

---

## Context

Projects differ in preferred local tooling (for example, make vs task usage).
Rendering all formats by default can add unnecessary files.

---

## Decision

Use `fileList` in `BaseLocalConfig` to control which script/output types are rendered.
`baseLocal` filters output generation based on that list.

---

## Advantages

- Users can limit generated local tooling to what they use.
- Reduces unused-file noise in generated repositories.
- Provides a clear config-level contract for local output selection.

---

## Disadvantages

- Rendering path is more conditional and requires broader test coverage.
- Registered entries may be ignored when corresponding file types are disabled.

---

## Alternatives Considered

- **Always render all local outputs:** rejected due to unnecessary output for many projects.
- **Conditional rendering hidden only in code with no config field:** rejected
  because behavior would be less explicit to users.
