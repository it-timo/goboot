# goboot

> A modular, versioned starting point for building production-grade Go projects.

[![License](https://img.shields.io/github/license/it-timo/goboot)](LICENSE)
[![Version](https://img.shields.io/github/v/release/it-timo/goboot?include_prereleases)](https://github.com/it-timo/goboot/releases)

---

## 📦 What is `goboot`?

> `goboot` is a deterministic scaffolding tool that provides a clean, modular foundation for real-world Go projects 
> — including applications, tools, and infrastructure code.
>
> It’s **not** a framework. It’s **not** an IDE.
> Instead, `goboot` is built for OSS maintainers, contributors, 
> and enterprise teams who care about long-term project hygiene, structure, and reproducibility.
>
> This repository focuses on **layered, progressive structure scaffolding**, not one-click demos or opinionated code generation.
>
> The goal isn’t just to get you started — it’s to help you grow Go projects that remain clean, consistent, 
> and scalable over time.

---

## 📁 Current Scope (v0.0.0)

This initial version (`v0.0.0`) establishes:

- A minimal CLI entry point (`cmd/goboot/main.go`)
- Basic configuration loading (`pkg/config`)
- Early repository scaffolding (`README.md`, `LICENSE`, `ROADMAP.md`)
- A clearly versioned development path under [ROADMAP.md](./ROADMAP.md)

> 👉 This project is in **early planning and layout phase**. Expect placeholder logic and evolving structure.

---

## 📐 Intended Design Principles

Even in early stages, `goboot` is being built with:

- Layered versioning and changelog visibility
- Clear module boundaries (`cmd/`, `pkg/`, `configs/`, etc.)
- Future support for Docker, CI/CD, and template-driven code generation

You can follow the structural milestones in [`ROADMAP.md`](./ROADMAP.md).

---

## 🛠️ Getting Started (For Contributors Only)

### Requirements

- Go 1.24+
- Make (for later local scripts)

### Clone and Explore

```bash
git clone https://github.com/it-timo/goboot.git
cd goboot
go run ./cmd/goboot
```

There’s not much output yet — but that’s intentional.

---

## 📚 Project Planning

This repository uses:

* [ROADMAP.md](./ROADMAP.md) for planned milestones
* [VERSIONING.md](./VERSIONING.md) for semantic version handling
* [WORKFLOW.md](./WORKFLOW.md) to define long-term contribution and CI logic
* [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) to track how the folder layout evolves over time
* [doc/adr/](./doc/adr) for architecture decisions in ADR format
* [doc/img/](./doc/img) for flow visualizations


These documents evolve alongside the project.

---

## ⚖️ License

Licensed under the MIT License. See [LICENSE](./LICENSE).
Includes attribution in [NOTICE](./NOTICE) (if applicable).

---

## 🚧 Status

> `goboot` is in **pre-alpha**.
> Intended for structural exploration and reproducible project setup. Not yet suitable for generating production-ready projects.

---

## 💖 Support This Project

If `goboot` helps you or saves you time, consider supporting its development:

- [💖 GitHub Sponsors](https://github.com/sponsors/it-timo)
- [🎁 Ko-Fi](https://ko-fi.com/ittimo)
- [☕ BuyMeACoffee](https://buymeacoffee.com/ittimo)

> No pressure — just a small way to say "thanks" if it brought you value.
