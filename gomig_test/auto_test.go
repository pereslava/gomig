package gomig_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pereslava/gomig"
)

func TestAuto(t *testing.T) {
	ctx := context.Background()
	t.Run("Fails if no backend set", func(t *testing.T) {
		r := gomig.NewRunner(nil, nil)

		if err := r.Auto(ctx); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})
	backend, migs, r := setup()
	b := backend.(*backend_mock)
	b.Reset(ctx)
	resetMigrations(migs)

	t.Run("Run all migrations", func(t *testing.T) {
		pattern := &verifyPattern{
			seq:    []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			migs:   []migDir{upRun, upRun, upRun, upRun, upRun, upRun, upRun, upRun, upRun, upRun},
			curVer: 10,
		}
		for i := 0; i < 10; i++ {
			b.Reset(ctx)
			b.currentVer = uint(i)
			resetMigrations(migs)
			if err := r.Auto(ctx); err != nil {
				t.Errorf("TestAuto failed: %v", err)
			}
			pattern.Set(i, 10, upRun, 10)
			verifyRunnerResults(ctx, t, b, migs, pattern)
		}
	})

	b.Reset(ctx)
	resetMigrations(migs)

	t.Run("With fail", func(t *testing.T) {
		pattern := &verifyPattern{}
		for i := 0; i < 10; i++ {
			pattern.Set(0, i, upRun, uint(i))
			b.Reset(ctx)
			resetMigrations(migs)
			migs[i].(*migration_mock).fail = errors.New("SomeError")
			if err := r.Auto(ctx); !errors.Is(err, gomig.ErrMigrationFailed) {
				t.Errorf("Want: %v, Have: %v", gomig.ErrMigrationFailed, err)
			}
			verifyRunnerResults(ctx, t, b, migs, pattern)
		}
	})

	b.Reset(ctx)
	resetMigrations(migs)

	t.Run("Silent exit if curVer == len of migrations array", func(t *testing.T) {
		b.currentVer = 10
		if err := r.Auto(ctx); err != nil {
			t.Errorf("TestAuto failed: %v", err)
		}
		verifyRunnerResults(ctx, t, b, migs, &verifyPattern{
			seq:    []uint{},
			migs:   []migDir{notRun, notRun, notRun, notRun, notRun, notRun, notRun, notRun, notRun, notRun},
			curVer: 10,
		})
	})

	b.Reset(ctx)
	resetMigrations(migs)

	t.Run("Returns ErrNoMigrations if len of migrations array less from current version", func(t *testing.T) {
		b.currentVer = 11
		if err := r.Auto(ctx); !errors.Is(err, gomig.ErrNoMigrations) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoMigrations, err)
		}
		verifyRunnerResults(ctx, t, b, migs, &verifyPattern{
			migs:   make([]migDir, 10),
			curVer: 11,
		})
	})
}
