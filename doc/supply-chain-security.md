# Supply-Chain Security

The optional `base_supplychain` service makes repository security controls explicit
and reproducible. It owns the generated `SUPPLY_CHAIN.md` policy and registers one
provider-native `security.yml` pipeline with `base_ci`.

## Controls

| Control | GitHub | GitLab |
| --- | --- | --- |
| Go source analysis | CodeQL | Not available; CodeQL is GitHub-specific |
| Reachable vulnerabilities | `govulncheck` | `govulncheck` |
| Dependency licenses | `go-licenses` allowlist | `go-licenses` allowlist |
| SBOM | CycloneDX JSON workflow artifact | CycloneDX JSON job artifact |

GitHub actions are pinned to immutable commit SHAs. Go scanners are pinned to
semantic versions in `configs/base_supplychain.yml`. The GitLab Syft image is
pinned to an explicit release and clears its container entrypoint for Docker
executor compatibility.

## Configuration

```yaml
sourcePath: "./templates/supplychain_base"
govulncheckVersion: "v1.1.4"
goLicensesVersion: "v2.0.1"
allowedLicenses:
  - "Apache-2.0"
  - "BSD-2-Clause"
  - "BSD-3-Clause"
  - "ISC"
  - "MIT"
```

Versions must be complete `vMAJOR.MINOR.PATCH` values. License values must be
safe SPDX-style identifiers and cannot be duplicated. Changes to either list
should be reviewed like source changes rather than applied as CI exceptions.

## Operational boundaries

- CodeQL requires GitHub code-scanning support and the workflow's
  `security-events: write` permission.
- GitHub default CodeQL setup and an advanced CodeQL workflow cannot run
  together. New generated projects use the advanced workflow; repositories
  with default setup already enabled must keep one setup and remove the other.
- The goboot repository disables GitHub default setup and uses the pinned
  advanced CodeQL job in its checked-in `security.yml` workflow.
- GitLab does not receive a fake CodeQL substitute. Its portable controls remain
  vulnerability, license, and SBOM scanning.
- Generated SBOMs are build artifacts; publishing or signing release assets is a
  separate release policy.
- Template profiles do not weaken these controls. Enable or disable the service
  explicitly in the root goboot configuration.
