# Release-Candidate Acceptance

Version 0.9 turns the v1 candidate into a repeatable acceptance process. A
release candidate is evidence that the frozen CLI and configuration contracts
can be built, installed, upgraded, generated, and verified; it is not merely a
version label.

## Pull-request gate

The `Release Candidate` workflow runs on every pull request and requires:

1. A current candidate binary built with the repository's pinned Go version.
2. Replacement of a binary built from the pull request base revision.
3. Successful `--version`, `--help`, and JSON configuration validation after
   replacement.
4. Two complete generations from the committed goboot configuration with
   byte-identical outputs.
5. A GoReleaser snapshot for every supported OS and architecture.
6. Valid SHA-256 checksums, archive extraction, executable smoke tests, and an
   SBOM for each release archive.
7. A documentation audit covering the public CLI, performance, regeneration,
   release, supply-chain, and threat-model contracts with no stale status claim.

The normal test workflow also compares base-project output paths with the
checked-in `testdata/golden/base_project.paths` contract. Intentional output
shape changes must update that file in the same pull request.

## Release assets

Tag releases publish Linux, macOS, and Windows archives for AMD64 and ARM64.
GoReleaser generates SHA-256 checksums and archive SBOMs, signs binaries and the
checksum manifest through keyless Sigstore identities, and GitHub records build
provenance attestations for the published assets.

Verify a checksum signature and GitHub provenance after downloading an asset:

```bash
cosign verify-blob \
  --bundle checksums.txt.sigstore.json \
  checksums.txt
gh attestation verify goboot_VERSION_linux_amd64.tar.gz -R it-timo/goboot
```

Consumers must still compare the archive digest with `checksums.txt`. A valid
signature identifies the release workflow; it does not replace vulnerability
assessment or local policy.

## v1 candidate sequence

Candidate tags use `v1.0.0-rc.N`. GoReleaser marks semantic-version prerelease
tags as prereleases automatically. For every candidate:

- run the full branch-protection suite and release-candidate workflow;
- review golden-output changes and the threat model;
- install one archive on a clean supported host;
- upgrade one existing installation;
- verify checksums, Sigstore bundles, provenance, and SBOM presence;
- record dogfood or external-repository feedback in the release notes.

Promote to `v1.0.0` only when no unresolved release-blocking defect remains.
Correct a published candidate with a new `rc.N`; never replace its tag or assets.

## Known pre-v1 boundary

The acceptance suite proves repository-controlled configurations and supported
release targets. It does not claim six months of compatibility history or broad
multi-project adoption; those remain explicit v1 exit criteria in the roadmap.
