## 📄 ADR-016: Error Handling Style

**Tags:** `style`, `errors`, `readability`

## Status
✅ Accepted

### Context

Go allows short variable declarations (`if err := ...`) and scoped error checks, 
but these can reduce traceability in larger or more modular functions.

### Decision

Use consistent:

```go
err = ...
if err != nil {
    return ...
}
```

for all error paths — no inline `if err :=` declarations.

### Advantages

* Easier to grep or scan for `err` across large files
* Error lifetime is clear across blocks
* Helps standardize flow and reduce surprises in reviews

### Disadvantages

* Slightly more verbose
* Requires discipline and enforcement in contributions
