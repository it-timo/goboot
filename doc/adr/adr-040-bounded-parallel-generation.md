# ADR-040: Bounded Parallel Generation with Measured Baselines

**Tags:** `performance`, `concurrency`, `benchmarks`, `determinism`

---

## Status

Accepted — 2026-07-14

## Context

Regular generation services are independent after `base_project` creates the
project root, but goboot previously ran every service serially. Larger template
sets need reproducible measurements before optimizations can be evaluated.
Unbounded goroutines would make resource use unpredictable, while concurrent
registrations into the shared CI and local maps would introduce data races.

## Decision

- Add a root `parallelism` setting with a valid range of `1` through `32`.
- Default omitted values to `1` for backward-compatible serial execution.
- Keep the prior and subsequent execution phases serial and ordered.
- Sort configured regular service IDs before dispatch.
- Run regular services through a bounded worker pool when parallelism exceeds 1.
- Wait for the complete regular-service batch and report errors in sorted service
  order.
- Protect the CI and local registries with mutexes and copy input slices at the
  registration boundary.
- Maintain benchmarks for template execution, atomic file rendering, large
  template sets, large root configs, and service orchestration.
- Record pull-request benchmark output as an artifact without imposing a noisy
  hosted-runner timing threshold.

## Consequences

- Users can trade resource use for faster generation with an explicit setting.
- Existing configurations retain their former serial behavior.
- Aggregate CI and local outputs are generated only after all producers finish.
- Concurrent service completion order may vary, but selected failures and
  aggregate rendering remain deterministic.
- The manager completes already scheduled work before returning a regular-service
  error; there is no partial cancellation contract.
- Performance changes have repeatable local entry points and reviewable CI data.

## Alternatives considered

### Always run serially

This is simplest but leaves independent service work underutilized and does not
meet the scale milestone.

### Start one goroutine per service

This reduces orchestration code but provides no resource bound as the service
registry grows.

### Enable parallel execution by default

This could improve default runtime but would change scheduling for every existing
configuration. An explicit opt-in is safer before a stable release.

### Gate pull requests on a fixed benchmark threshold

Hosted-runner variance makes a fixed timing threshold prone to false failures.
Recorded artifacts support review while stable benchmark infrastructure can be
introduced later.
