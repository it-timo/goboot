## 📄 ADR-001: Minimalism over Generalization

**Tags:** `design-philosophy`, `minimalism`, `abstraction`

## Status
✅ Accepted

### Context

Some Go projects tend toward early interface design, reflection, or generic pipelines. 
This can lead to unnecessary abstraction and complexity.

### Decision

This package uses:

* Concrete types
* Clear linear flow (no plugin maps, no injected handlers)
* Utility helpers (`utils.EnsureDir`) where repetition warrants it

Interfaces are deferred until actually needed.

### Advantages

* High readability and traceability for OSS contributors
* No unnecessary indirection for trivial or unique logic
* Reduces the surface area for bugs or misuse

### Disadvantages

* If support for pluggable templates or dynamic output is added later, refactoring may be needed
* Contributors unfamiliar with Go’s idioms might expect more abstraction
