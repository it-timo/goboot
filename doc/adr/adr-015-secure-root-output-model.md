# 📄 ADR-015: Scoped Filesystem Output via `os.Root`

**Tags:** `filesystem`, `security`, `sandboxing`, `go-1.23+`

---

## Status

✅ Accepted

---

## Context

Filesystem generation often risks unintentional overwrites, escapes (`../../`), or host-level side effects.
Go 1.23 introduced `*os.Root`, a new API that securely confines I/O operations to a specific root path,
preventing traversal or unsafe writes.

---

## Decision

Use `os.OpenRoot()` to establish a **safe, bounded output directory** for each scaffolding run.
All rendering operations occur inside this confined space — no raw `os.Create`, `filepath.Walk`,
or unbounded path logic is allowed.

---

## Advantages

- Strong guarantees against accidental or malicious writings outside the project dir
- Aligns with modern Go sandboxing and testability goals
- Easier to mock or simulate in test environments

---

## Disadvantages

- Slightly more verbose than traditional file APIs
- Requires contributors to understand the `*os.Root` model
- Cannot reuse legacy `os.*` and `filepath.*` code without adapting
