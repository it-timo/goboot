# Configuration Schemas

The JSON Schema files in this directory define the frozen v1 candidate shape
for the root goboot YAML file and every built-in service configuration.

Each committed example under `configs/` declares its schema with a
`yaml-language-server` modeline. Editors with YAML language-server support use
that declaration for completion, enum suggestions, and inline diagnostics
without repository-specific editor settings.

Schemas use JSON Schema draft 2020-12, reject unknown properties, and track the
same field names accepted by goboot's strict YAML decoder. Runtime validation
remains authoritative for rules that depend on multiple files or enabled
services.
