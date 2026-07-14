# ADR-039: Explicit Supply-Chain Security Service

**Tags:** `security`, `ci`, `codeql`, `dependencies`, `sbom`

---

## Status

Accepted

## Context

Generated repositories need auditable source, dependency, license, and artifact
controls. Hiding those controls inside `base_ci` would mix pipeline rendering with
security policy and make versions or license exceptions difficult to review.

## Decision

Add `base_supplychain` as an explicit regular service. It validates scanner
versions and the license allowlist, renders `SUPPLY_CHAIN.md`, and registers a
provider-native `security.yml` file with the existing CI registrar.

GitHub uses CodeQL plus portable Go dependency checks and CycloneDX SBOM
generation. GitLab uses the portable checks and SBOM generation, while clearly
documenting that CodeQL is GitHub-specific. Actions use immutable commit SHAs;
command-line scanners and images use explicit release versions.

## Advantages

- Security policy is visible in config, generated documentation, and CI.
- Provider differences are explicit instead of silently approximated.
- Scanner updates and license changes produce reviewable diffs.
- The existing registrar boundary remains unchanged.

## Disadvantages

- Security jobs add CI time and download external scanners.
- GitHub and GitLab do not have identical source-analysis capabilities.
- Version pins require routine maintenance.

## Alternatives Considered

- Put the controls directly in `base_ci`; rejected because CI would own security policy.
- Run CodeQL on GitLab; rejected because CodeQL analysis and upload are GitHub-native.
- Use floating tool versions; rejected because generated pipelines must be reproducible.
