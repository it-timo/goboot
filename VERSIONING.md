# Versioning Strategy

The `goboot` project follows [Semantic Versioning 2.0.0](https://semver.org), 
with a **structured feature build-up** during its `v0.x.x` phase.

---

## 📦 Format

`MAJOR.MINOR.PATCH` — e.g., `0.4.0`, `1.0.0`, `1.2.3`

---

## 🔁 Pre-1.0 Logic

Versions below `v1.0.0` **do not imply instability** — instead, they represent a **deliberate, 
layered rollout of features** and foundational tooling.

| Version Range | Purpose                                                |
|---------------|--------------------------------------------------------|
| `v0.0.x`      | Structural setup: config, CI, lint, testing, basics    |
| `v0.1.x`      | Infrastructure: contributors, community, docs          |
| `v0.2.x`      | Packaging: Docker, release systems                     |
| `v0.3.x`      | Contribution: issue templates, reuse, conduct          |
| `v0.4.x`      | Security: dependency checks, code scanning             |
| `v0.5.x`      | Benchmarking support                                   |
| `v1.0.0`+     | Feature-ready CLI, stable bootstrapping system         |

The project intentionally uses `v0.x.x` **to incrementally ship working layers** — 
not as placeholders or “unstable previews.”

---

## 🔖 Release Rules

* All releases are tagged using the `vX.Y.Z` format
* Pre-release identifiers like `-beta.1`, `-rc.1`, or `-dev.2` are allowed when necessary

---

## 🧪 Experimental Versions (Optional)

If needed, pre-release variants may follow this syntax:

* `v0.6.0-beta.1`
* `v1.0.0-rc.1`

…but the main development remains within the progressive `v0.x.x` scope until `v1.0.0` is fully realized.

---

## 🤝 Philosophy

Every `v0.x.x` release is a **public and deliberate step** toward maturity, reflecting meaningful evolution.

---

## 📚 See Also

* [ROADMAP.md](./ROADMAP.md)

