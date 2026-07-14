# Governance and Contribution Generation

The `base_governance` service creates deterministic repository ownership,
contribution, and security-policy files. It is optional and never changes
provider settings, permissions, branch protection, or repository visibility.

## Configuration

```yaml
sourcePath: "./templates/governance_base"
maintainers:
  - "@example-maintainer"
defaultBranch: "main"
```

Maintainers must be explicit provider handles. GitHub team handles such as
`@example-org/maintainers` and GitLab group handles are supported. Duplicate or
malformed entries fail validation instead of producing ambiguous ownership.

## Common output

Every enabled governance profile receives:

- `CODEOWNERS` with repository-wide ownership;
- `CONTRIBUTING.md` with local validation and review expectations;
- `SECURITY.md` directing vulnerability reports to the provider's private channel.

Security templates deliberately do not suggest public issues or store private
contact details in scaffold configuration.

## Provider output

| Provider | Review template                              | Issue templates        |
| -------- | -------------------------------------------- | ---------------------- |
| GitHub   | `.github/PULL_REQUEST_TEMPLATE.md`           | GitHub issue forms     |
| GitLab   | `.gitlab/merge_request_templates/Default.md` | GitLab issue templates |

## Profile baselines

| Profile    | Generated issue workflows                                      |
| ---------- | -------------------------------------------------------------- |
| Minimal    | Bug report only                                                |
| Standard   | Bug report and feature request                                 |
| Enterprise | Standard plus controlled change request                        |
| OSS        | Standard plus documentation report and GitHub chooser policy   |

Profiles select only files owned by `base_governance`. They do not create labels,
enable private vulnerability reporting, enforce approvals, or overwrite explicit
repository settings.

## Operational boundary

Generated `CODEOWNERS` expresses intended reviewers, but enforcement still belongs
to the hosting provider. Maintainers must configure protected branches, required
approvals, and private reporting in GitHub or GitLab after repository creation.
