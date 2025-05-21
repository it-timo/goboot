# 📁 Project Structure — `{{.ProjectName}}`

This document describes the directory and file structure used in `{{.ProjectName}}`, 
intended as a consistent, high-standard Go project layout. 
All components listed here are **expected to exist and be functional** when released.

---

## 📂 Root Directory (Top-Level)

| File / Dir                 | Purpose                                                           |
|----------------------------|-------------------------------------------------------------------|
| `.editorconfig`            | Enforces consistent editor formatting settings                    |
| `.env.example`, `.env.ci`  | Environment configuration templates for local and CI environments |
| `.gitattributes`           | Git settings for line endings and diff behavior                   |
| `.gitignore`               | Files and directories to be excluded from version control         |
| `.nvmrc`, `.tool-versions` | Toolchain version hints (e.g., Go, Node, Python)                  |
| `go.mod`, `go.sum`         | Go module configuration                                           |
| `LICENSE`, `NOTICE`        | Licensing and legal compliance files                              |
| `README.md`                | Overview and quickstart guide                                     |
| `ROADMAP.md`               | (Optional) Planned features and layer milestones                  |
| `VERSIONING.md`            | Describes the versioning and tagging strategy                     |
| `WORKFLOW.md`              | Development, CI/CD, and contribution workflow description         |

---

## 📂 Key Directories

### `cmd/`

Contains CLI entry points.

* `{{.ProjectName}}`: The root executable command.
* Subdirectories may define other command-line tools or utilities as needed.

### `pkg/`

Reusable Go packages structured for modularity.

* `config/`: Configuration loading and validation logic.
* `utils/`: Generic helper functions used across the codebase.
* `{{.ProjectName}}/`: Core project logic, structured by domain.

---

## 📌 Structure Principles

This layout follows these principles:

* **Clarity**: All files and folders serve a clear, active role.
* **Modularity**: Separation by domain and concern.
* **No placeholders**: All items are meaningful and ready to be used.
* **Best Practices First**: Includes tooling and metadata expected in real-world projects.
