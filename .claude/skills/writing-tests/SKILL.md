---
name: writing-tests
description: Writes unit and integration tests using this project's helpers (MockDB, Serve, integration.Setup) and the test-first approval workflow. Use when the user asks for tests, test cases, TDD or coverage for a function, middleware, service, handler or endpoint, and before implementing any new behavior.
---

# Writing tests

```
- [ ] 1. Choose the type
- [ ] 2. Read the matching rule and the nearest existing test
- [ ] 3. Write the tests
- [ ] 4. Run them: they must fail for the expected reason
- [ ] 5. Show them and STOP until the user approves
- [ ] 6. Implement until green
```

## 1. Type

| Target                                                             | Type        | Location                                          |
| ------------------------------------------------------------------ | ----------- | ------------------------------------------------- |
| function, method, validation, service, middleware, errors, logging | unit        | `tests/unit/{path as in internal}/{file}_test.go` |
| HTTP endpoint through the real app and database                    | integration | `tests/integration/features/{feature}_test.go`    |

## 2. Context

Read `.claude/rules/unit-tests.md` or `.claude/rules/integration-tests.md`, then one existing test of the same kind. Skeletons are in [templates.md](templates.md).

## 3. Cases to cover

- **Handler/endpoint:** success, validation failure 400, missing or invalid token 401 (private routes), forbidden 403 (role/permission routes), not found 404, duplicate 409, and that the data was stored.
- **Service:** every branch, and the repository-error branch returning a 500-class error.
- **Validation:** each rule valid and invalid, plus normalization (trim, lowercase).
- **Middleware:** every rejection path and the pass-through path.

## 4. Failing first

A compile error for a symbol that does not exist yet counts as failing for the right reason. A failure caused by the harness, config or Docker does not: fix that first.

## 5. Approval gate

Show the test files and a list of case names. Do not write implementation code before the user approves. After approval, do not weaken or edit approved tests to make them pass: propose the change instead.

## 6. Run

- `go test -count=1 ./tests/unit/...`
- `go test -count=1 -v ./tests/integration/...` (needs Docker)
