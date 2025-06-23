# goboot Project Roadmap

This roadmap defines the **planned evolution and goals** of the `goboot` project — a CLI tool
and scaffolding framework that aims to provide industrial-grade, flexible Go project setups.

Currently, the project is **in foundational development**.
Features are introduced incrementally through versioned releases.

---

## 🧱 Project Philosophy

`goboot` is built with:

- Modular, composable templates for scalable Go projects
- Clean project hygiene from the first commit
- Modular, layered feature releases
- Top-tier engineering practices
- Reproducible build and versioning systems
- A clean, extensible design for plugin and template evolution

---

## 🔖 Versioned Milestones

Each milestone incrementally adds functionality and structure. Until `v1.0.0`, breaking changes may still occur.

| Version  | Focus Area                             | Key Additions                                                |
|----------|----------------------------------------|--------------------------------------------------------------|
| `v0.0.0` | Bootstrap Setup                        | CLI entry point, basic config, README, LICENSE               |
| `v0.0.1` | Local Tooling & Linting                | GolangCI config, Makefile, Taskfile, scripts, lint/format    |
| `v0.0.2` | Tests & Utilities                      | `/pkg/utils`, `/test`, bootstrap/test/version scripts        |
| `v0.0.3` | Multi-Command CLI Base                 | `goboot init`, `goboot run`, `goboot validate` structure     |
| `v0.1.0` | CI Integration                         | GitHub/GitLab workflows, YAML lint, test runners             |
| `v0.1.1` | Community & Meta                       | CONTRIBUTORS, SUPPORT, FUNDING, SECURITY, CODEOWNERS         |
| `v0.2.0` | Dockerization                          | Dockerfile, `.dockerignore`, `docker-compose.yml`, test base |
| `v0.2.1` | Release Infrastructure                 | GoReleaser, release-local script, release CI setup           |
| `v0.3.0` | Contribution Templates                 | Issue/PR/MR templates for GitHub & GitLab                    |
| `v0.3.1` | Contributor Documentation & Compliance | CONTRIBUTING.md, REUSE.toml, HACKING.md, 3rd party           |
| `v0.4.0` | Supply Chain & Security Checks         | CodeQL, depcheck, license-checker, `SECURITY_CONTACTS`       |
| `v0.5.0` | Benchmarking Setup                     | Benchmarks for config/logger, BENCHMARKS.md                  |
| `v1.0.0` | CLI Generator: Public Ready            | CHANGELOG, `.version`, template system foundation            |

> 🚧 *Versions prior to `v1.0.0` are unstable and may change structure.*

---

## 🛠️ Technical Goals (Q2 – Q3 2025)

- [x] Tooling layer: Makefile, Taskfile, scripts, linting, formatting
- [ ] Template engine (placeholder parser)
- [ ] Finalize `pkg`, `config` and logger strategies
- [ ] Harden CLI init/check/validate flow

---

## 🔭 Midterm Goals (Q4 2025)

- [ ] Docker image build and publish automation
- [ ] GitHub/GitLab parity enforcement
- [ ] Structural Profiles (OSS vs. Enterprise-style layouts)

---

## 🧬 Long-Term Vision (2026+)

- [ ] Plugin-based CLI extensions
- [ ] Template registry with extension support
- [ ] Declarative init scripting and chaining

---

## 🧬 Long-Term Vision (2027+)

- [ ] Optional CLI frontend for scaffold preview or selection
- [ ] GitHub/GitLab visual sync tool
- [ ] Multi-user/team CI config injection

---

*This document evolves with the project.*
