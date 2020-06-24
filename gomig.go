// Package gomig TODO: add documentation
package gomig

import (
	"context"
	"errors"
)

var backend BackendAdapter

var (
	// ErrNoBackend returned by runners if no BackendAdapter set
	ErrNoBackend = errors.New("NO_BACKEND")
	// ErrNoMigrations returns by NewRunner if no migrations provided
	ErrNoMigrations = errors.New("NO_MIGRATIONS")
	// ErrMigrationFailed returned by runners if migration was failed
	ErrMigrationFailed = errors.New("MIGRATION_FAILED")
)

// Migration TODO: Add documentation
type Migration interface {
	Up(ctx context.Context) (messages []string, err error)
	Down(ctx context.Context) (messages []string, err error)
	Name() string
}

// BackendAdapter TODO: Add documentation
type BackendAdapter interface {
	GetVersion(ctx context.Context) (uint, error)
	SaveVersion(ctx context.Context, ver uint, name string, messages []string) error
	Reset(ctx context.Context) error
}
