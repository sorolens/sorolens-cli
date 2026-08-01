## Summary

<!-- What does this PR do and why? One to three sentences. -->

Closes #

## Type of change

- [ ] Bug fix
- [ ] New feature / command
- [ ] Output format change (table, JSON, CSV)
- [ ] Refactor (no behavior change)
- [ ] Documentation only
- [ ] CI / build / tooling

## Changes

<!-- Bullet list of the concrete things you changed. -->

-

## Testing

<!-- How did you verify this works? Paste relevant commands or test output. -->

```bash

```

## Checklist

- [ ] `go build ./...` passes
- [ ] `go test -race ./...` passes
- [ ] `make lint` passes with no new warnings
- [ ] New commands have tests in `internal/client/client_test.go`
- [ ] New table renderers have a golden file under `testdata/golden/`
- [ ] `README.md` updated if new flags or commands were added
- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org)
- [ ] No secrets or real contract IDs in test fixtures
