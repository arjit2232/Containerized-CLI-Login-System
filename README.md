# osto-login-cli

A containerized CLI login system in Go, backed by MySQL, with bcrypt password
hashing, TOTP-based 2FA, and login lockout after repeated failures.

## Prerequisites

- Go 1.22+ (only needed if running outside Docker)
- Docker + Docker Compose

## Quick Start (Docker)

```bash
git clone <repo>
cd osto-login-cli
docker compose -f docker/docker-compose.yml up --build -d db
docker compose -f docker/docker-compose.yml run --rm cli
```

No manual `.env` creation or config edits needed — `.env` ships with dummy
local credentials (see `.env.example` for the documented contract) and the
CLI runs pending migrations automatically on startup.

## Running locally (without Docker)

```bash
go mod tidy      # requires network access to the Go module proxy
make build
./cli
```

You'll need a reachable MySQL 8 instance and a `.env`/exported `DB_URL`
pointing at it.

## Command Reference

**Before login**
| Command | Description |
|---|---|
| `register` | Create a new account (username + password) |
| `login` | Sign in (prompts for 2FA code if enabled) |
| `help` | Show available commands |
| `exit` | Quit the CLI |

**After login**
| Command | Description |
|---|---|
| `whoami` | Show the current user and session info |
| `enable-2fa` | Turn on TOTP 2FA (scan QR / enter secret manually, then confirm with a code) |
| `disable-2fa` | Turn off 2FA (requires password + a current code) |
| `logout` | End the session |

## Setting up 2FA

Run `enable-2fa` while logged in. The CLI prints an `otpauth://` URI —
paste it into any online QR code generator and scan it with Google
Authenticator, Authy, or similar, or type the printed secret into the app
manually. Enter the 6-digit code it shows you to confirm and activate 2FA.

## Lockout policy & session timeout

Configurable via `.env`:

| Variable | Default | Meaning |
|---|---|---|
| `LOCKOUT_THRESHOLD` | 5 | Failed login attempts before lockout |
| `LOCKOUT_DURATION_MINUTES` | 15 | How long an account stays locked |
| `SESSION_TIMEOUT_MINUTES` | 20 | How long a session stays valid after login |

## Running migrations manually

Migrations run automatically on every CLI startup via `golang-migrate`.
There's no separate migrate step to run by hand.

## Tests

```bash
make test
```

`tests/auth_password_test.go`, `auth_totp_test.go`, and `auth_lockout_test.go`
are pure unit tests with no external dependencies. `tests/repository_test.go`
is skipped automatically unless `DB_URL` is set and reachable — run it
against the Dockerized `db` service or any MySQL 8 instance with the schema
migrated.

## Notes

- `go build ./...`, `go vet ./...`, and the pure-logic unit tests
  (`auth_password_test.go`, `auth_totp_test.go`, `auth_lockout_test.go`) were
  run and passed in the environment that generated this code. `go.mod`
  includes a small `replace` block redirecting a few `golang.org/x/*`
  modules to their GitHub mirrors — needed only because that sandbox had no
  route to `proxy.golang.org`; delete it freely if building somewhere with
  normal network access.
- `repository_test.go` needs a live MySQL instance (set `DB_URL`) and was
  not exercised in the sandbox, since no MySQL server was running there.
