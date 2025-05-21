# 🛠 Development Workflow — `{{.ProjectName}}`

This file documents the typical development and release workflow for this Go-based project.

---

## 1. Project Standards

| Area         | Best Practice                                                            |
|--------------|--------------------------------------------------------------------------|
| Structure    | Use `cmd/`, `pkg/` and `configs/` directories                            |
| Module       | Set up `go.mod` early with correct module path                           |
| Versioning   | Use [SemVer](https://semver.org) and git tags (`v1.0.0`, `v1.2.3`, etc.) |
| Licensing    | Include MIT or Apache 2.0 license in the root                            |

---

## 2. Development Flow

### ✨ Feature Workflow

```bash
git checkout -b feat/short-name
# Write feature with tests and docs
# Commit: feat: Add X support to CLI
```

### 🐛 Bugfix Workflow

```bash
git checkout -b fix/short-name
# Write failing test, then fix
# Commit: fix: Prevent panic on malformed config
```

### 🔤 Commit Types

| Prefix   | Meaning                         |
|----------|---------------------------------|
| `feat:`  | New feature                     |
| `fix:`   | Bugfix                          |
| `ref:`   | Refactor (no functional change) |
| `test:`  | Test addition or improvement    |
| `docs:`  | Docs or examples                |
| `chore:` | Meta config, CI, tooling        |

---

## 3. Branching & Tagging

| Action      | Convention            |
|-------------|-----------------------|
| Features    | `feat_*`              |
| Hotfixes    | `fix_*` or `hotfix_*` |
| Releases    | Tag with `vX.Y.Z`     |
| Main branch | Always stable         |

---

## 4. Release Checklist

1. Finalize and merge changes
2. Bump version manually or via script

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. Update `CHANGELOG.md`
4. Create a GitHub release with highlights

---

## 5. Changelog Format (Example)

```
## v1.2.0 - 2025-05-20
### Added
- `--unix` mode for socket-only apps
- Metrics export via `--metrics`

### Fixed
- Crash on empty config
- Incorrect CLI error on an invalid flag

### Changed
- Default config path moved to `~/.config/toolname`
```

---

## 6. Planning & Backlog

Use [`ROADMAP.md`](./ROADMAP.md) or GitHub Projects for planning.

Example roadmap block:

```
## Upcoming Features
- [ ] Plugin system
- [ ] Config validator
- [ ] Sponsor-only templates
```

---

## 7. CLI Tool Release Checklist

* [ ] `--help` flag is descriptive
* [ ] Reasonable defaults work without configuration
* [ ] Flags: `--version`, `--config` supported
* [ ] Edge case tests added
* [ ] README includes install/run instructions
