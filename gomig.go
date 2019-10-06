// Package gomig TODO: add documentation
package gomig

import (
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
	Up() (messages []string, err error)
	Down() (messages []string, err error)
}

// BackendAdapter TODO: Add documentation
type BackendAdapter interface {
	GetVersion() (uint, error)
	SaveVersion(ver uint, messages []string) error
	Reset() error
}
