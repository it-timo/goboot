# Versioning Strategy

`goboot` follows [Semantic Versioning 2.0.0](https://semver.org),
with structured feature rollout during `v0.x.x`.

---

## 📦 Format

`MAJOR.MINOR.PATCH` — e.g., `0.4.0`, `1.0.0`, `1.2.3`

---

## 🔁 Pre-1.0 Logic

Versions below `v1.0.0` represent a deliberate,
layered rollout of features, tooling, and project hygiene.

| Version Range | Purpose                                                        |
|---------------|----------------------------------------------------------------|
| `v0.0.x`      | Structural core: services, config, linting, testing, FS safety |
| `v0.1.x`      | CI/CD & infrastructure: pipelines, release flow, templates     |
| `v0.2.x`      | Packaging: Docker images, release artifacts                    |
| `v0.3.x`      | Template profiles and profile-aware scaffolding                |
| `v0.4.x`      | Governance: contribution workflow, CODEOWNERS, security docs   |
| `v0.5.x`      | Supply chain & security checks                                 |
| `v0.6.x`      | Benchmarking & performance                                     |
| `v1.0.0`+     | Stable public CLI with backward compatibility guarantees       |

The project intentionally uses `v0.x.x` **to incrementally ship working layers** —
not as placeholders or “unstable previews.”

---

## 🔖 Release Rules

- All releases are tagged using the `vX.Y.Z` format
- Pre-release identifiers like `-beta.1`, `-rc.1`, or `-dev.2` are allowed when necessary

---

## 🧪 Experimental Versions (Optional)

If needed, pre-release variants may follow this syntax:

- `v0.6.0-beta.1`
- `v1.0.0-rc.1`

…but the main development remains within the progressive `v0.x.x` scope until `v1.0.0` is fully realized.

---

## 🤝 Philosophy

Every `v0.x.x` release is a **public and deliberate step** toward maturity, reflecting meaningful evolution.

---

## 📚 See Also

- [ROADMAP.md](../ROADMAP.md)
