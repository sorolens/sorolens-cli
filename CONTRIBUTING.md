# Contributing to sorolens-cli

Thank you for helping improve sorolens-cli. This guide covers everything you
need to open a pull request.

---

## Table of contents

- [Development setup](#development-setup)
- [Running tests](#running-tests)
- [Adding a new command](#adding-a-new-command)
- [Branch naming](#branch-naming)
- [Commit messages](#commit-messages)
- [Pull request checklist](#pull-request-checklist)

---

## Development setup

Requirements: Go 1.23 or later, `golangci-lint`, `goreleaser` (optional,
for release dry-runs).

```bash
git clone https://github.com/sorolens/sorolens-cli.git
cd sorolens-cli
go mod download
make build        # compile the binary
make install      # install to $GOPATH/bin
```

Copy the example environment file and point it at a local Sorolens API:

```bash
cp .env.example .env
# edit .env and set SOROLENS_API_URL
```

---

## Running tests

```bash
# All tests
make test

# With race detector (matches CI)
go test -race ./...

# Regenerate golden files after an intentional output change
make update-golden

# Lint
make lint
```

---

## Adding a new command

1. Create `cmd/<name>.go`. Skeleton:

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/sorolens/sorolens-cli/internal/format"
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand <contract-id>",
    Short: "One-line description shown in --help",
    Args:  cobra.ExactArgs(1),
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    ctx := cmd.Context()
    data, err := apiClient.GetSomething(ctx, args[0])
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        return err
    }
    if flagJSON {
        return format.PrintJSON(nil, data)
    }
    // render a table
    return nil
}
```

2. Register it in `init()` with `rootCmd.AddCommand(myCmd)`.

3. Honour the global `--json` flag: when set, call `format.PrintJSON(nil,
   data)` instead of rendering a table.

4. Write all errors to `os.Stderr` and return a non-nil error (exit code 1).

5. Add a method to `internal/client/client.go` for each new API endpoint,
   following the pattern of `GetContract`, `GetEvents`, and `GetStorage`.

6. Add an `httptest` test for the new client method in
   `internal/client/client_test.go`.

7. Add usage examples to `README.md` under the "Commands" section.

---

## Branch naming

```
feat/<short-description>     new feature
fix/<short-description>      bug fix
docs/<short-description>     documentation only
test/<short-description>     tests only
chore/<short-description>    build, deps, tooling
```

Examples:

```
feat/events-csv-output
fix/ttl-color-threshold
docs/shell-completion-readme
```

---

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

```
feat(cmd): add --csv flag to events command

fix(client): retry on 429 with Retry-After header

docs(readme): add shell completion instructions

test(format): add golden file for TTL table output

chore(deps): upgrade cobra to v1.8.1
```

Rules:
- Subject line: 72 characters or fewer
- Use imperative mood: "add", "fix", "update" -- not "added" or "fixes"
- Reference the relevant issue with `Closes #N` in the commit body

---

## Release secrets

Two repository secrets are required for the release workflow:

| Secret | Purpose |
|--------|---------|
| `HOMEBREW_TAP_TOKEN` | GitHub PAT with write access to `sorolens/homebrew-tap` |
| `SCOOP_BUCKET_TOKEN` | GitHub PAT with write access to `sorolens/scoop-bucket` |

Contact a maintainer to be granted access.

---

## Pull request checklist

- [ ] `go build ./...` passes
- [ ] `go test -race ./...` passes
- [ ] `make lint` passes with no new warnings
- [ ] New commands have `httptest` coverage in `internal/client/client_test.go`
- [ ] New table renderers have a golden file under `testdata/golden/`
- [ ] `README.md` updated with any new flags or commands
- [ ] Commit messages follow Conventional Commits
