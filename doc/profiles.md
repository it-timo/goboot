# Template Profiles

Template profiles provide named lint, test, governance, and documentation baselines without
changing goboot's explicit service composition model.

Select a profile in the root configuration:

```yaml
profile: standard
```

Supported values are `minimal`, `standard`, `enterprise`, and `oss`. If the
field is omitted, goboot uses `standard`.

## Baselines

| Profile      | Default tests    | Test command                                     | Governance                         |
| ------------ | ---------------- | ------------------------------------------------ | ---------------------------------- |
| `minimal`    | Standard library | Fast tests without race or coverage              | Bug workflow                       |
| `standard`   | Ginkgo/Gomega    | Race detection and coverage                      | Bug and feature workflows          |
| `enterprise` | Ginkgo/Gomega    | Race detection, shuffled order, coverage         | Standard plus controlled changes   |
| `oss`        | Ginkgo/Gomega    | Race detection and atomic coverage               | Standard plus documentation intake |

Go lint baselines remain a small complexity-20 set for `minimal`, a balanced
complexity-10 set for `standard` and `oss`, and a stricter complexity-8 set for
`enterprise`.

## Precedence

Profiles fill defaults and provide template baselines. Explicit service values
remain authoritative. For example, setting `useStyle: go` or a custom
`testCmd` in `base_test.yml` overrides the selected profile default.

Profiles do not enable or disable services. The root `services` list remains the
only source of service composition, preserving deterministic and auditable runs.

Generated projects contain `PROFILE.md` so the chosen baseline remains visible
after generation.
