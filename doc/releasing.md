# Release Automation

`goboot` and projects generated with `base_release` use GoReleaser v2 for
tag-driven binary releases. Version calculation is deliberately explicit: the
Git tag is the release version. The automation does not guess the next version
or commit version bumps to the repository.

## Generated Output

The service writes:

- `.goreleaser.yml` with reproducible build targets, archives, checksums, and
  changelog rules
- `RELEASE.md` with the project release procedure
- a provider job registered with `base_ci` as either
  `.github/workflows/release.yml` or `.gitlab/ci/release.yml`

Linux, macOS, and Windows binaries are built for AMD64 and ARM64. Release notes
are derived from Git history. Archives include SHA-256 checksums and SBOMs;
binaries and the checksum manifest receive keyless Sigstore bundles, and GitHub
publishes build-provenance attestations for the release assets.

## Creating a Release

Run the normal release gates, create an annotated semantic-version tag, and
push the tag:

```bash
make release_check
git tag -a v0.2.1 -m "Release v0.2.1"
git push origin v0.2.1
```

GitHub publishes through the repository `GITHUB_TOKEN`. GitLab requires a
masked `GITLAB_TOKEN` with API permission for GoReleaser publishing.

Test the configuration without publishing:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

Snapshot acceptance, installation, upgrade, and dogfood procedures are defined
in [`release-candidate.md`](./release-candidate.md). Snapshot builds skip signing
because pull requests do not receive release identity; real tag releases require
OIDC signing and attestation permissions.

After downloading a release, verify the checksum bundle and provenance:

```bash
cosign verify-blob --bundle checksums.txt.sigstore.json checksums.txt
gh attestation verify goboot_VERSION_linux_amd64.tar.gz -R it-timo/goboot
sha256sum --check checksums.txt
```

Tags are immutable release inputs. Correct mistakes with a new patch version;
do not replace a published tag.
