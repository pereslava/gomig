package gomig

import "fmt"

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
	return r.runUp(v, len(r.migs))
}

// SetVer runs migrations up or down to the ver number
func (r *Runner) SetVer(ver int) error {
	if r.storage == nil {
		return ErrNoBackend
	}
	return nil
}

// ForceVer runs all down migrations from the current version to vertion 0
// then run migrations up to the ver number
func (r *Runner) ForceVer(ver int) error {
	if r.storage == nil {
		return ErrNoBackend
	}
	return nil
}

func (r *Runner) runUp(from, to int) error {
	if r.storage == nil {
		return ErrNoBackend
	}

	for i := from; i < to; i++ {
		log, err := r.migs[i].Up()
		if err != nil {
			return fmt.Errorf("%v, %w", err, ErrMigrationFailed)
		}
		r.storage.SaveVersion(i+1, log)
	}
	return nil
}
