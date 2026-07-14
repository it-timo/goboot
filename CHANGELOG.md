# Changelog

All notable changes to `goboot` are tracked here.

The project follows semantic versioning during the pre-1.0 rollout described in
[`doc/VERSIONING.md`](./doc/VERSIONING.md).

## Unreleased

### v0.5.0 — Supply Chain Security

Added:

- `base_supplychain` service with validated, pinned scanner versions and an
  explicit dependency-license allowlist.
- Generated GitHub CodeQL, vulnerability, license, and CycloneDX SBOM jobs.
- Generated GitLab vulnerability, license, and CycloneDX SBOM jobs.
- `SUPPLY_CHAIN.md` policy output for generated repositories.
- Supply-chain security guide and ADR-039.

Changed:

- The default goboot configuration enables supply-chain security generation.
- GitLab pipelines include a security stage before release automation.
- The root and generated-project Go baseline is 1.26.5, which includes the fix
  for GO-2026-4970.

### v0.4.0 — Governance & Contribution

Added:

- `base_governance` service with validated maintainer and default-branch inputs.
- Shared `CODEOWNERS`, `CONTRIBUTING.md`, and `SECURITY.md` outputs.
- GitHub issue forms and pull request templates.
- GitLab issue and merge request templates.
- Profile-aware governance baselines for minimal, standard, enterprise, and OSS projects.
- Governance guide and ADR-038.

Changed:

- The default goboot configuration enables repository governance generation.
- Contributor-facing security guidance directs vulnerability reports to private
  provider channels instead of public issues.

### v0.3.0 — Template Profiles

Added:

- Validated root `profile` selection with `minimal`, `standard`, `enterprise`,
  and `oss` values.
- Profile-specific Go lint baselines and cyclomatic-complexity thresholds.
- Profile-default test styles and commands with explicit service overrides.
- Generated `PROFILE.md` documenting the selected baseline.
- Profile guide and ADR-037.

Changed:

- Omitted profiles default to `standard` for backward compatibility.
- Generated project README and structure documentation identify the active
  profile.

### v0.2.1 — Release Automation

Added:

- `base_release` service with validated GoReleaser settings.
- Tag-driven GitHub and GitLab release jobs for generated projects.
- Cross-platform AMD64/ARM64 archives and checksum manifests.
- Automated changelog generation from Git history.
- GoReleaser configuration and release workflow for `goboot` itself.
- Release guide and ADR-036.

Changed:

- `base_ci` can aggregate a release job in addition to build, test, lint, and
  container jobs.
- Version selection is explicitly derived from immutable semantic-version tags;
  goboot does not guess or commit version changes.

### v0.2.0 — Containerization

Added:

- `base_docker` service for generated CLI-style Go project container packaging.
- `configs/base_docker.yml` with explicit Docker template and build settings.
- Multi-stage `Dockerfile` template with non-root runtime stage.
- Minimal `docker-compose.yml` template with empty ports by default.
- `.dockerignore` template.
- Generated local Docker workflows:
  - `make docker-build`
  - `make compose-up`
  - `make container-check`
  - matching Taskfile tasks
  - `scripts/docker.sh`
- GitLab and GitHub CI container jobs that validate compose config and build the
  generated image.
- Root `Dockerfile` and `.dockerignore` for building the `goboot` CLI image.
- `make docker_build`, `make docker_smoke`, `task docker:build`, and
  `task docker:smoke` for local generator image validation.
- GitHub Actions container workflow for the `goboot` CLI image.
- ADR-035 documenting the CLI-containerization scope and trade-offs.
- Containerization documentation and examples.

Changed:

- Default `configs/goboot.yml` now enables `base_docker`.
- Generated-project verification now includes Docker outputs and helper scripts.
- Release checks now include root `goboot` image smoke testing and generated
  project container build/compose validation.

Notes:

- The generated project remains CLI-style. `base_docker` does not assume a
  network service or default exposed port.
