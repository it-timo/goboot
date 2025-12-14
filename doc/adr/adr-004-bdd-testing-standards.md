# 📄 ADR-004: BDD Testing Standards and Tooling

**Tags:** `testing`, `bdd`, `standards`, `quality`

---

## Status

✅ Accepted

---

## Context

Go's standard `testing` package is lightweight and sufficient for simple libraries.
However, `goboot` involves complex integration logic, file system permutations, and configuration variations.
Using standard table-driven tests for these scenarios often leads to:

- Deeply nested `if` statements in test loops.
- Difficulty in isolating "setup" and "teardown" per test case.
- Lack of expressive assertions for complex states (e.g., file existence, content matching).

---

## Decision

We adopt **Behavior-Driven Development (BDD)** style testing using **Ginkgo** and **Gomega**.

### 1. Framework & Assertions

- **Ginkgo**: Use for structural organization (`Describe`, `Context`, `It`, `BeforeEach`).
- **Gomega**: Use for fluent assertions (`Expect(val).To(Equal(x))`).

### 2. Required Patterns

- **DescribeTable**: ALL permutation tests must use `DescribeTable` instead of manual `for` loops over structs.
This ensures that one failed regression doesn't crash the entire loop and provides clear output.
- **Isolation**: Every test spec (`It`) must be independent.
Use `GinkgoT().TempDir()` or `os.MkdirTemp` for **all** file system operations.
- **Prohibited**: Manually creating/cleaning up directories in `var _` blocks or using shared global state variables.

### 3. Golden Standards

- **No Mystery constants**: Use `goboottypes` constants where possible.
- **Explicit Setup**: `BeforeEach` should establish the "known good state" for the `Context`.

---

## Advantages

- **Readability**: Tests read like documentation ("It handles empty config by returning error").
- **Debuggability**: `DescribeTable` reports exactly which entry failed.
- **Safety**: Strict isolation prevents flaky tests caused by lingering files.

---

## Disadvantages

- **Learning Curve**: Contributors must learn the Ginkgo DSL.
- **Verbosity**: BDD specifications can be more verbose than simple Go test functions.
