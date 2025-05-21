## 📄 ADR-011: Service Responsibility & Scope – `baseProject`

**Tags:** `service`, `responsibility`, `structure`

## Status
✅ Accepted

### Context

The `baseProject` package defines a service that bootstraps a new Go project from a predefined directory of templates. 
It is the default entry point in the `goboot` project system and responsible for rendering structure and injecting metadata.

### Decision

Implement the `base_project` service as a concrete type that:

* Operates based on `config.BaseProjectConfig`
* Renders both paths and file contents using Go templates
* Performs all operations within an isolated `*os.Root`

### Advantages

* Keeps responsibility narrow and testable
* Integrates cleanly into the `goboot` service lifecycle
* Supports isolated, deterministic generation of project output
* Enables extensibility in a controlled, concrete way

### Disadvantages

* Any future support for different output formats or modes will need manual extension
* No plugin support (by design) means reduced flexibility unless rewritten
