# ADR-033: CI Policy Modes and Provider-Scoped Template Layout

**Tags:** `ci`, `templates`, `config`, `security`, `provider-layout`

---

## Status

Accepted

---

## Context

CI generation needed a clearer contract for policy behavior, provider layout, and render data consistency.
Trade-offs include reproducibility, usability, and provider-specific runtime constraints.

---

## Decision

- Support three image policies: `balanced` (default), `strict`, and `simple`.
- Keep provider templates under `sourcePath/<provider>/...`.
- Require typed config fields for policy-relevant behavior (`goVersions`, `autoBranches`, `imagePolicy`, job `allowFailure`).
- Keep strict digest placeholders explicit in generated output (`REPLACE_ME`) instead of hidden auto-fetch logic.
- Preserve deterministic render ordering.

---

## Advantages

- Clear CI policy contract across providers.
- Explicit behavior avoids hidden runtime fetch steps.
- Deterministic output simplifies review and diffs.

---

## Disadvantages

- Strict mode requires follow-up digest replacement workflow.
- Provider-specific limitations prevent fully identical behavior.
- Policy and template contract surface area increases maintenance work.

---

## Alternatives Considered

- **Strict-only mode with fixed digests:** rejected due to update burden and stale references.
- **Tag-only mode:** rejected because reproducibility/security controls are weaker.
- **Automatic digest resolution at generation/runtime:** rejected to keep behavior explicit and deterministic.
