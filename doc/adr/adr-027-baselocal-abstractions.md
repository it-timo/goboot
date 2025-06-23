# 📄 ADR-027: Script Type Abstractions and FileList Control (`baseLocal`)

**Tags:** `baseLocal`, `scripts`, `filelist`, `conditional-rendering`, `modularity`, `extensibility`

---

## Status

✅ Accepted

---

## Context

The `baseLocal` service supports multiple output types:

- `Makefile`
- `Taskfile.yml`
- `.pre-commit-config.yaml`
- Scripts in `scripts/`

However, not all projects want or need all of these. For example, a user might prefer `make` but not use `task`,
or might skip `pre-commit` entirely. We want to:

- Avoid rendering unnecessary files
- Give users control over which files get scaffolded
- Allow clean expansion in the future (e.g., `scripts/windows`, `fishfile`)

This is managed through a `FileList` field in the `BaseLocalConfig`.

---

## Decision

The `baseLocal` service interprets a config-provided `FileList` that defines **which script types** should be rendered.

```yaml
fileList:
    - make
    - task
    - script
    - commit
```

Internally, this:

- Filters what is rendered in `copyFiles()`
- Prevents unnecessary output or empty files

The available options are aligned with constants in `types/`, e.g. `ScriptNameMake`, `ScriptNameTask`.

---

## Advantages

- Users have fine-grained control over rendered outputs
- Prevents noise in projects that don’t need all formats
- Enables clean support for future formats (e.g., `justfile`, `scripts/windows`)
- Makes the `baseLocal` service more declarative and modular

---

## Disadvantages

- Slight increase in complexity in `copyFiles()` and registration logic
- Inconsistent `Register*()` calls from services may silently be ignored if a file type is disabled

---

## Alternatives Considered

- **Always render all script formats**: Violates minimalism, bloat project with unused files
- **Make rendering conditional only inside logic**: Less user-friendly, harder to trace at config level
