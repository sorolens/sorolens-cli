# sorolens-cli

Terminal access to Soroban smart contract observability data.

[![CI](https://github.com/sorolens/sorolens-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/sorolens/sorolens-cli/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

---

## Install

**go install** (requires Go 1.23+)

```bash
go install github.com/sorolens/sorolens-cli@latest
```

**Direct binary download**

Download a pre-built binary for your platform from the
[releases page](https://github.com/sorolens/sorolens-cli/releases),
extract the archive, and place the `sorolens` binary somewhere on your `PATH`.

```bash
# Example: Linux amd64
curl -L https://github.com/sorolens/sorolens-cli/releases/latest/download/sorolens_linux_x86_64.tar.gz \
  | tar xz
sudo mv sorolens /usr/local/bin/
```

---

## Configuration

Set the API base URL with `--api-url` or `SOROLENS_API_URL` (default: `http://localhost:8080`).

```bash
export SOROLENS_API_URL=https://api.sorolens.dev
```

Copy `.env.example` to `.env` to persist local settings.

---

## Commands

### inspect

Show a summary table for a tracked contract.

```bash
sorolens inspect CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD2KM
sorolens inspect CAAAA... --json
```

Output: alias, network, first seen, event count, invocation count, success
rate, storage entries, nearest TTL expiry ledger.

### events

List or stream events for a contract.

```bash
# Latest 50 events
sorolens events CAAAA...

# Filter by event type
sorolens events CAAAA... --type invoke

# Cap the result set
sorolens events CAAAA... --limit 100

# Stream new events every 3 seconds (Ctrl+C to stop)
sorolens events CAAAA... --follow

# Export to RFC 4180 CSV
sorolens events CAAAA... --csv > events.csv
```

### ttl

Show storage entry TTLs with color-coded urgency.

```bash
sorolens ttl CAAAA...
sorolens ttl CAAAA... --json
```

Row colors:
- Green: safe (> 10,000 ledgers remaining)
- Yellow: warning (1,000-9,999 ledgers remaining)
- Red: danger (< 1,000 ledgers or expired)

### track

Register a contract for tracking.

```bash
sorolens track CAAAA...
sorolens track CAAAA... --alias my-defi-contract
sorolens track CAAAA... --alias my-defi-contract --json
```

### version

Print the CLI version.

```bash
sorolens version
```

---

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--api-url` | `http://localhost:8080` | API base URL (env: `SOROLENS_API_URL`) |
| `--json` | false | Machine-readable JSON output |
| `--no-color` | false | Disable color (also respects `NO_COLOR` env var) |
| `--timeout` | `10s` | Per-request timeout |

---

## Shell completion

```bash
# Bash
sorolens completion bash > /etc/bash_completion.d/sorolens

# Zsh
sorolens completion zsh > "${fpath[1]}/_sorolens"

# Fish
sorolens completion fish > ~/.config/fish/completions/sorolens.fish
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Contributors

Thanks to everyone who has contributed to sorolens-cli!

[![Contributors](https://contrib.rocks/image?repo=sorolens/sorolens-cli)](https://github.com/sorolens/sorolens-cli/graphs/contributors)

## License

Apache 2.0 -- see [LICENSE](LICENSE).
