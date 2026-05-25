# Changelog

All notable changes to `goboot` are tracked here.

The project follows semantic versioning during the pre-1.0 rollout described in
[`doc/VERSIONING.md`](./doc/VERSIONING.md).

## Unreleased

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
