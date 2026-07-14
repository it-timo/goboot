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
are derived from Git history and archives include a checksum manifest.

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

Tags are immutable release inputs. Correct mistakes with a new patch version;
do not replace a published tag.
