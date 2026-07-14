# Threat Model

This threat model covers the goboot CLI, its committed configuration and
templates, generated output, and the release pipeline. It records the v0.9
security audit boundary before the first v1 release candidate.

## Assets and trust boundaries

Protected assets are user-owned files in an existing target, configuration
integrity, generated source and CI policy, release binaries, signing identity,
and provenance evidence.

Trust crosses these boundaries:

- YAML and template files enter the local generator process.
- Generated paths cross from staging into a user-selected target directory.
- Generated CI invokes third-party actions, images, and dependency tools.
- Tag workflows receive an OIDC identity and publish immutable release assets.
- Consumers download archives and verify them outside the repository.

## Threats and controls

| Threat | Existing control | Residual risk |
| --- | --- | --- |
| Target traversal or unsafe file type | Strict target joins, staging, symlink and special-file rejection | A compromised local account can alter files after validation |
| Overwriting user work | Ownership manifest, content/mode digests, managed conflicts, transactional rollback | `replace` is intentionally destructive when explicitly selected |
| Malicious or malformed configuration | Strict YAML decoding, validators, schemas, validation-only mode | Trusted templates can still generate unsafe application logic |
| Nondeterministic generation | Stable service ordering, bounded concurrency, golden paths, repeat-generation dogfood gate | Tool output such as `go mod tidy` can vary across dependency ecosystems |
| Dependency or CI compromise | Pinned action SHAs, pinned scanner/image versions, CodeQL, vulnerability and license checks | Upstream compromise before pinning and transitive build-tool risk remain |
| Release asset substitution | SHA-256 manifest, keyless signatures, GitHub provenance attestations, immutable tag policy | Consumers who skip verification receive no integrity guarantee |
| Credential disclosure | Least-privilege job permissions and GitHub OIDC; no long-lived signing key | Repository administrators and compromised workflows remain privileged |
| Denial of service | Parallelism cap, documented benchmark envelope, test timeouts | Input size is not globally capped; local disk and memory can be exhausted |

## Security invariants

- Generation never follows symbolic links in staged or target trees.
- Managed regeneration never replaces changed user-owned content silently.
- JSON command output remains isolated from logs.
- Release signing uses ephemeral OIDC identity; no signing private key is stored
  in the repository or Actions secrets.
- Every release archive is covered by a checksum and SBOM, and published assets
  receive workflow provenance.
- Pull-request workflows do not receive release write permissions.

## Out of scope

Goboot does not sandbox templates, audit the application behavior of generated
code, secure a compromised workstation, or guarantee third-party runner and
registry availability. Generated repositories retain responsibility for secret
management, deployment policy, and review of enabled automation.

## Audit cadence

Review this model for each new trust boundary, release publisher, template
execution capability, or post-v1 breaking change. Security-relevant changes must
update this document, tests, and the private-reporting guidance together.
