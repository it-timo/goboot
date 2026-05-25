# Containerization

`goboot` has two containerization paths:

- a root container image for the `goboot` generator CLI,
- generated Docker output for scaffolded Go projects through `base_docker`.

## `goboot` CLI Image

The repository root contains a `Dockerfile` for building `goboot` itself:

```bash
make docker_build
```

Equivalent direct command:

```bash
docker build --build-arg VERSION="$(cat .version)" -t goboot:local .
```

`make docker_build` tags the image as `goboot:<contents-of-.version>`.
The direct command above uses `goboot:local` only as a manual example.

Smoke test:

```bash
make docker_smoke
```

The runtime image intentionally uses a Go image, not distroless, and the default
Go image tag tracks the patch version required by `go.mod`. `goboot` currently
runs `go mod tidy` after generation by default when the generated project
contains a `go.mod`, so the container needs the Go toolchain unless the caller
passes `--skip-go-mod-tidy`.

Example usage with a mounted repository workspace:

```bash
docker run --rm \
  -v "$PWD":/workdir \
  -w /workdir \
  "goboot:$(cat .version)" \
  --config ./configs/goboot.yml
```

Generated output paths are resolved inside the mounted `/workdir`.

## Generated Project Containerization

`base_docker` adds container packaging for generated Go projects.

The generated application shape is currently CLI-style. The Docker output is
therefore intended to build and run the CLI binary in a container, not to model a
network service with default ports, health checks, or service discovery.

## Config

Enable the service in `goboot.yml`:

```yaml
services:
  - id: "base_docker"
    confPath: "./configs/base_docker.yml"
    enabled: true
```

Minimal `base_docker.yml`:

```yaml
sourcePath: "./templates/docker_base"
goVersion: "1.26.3"
runtimeImage: "gcr.io/distroless/static-debian12:nonroot"
fileList:
  - dockerfile
  - compose
  - dockerignore
ports: []
```

Optional fields:

- `binaryName`: defaults to the lowercase project name.
- `mainPackage`: defaults to `./cmd/<lowercase-project-name>`.
- `ports`: optional compose mappings such as `"8080:8080"` for projects that
  later become long-running services.

## Generated Files

When all outputs are enabled, generated projects receive:

- `Dockerfile`
- `docker-compose.yml`
- `.dockerignore`

With `base_local` enabled, generated projects also receive:

- `make docker-build`
- `make compose-up`
- `make container-check`
- matching Taskfile tasks
- `scripts/docker.sh`

With `base_ci` enabled, generated CI receives a container job that validates the
compose file and builds the image.

## Validation Scope

Current validation checks that:

- Docker templates render without unresolved template markers.
- generated Docker files are present when `base_docker` is enabled.
- local lint/test workflows include the generated Docker helper script.
- generated `make container-check` and `task container-check` workflows validate
  compose config, image builds, and the CLI help smoke run during the
  intro-project verification flow.
- GitLab and GitHub provider CI can render container build jobs, and the CI
  canary executes provider container jobs when those workflows exist. GitLab
  container validation uses `gitlab-ci-local --privileged`, and GitHub
  validation can generate an additional `IntroGitHubCanary` project when the
  default generated output targets GitLab.
- generated project verification passes across the supported test/logger/provider
  matrix.

The generated container is not treated as a network-service integration test
target because the scaffolded application is not a network service. The release
gate validates buildability, compose configuration, and CLI help startup instead.
