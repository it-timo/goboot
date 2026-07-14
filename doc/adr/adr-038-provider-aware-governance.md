# ADR-038: Provider-Aware Governance as an Explicit Service

**Tags:** `governance`, `templates`, `security`, `github`, `gitlab`

---

## Status

Accepted

---

## Context

Contribution files are part of a repository's behavioral contract. Keeping them
only in the goboot source repository would leave generated projects without clear
ownership, reporting, and review paths. Provider-specific file locations also make
them unsuitable for unconditional generation by `base_project`.

---

## Decision

Introduce an optional `base_governance` service. It validates maintainer handles,
renders common ownership and security documents, and selects GitHub or GitLab
contribution templates from the validated root provider. The selected project
profile controls only the breadth of issue workflows.

The service generates policy text but never mutates remote repository settings.
Vulnerability guidance uses private provider channels and does not encourage public
security issues.

---

## Advantages

- Generated repositories begin with explicit ownership and contribution contracts.
- GitHub and GitLab receive native template layouts without unused provider files.
- Profile differences remain deterministic and documented.
- Remote permission and branch-protection changes stay outside generator scope.

---

## Disadvantages

- Maintainers must still enable provider-side enforcement and security features.
- Provider template conventions add files that require ongoing compatibility review.
- A dedicated service adds configuration and orchestration surface.

---

## Alternatives Considered

- Put all governance files in `base_project`; rejected because provider-specific
  directories would be emitted together or require hidden deletion rules.
- Configure provider settings through APIs; rejected because it requires credentials,
  network access, and permissions outside deterministic filesystem generation.
- Generate only generic Markdown; rejected because native issue forms and merge
  request templates materially improve contribution quality.
