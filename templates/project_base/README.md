# {{.ProjectName}}

> ...

{{- if eq .GitProvider "github" }}
[![License](https://img.shields.io/github/license/{{.GitHubUser}}/{{.LowerProjectName}})](LICENSE)
[![Version](https://img.shields.io/github/v/release/{{.GitHubUser}}/{{.LowerProjectName}}?include_prereleases)](https://github.com/{{.GitHubUser}}/{{.LowerProjectName}}/releases)
{{- else if eq .GitProvider "gitlab" }}
[![License](https://img.shields.io/gitlab/license/{{.GitLabUser}}/{{.LowerProjectName}})](LICENSE)
[![Version](https://img.shields.io/gitlab/v/release/{{.GitLabUser}}/{{.LowerProjectName}}?include_prereleases)](https://gitlab.com/{{.GitLabUser}}/{{.LowerProjectName}}/-/releases)
{{- end }}

---

## 📦 Overview

`{{.ProjectName}}` ...

---

## 📁 Project Structure

This repository provides:

- A CLI entry point (`/cmd/{{.LowerProjectName}}/main.go`)
- Modular packages (`/pkg/...`) for domain and utility code
- Core documentation (`README.md`, `ROADMAP.md`, `VERSIONING.md`, etc.)

For a full breakdown, see [`PROJECT_STRUCTURE.md`](./PROJECT_STRUCTURE.md) and [`ROADMAP.md`](./ROADMAP.md).

---

## 🚀 Usage

### Requirements

- Go {{.UsedGoVersion}}+

### Run

```bash
go run ./cmd/{{.ProjectName}}
```

---

## 🧭 Project Planning

This repository includes the following planning documents:

* [ROADMAP.md](./ROADMAP.md) — Feature planning and milestone tracking
* [VERSIONING.md](./VERSIONING.md) — Semantic versioning and release rules
* [WORKFLOW.md](./WORKFLOW.md) — Contributor and CI process

---

## 📦 Releases

{{- if eq .GitProvider "github" }}
All releases are tagged and listed under the [Releases](https://github.com/{{.GitHubUser}}/{{.LowerProjectName}}/releases) page.
{{- else if eq .GitProvider "gitlab" }}
All releases are tagged and listed under the [Releases](https://gitlab.com/{{.GitLabUser}}/{{.LowerProjectName}}/-/releases) page.
{{- end }}

To see current goals, refer to [`ROADMAP.md`](./ROADMAP.md).

---

## ⚖ License

Licensed under the MIT License. See [LICENSE](./LICENSE).
Additional notices may be found in [NOTICE](./NOTICE) if applicable.

---

## 🙌 Acknowledgements

This project was scaffolded using [goboot](https://github.com/it-timo/goboot), 
a CLI tool for generating reproducible Go project structures.
