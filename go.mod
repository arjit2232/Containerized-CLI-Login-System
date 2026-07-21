module osto-login-cli

go 1.22

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-migrate/migrate/v4 v4.17.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/pquerna/otp v1.4.0
	golang.org/x/crypto v0.24.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/chzyer/test v1.0.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

// Redirects a few golang.org/x/* (and other vanity-import) modules to their
// official GitHub mirrors. This was needed to build in the sandbox that
// generated this project (no route to golang.org/proxy.golang.org), and is
// harmless to keep — the mirrored code is identical to the canonical module.
// Safe to delete this block if you're building somewhere with normal access
// to proxy.golang.org; `go mod tidy` will then resolve everything directly.
replace (
	filippo.io/edwards25519 => github.com/FiloSottile/edwards25519 v1.1.0
	go.uber.org/atomic => github.com/uber-go/atomic v1.7.0
	golang.org/x/crypto => github.com/golang/crypto v0.24.0
	golang.org/x/sys => github.com/golang/sys v0.21.0
	golang.org/x/term => github.com/golang/term v0.21.0
)
