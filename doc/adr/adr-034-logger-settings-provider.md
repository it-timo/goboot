# ADR-034: Logger Settings Provider

**Tags:** `logger`, `scaffolding`, `templates`, `service-boundaries`, `ownership`

---

## Status

Accepted

---

## Context

Generated projects can optionally include runtime logging.
Logger support affects files owned by the base project scaffold, including:

- `cmd/<project>/main.go`
- `pkg/<project>/<project>.go`
- `go.mod`

A dedicated logger service that writes those files directly creates unstable ownership.
It can overwrite base project output, depend on service execution order, and conflict with
test scaffolding that imports or calls generated application APIs.

This is the same class of problem solved for local script generation by ADR-024:
contributors should declare intent through typed contracts, while the file-owning service
renders the final artifact.

---

## Decision

`base_logger` provides validated logger settings only.

`base_project` owns runtime project files and renders logger-aware variants from its own
templates using those settings.

The orchestration layer wires logger settings through explicit interfaces:

- `LoggerSettingsProvider` exposes validated logger settings.
- `LoggerSettingsReceiver` accepts logger settings before rendering.

The absence of `base_logger` means logger wiring is disabled and `base_project` renders the
plain baseline scaffold.

---

## Advantages

- Preserves strict file ownership.
- Avoids services overwriting files they do not own.
- Keeps generated runtime code coherent across `main.go`, package code, and `go.mod`.
- Supports multiple logger implementations without changing service execution order.
- Keeps extension points explicit and testable.

---

## Disadvantages

- `base_project` templates contain conditional logger branches.
- Adding a logger implementation requires updating the owning project templates.
- Logger settings are currently a narrow contract and may need extension for format,
  output, sampling, or environment-specific defaults.

---

## Alternatives Considered

- **Dedicated logger templates that overwrite project files:** rejected because this
  violates service ownership and creates order-dependent behavior.
- **Hardwire one logger into base project:** rejected because logger choice is a scaffold
  option and the project already supports multiple logger types.
- **Create separate whole-project template trees per logger:** rejected for now because it
  would duplicate most base project files and increase drift risk.
