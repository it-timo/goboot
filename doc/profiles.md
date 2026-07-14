# Template Profiles

Template profiles provide named lint, test, and documentation baselines without
changing goboot's explicit service composition model.

Select a profile in the root configuration:

```yaml
profile: standard
```

Supported values are `minimal`, `standard`, `enterprise`, and `oss`. If the
field is omitted, goboot uses `standard`.

## Baselines

| Profile | Default tests | Test command | Go lint baseline |
| ------- | ------------- | ------------ | ---------------- |
| `minimal` | Standard library | Fast tests without race or coverage | Small correctness/security set; complexity 20 |
| `standard` | Ginkgo/Gomega | Race detection and coverage | Full balanced set; complexity 10 |
| `enterprise` | Ginkgo/Gomega | Race detection, shuffled order, coverage | Full strict set; complexity 8 |
| `oss` | Ginkgo/Gomega | Race detection and atomic coverage | Full public-project set; complexity 10 |

## Precedence

Profiles fill defaults and provide template baselines. Explicit service values
remain authoritative. For example, setting `useStyle: go` or a custom
`testCmd` in `base_test.yml` overrides the selected profile default.

Profiles do not enable or disable services. The root `services` list remains the
only source of service composition, preserving deterministic and auditable runs.

Generated projects contain `PROFILE.md` so the chosen baseline remains visible
after generation.
