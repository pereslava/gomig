package gomig

import (
	"fmt"
)

// Runner the runner of migrations
type Runner struct {
	migs    []Migration
	storage BackendAdapter
}

// NewRunner creates the runner
func NewRunner(migs []Migration, storage BackendAdapter) *Runner {
	return &Runner{migs, storage}
}

// Auto runs the migrations from the current version up to the latest
func (r *Runner) Auto() error {
	if r.storage == nil {
		return ErrNoBackend
	}
	v, err := r.storage.GetVersion()
	if err != nil {
		return fmt.Errorf("storage GetVersion failed: %v, %w", err, ErrMigrationFailed)
	}
	if v > uint(len(r.migs)) {
		return fmt.Errorf("Auto failed: %w", ErrNoMigrations)
	}
	return r.runUp(v, uint(len(r.migs)))
}

// Reset runs all migrations down to the clean state and calls to Reset of storage
func (r *Runner) Reset() error {
	if r.storage == nil {
		return ErrNoBackend
	}
	if err := r.SetVer(0); err != nil {
		return err
	}
	return r.storage.Reset()
}

// SetVer runs migrations up or down to the ver number
func (r *Runner) SetVer(ver uint) error {
	if r.storage == nil {
		return ErrNoBackend
	}
	v, err := r.storage.GetVersion()
	if err != nil {
		return fmt.Errorf("storage GetVersion failed: %v, %w", err, ErrMigrationFailed)
	}
	switch {
	case int(ver) > len(r.migs):
		return fmt.Errorf("SetVer failed: %w", ErrNoMigrations)
	case v > ver:
		return r.runDown(v-1, ver)
	case v < ver:
		return r.runUp(v, ver)
	default:
		return fmt.Errorf("SetVer failed: %w", ErrNoMigrations)
	}
}

// ForceVer calls Reset then runs migrations up to the ver number
func (r *Runner) ForceVer(ver uint) error {
	if r.storage == nil {
		return ErrNoBackend
	}
	if err := r.Reset(); err != nil {
		return err
	}
	return r.SetVer(ver)
}

func (r *Runner) runUp(from, to uint) error {
	if r.storage == nil {
		return ErrNoBackend
	}

	for i := from; i < to; i++ {
		log, err := r.migs[i].Up()
		if err != nil {
			return fmt.Errorf("%v, %w", err, ErrMigrationFailed)
		}
		r.storage.SaveVersion(i+1, r.migs[i].Name(), log)
	}
	return nil
}

func (r *Runner) runDown(from, to uint) error {
	if r.storage == nil {
		return ErrNoBackend
	}
	for i := int(from); i >= int(to); i-- {
		log, err := r.migs[i].Down()
		if err != nil {
			return fmt.Errorf("%v, %w", err, ErrMigrationFailed)
		}
		r.storage.SaveVersion(uint(i), r.migs[i].Name(), log)
	}
	return nil
}
