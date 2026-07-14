# Performance and Scale

`goboot` measures the work that dominates larger generation runs and supports
bounded concurrency for independent services.

## Service parallelism

Set `parallelism` in the root configuration:

```yaml
parallelism: 4
```

The supported range is `1` through `32`. When the field is omitted, goboot uses
`1` to preserve the serial behavior of existing configurations.

The execution phases remain explicit:

1. `base_project` initializes the project root.
2. Configured independent services run in a bounded worker pool.
3. `base_ci` and `base_local` render their aggregated registrations.

Only the middle phase runs concurrently. The manager waits for every scheduled
service before starting aggregate output. If multiple services fail, goboot
returns the error for the lexicographically first service ID so the reported
failure does not depend on goroutine scheduling.

The CI and local command registries serialize concurrent writes and copy
registered command slices. Services therefore cannot mutate shared aggregate
state after registration.

## Benchmarks

Run the complete benchmark suite locally:

```bash
make benchmark
# or
task benchmark
```

Override the sample duration when comparing changes:

```bash
make benchmark BENCH_TIME=3s
```

The suite covers:

| Benchmark | Scale contract |
| --------- | -------------- |
| `BenchmarkExecuteTemplateText` | Small and large in-memory templates |
| `BenchmarkRenderTemplateToFile` | Parse, render, and atomic file replacement |
| `BenchmarkLargeProjectRendering` | 500 templates per generated-project operation |
| `BenchmarkGoBootConfigParsing` | Root configuration with 500 service entries |
| `BenchmarkRegularServiceExecution` | 64 services at parallelism 1, 4, and 8 |

Pull requests run the same suite in `.github/workflows/performance.yml` and
retain `benchmark-results.txt` for 14 days. Hosted-runner measurements are
recorded as evidence rather than used as a hard pass/fail threshold because
shared runners do not provide stable enough timing for a trustworthy gate.

## CPU and memory profiles

Capture profiles for the large-project rendering benchmark:

```bash
make profile_generation BENCH_TIME=3s
# or
task profile:generation
```

This writes:

- `profiles/generation-cpu.out`
- `profiles/generation-memory.out`

Inspect them with the standard Go profiler:

```bash
go tool pprof profiles/generation-cpu.out
go tool pprof profiles/generation-memory.out
```

The `profiles/` directory is ignored by Git. Commit benchmark code and documented
conclusions, not machine-specific profile binaries.

## Comparing results

Use the same Go version, host, `GOMAXPROCS`, and benchmark duration for both
revisions. Compare `ns/op`, `B/op`, and `allocs/op`; do not infer a regression
from a single hosted-runner sample. Functional tests and the race detector remain
the correctness gates for the parallel execution path.
