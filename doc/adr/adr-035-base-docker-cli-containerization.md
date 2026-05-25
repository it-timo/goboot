# ADR-035: CLI Containerization for `goboot` and Generated Projects

**Tags:** `docker`, `containerization`, `service`, `templates`, `ci`, `cli`

---

## Status

Accepted

---

## Context

`goboot` needs a containerization layer as part of the `v0.2.0` milestone.
That layer has two related but separate concerns:

- generated Go projects should receive container packaging,
- the `goboot` generator CLI itself should be runnable from a container.

Generated projects are currently CLI-style applications, not network services.
The containerization contract therefore needs to package and validate the CLI
binary without implying HTTP ports, long-running server behavior, health checks,
or service discovery.

The `goboot` CLI also runs `go mod tidy` after generation by default when the
generated project contains a `go.mod`. A containerized `goboot` runtime therefore
needs access to the Go toolchain unless the user passes `--skip-go-mod-tidy`.

The existing architecture already provides:

- explicit service registration,
- typed service configs,
- template-owned output files,
- registrar-based integration with local scripts and CI.

Containerization should follow those patterns instead of being folded into
`base_project`, `base_local`, or `base_ci`.

---

## Decision

- Add a dedicated `base_docker` service.
- Keep Docker templates under `templates/docker_base`.
- Generate only container-owned files:
  - `Dockerfile`
  - `docker-compose.yml`
  - `.dockerignore`
- Use a multi-stage Dockerfile:
  - Go build stage
  - small non-root runtime stage
- Treat `docker-compose.yml` as local CLI-container execution scaffolding, not a
  network-service deployment manifest.
- Keep port mappings optional and empty by default.
- Register Docker commands with `base_local`:
  - container image build
  - compose up
  - compose/build validation
- Register container validation commands with `base_ci` so generated GitLab and
  GitHub CI can build the image.
- Keep the service explicit in `configs/goboot.yml` and configurable through
  `configs/base_docker.yml`.
- Add a root `Dockerfile` for the `goboot` CLI.
- Use a Go-based runtime image for the `goboot` container so default generation
  can still run `go mod tidy`.
- Add a root `.dockerignore`.
- Add local `docker_build` and `docker_smoke` targets for the `goboot` image.
- Add a GitHub Actions container workflow that builds and smoke-tests the
  `goboot` image.

---

## Advantages

- Container behavior is isolated behind a clear service boundary.
- CLI-style generated projects get useful packaging without pretending to be
  network services.
- Existing local and CI registrars remain the coordination mechanism.
- Templates remain deterministic and reviewable.
- Users can disable containerization by disabling `base_docker`.
- The generator itself can run from a container in environments where a local
  `goboot` binary is not installed.
- Keeping Go in the `goboot` runtime image preserves current CLI defaults.

---

## Disadvantages

- The generated compose file is intentionally minimal.
- Runtime integration testing is limited to compose config and image build
  validation until generated projects expose stable behavior worth asserting.
- Provider CI templates need one more job file per provider.
- The `goboot` container image is larger than a distroless image because it
  includes the Go toolchain.

---

## Alternatives Considered

- **Add Docker files to `base_project`:** rejected because containerization is an
  optional deployment concern, not core project structure.
- **Only generate a Dockerfile:** rejected because local parity also benefits
  from compose validation and local command registration.
- **Assume an HTTP service and expose a default port:** rejected because the
  generated project is CLI-style today.
- **Run generated containers as integration tests:** deferred until generated
  projects have stable, meaningful runtime behavior beyond command startup.
- **Use a distroless runtime for the `goboot` image:** rejected because the
  current CLI runs `go mod tidy` by default.
