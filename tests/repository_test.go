package tests

// These tests hit a real MySQL database and are skipped unless DB_URL is
// set (e.g. when running via `docker compose run --rm cli go test ./tests/...`
// against the docker-compose `db` service, or any other reachable MySQL
// instance with the schema already migrated).

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"osto-login-cli/internal/auth"
	"osto-login-cli/internal/repository"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		t.Skip("DB_URL not set, skipping repository integration tests")
	}
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("opening db: %v", err)
	}
	if err := pool.Ping(); err != nil {
		t.Skipf("could not reach db, skipping: %v", err)
	}
	return pool
}

func TestUserRepositoryCreateAndFetch(t *testing.T) {
	pool := openTestDB(t)
	defer pool.Close()

	repo := repository.NewUserRepository(pool)
	username := "testuser_" + uuid.NewString()[:8]
	hash, _ := auth.HashPassword("password123")

	id, err := repo.CreateUser(username, hash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if id == "" {
		t.Fatal("expected a non-empty generated ID")
	}

	fetched, err := repo.GetByUsername(username)
	if err != nil {
		t.Fatalf("GetByUsername failed: %v", err)
	}
	if fetched.ID != id {
		t.Errorf("expected ID %s, got %s", id, fetched.ID)
	}
	if fetched.TOTPEnabled {
		t.Error("expected TOTPEnabled to default to false")
	}
}

func TestSessionRepositoryCreateAndFetch(t *testing.T) {
	pool := openTestDB(t)
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	sessRepo := repository.NewSessionRepository(pool)

	username := "testuser_" + uuid.NewString()[:8]
	hash, _ := auth.HashPassword("password123")
	userID, err := userRepo.CreateUser(username, hash)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	sess := auth.NewSession(userID, 20*time.Minute)
	if err := sessRepo.Create(sess); err != nil {
		t.Fatalf("Create session failed: %v", err)
	}

	active, err := sessRepo.GetActiveSession(sess.Token)
	if err != nil {
		t.Fatalf("GetActiveSession failed: %v", err)
	}
	if active.UserID != userID {
		t.Errorf("expected session user_id %s, got %s", userID, active.UserID)
	}

	if err := sessRepo.RevokeSession(sess.ID); err != nil {
		t.Fatalf("RevokeSession failed: %v", err)
	}
	if _, err := sessRepo.GetActiveSession(sess.Token); err == nil {
		t.Error("expected revoked session to no longer be active")
	}
}
