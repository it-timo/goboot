# Contributing to goboot

Thank you for contributing. goboot favors deterministic output, explicit
configuration, and narrowly owned services over hidden defaults.

## Before opening a change

1. Search existing issues and pull requests.
2. Open an issue before substantial behavioral or architectural changes.
3. Keep the change focused and add tests for observable behavior.
4. Run the full local quality gates:

   ```bash
   make lint
   make test
   make build
   ```

Changes to generated output should include the relevant template, config,
service tests, end-to-end assertions, documentation, and an ADR when they alter
an architectural contract.

## Pull requests

- Use a short, imperative commit subject.
- Explain motivation, user impact, and compatibility considerations.
- Link the relevant issue when one exists.
- Include validation evidence and update documentation.
- Never commit credentials, tokens, private keys, or production data.

## Security reports

Do not open a public issue for a suspected vulnerability. Follow
[SECURITY.md](./SECURITY.md) to report it privately.
