package gomig_test

import (
	"errors"
	"testing"

	"github.com/pereslava/gomig"
)

func TestSetVer(t *testing.T) {
	t.Run("SetVer fails if no backend set", func(t *testing.T) {
		r := gomig.NewRunner(nil, nil)

		if err := r.SetVer(0); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})

	t.Run("Returns ErrNoMigrations if len of migrations array is less or erqual to target version", func(t *testing.T) {
		backend, migs, r := setup()
		resetMigrations(migs)
		backend.Reset()
		b := backend.(*backend_mock)
		b.currentVer = 5
		if err := r.SetVer(11); !errors.Is(err, gomig.ErrNoMigrations) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoMigrations, err)
		}
		verifyRunnerResults(t, b, migs, &verifyPattern{
			migs:   make([]migDir, 10),
			curVer: 5,
		})
	})

	t.Run("Run up", func(t *testing.T) {
		backend, migs, r := setup()
		b := backend.(*backend_mock)
		pattern := &verifyPattern{}
		t.Run("Run all", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				b.Reset()
				b.currentVer = uint(i)
				resetMigrations(migs)
				if err := r.SetVer(10); err != nil {
					t.Errorf("TestAuto failed: %v", err)
				}
				pattern.Set(i, 10, upRun, 10)
				verifyRunnerResults(t, b, migs, pattern)
			}
		})
		t.Run("Till version", func(t *testing.T) {
			for i := 1; i <= 10; i++ {
				pattern.Set(0, i, upRun, uint(i))
				b.Reset()
				resetMigrations(migs)
				if err := r.SetVer(uint(i)); err != nil {
					t.Error(err)
				}
				verifyRunnerResults(t, b, migs, pattern)
			}
		})
		t.Run("With fail", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				pattern.Set(0, i, upRun, uint(i))
				b.Reset()
				resetMigrations(migs)
				migs[i].(*migration_mock).fail = errors.New("SomeError")
				if err := r.SetVer(10); !errors.Is(err, gomig.ErrMigrationFailed) {
					t.Errorf("Want: %v, Have: %v", gomig.ErrMigrationFailed, err)
				}
				verifyRunnerResults(t, b, migs, pattern)
			}
		})
	})

	t.Run("Run down", func(t *testing.T) {
		backend, migs, r := setup()
		b := backend.(*backend_mock)
		pattern := &verifyPattern{}
		t.Run("Run all", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				b.Reset()
				b.currentVer = 10
				resetMigrations(migs)
				t.Log("Iteration", i)
				if err := r.SetVer(uint(i)); err != nil {
					t.Errorf("TestAuto failed: %v", err)
				}
				pattern.Set(10, i, downRun, uint(i))
				verifyRunnerResults(t, b, migs, pattern)
			}
		})
		t.Run("From version", func(t *testing.T) {
			for i := 1; i <= 10; i++ {
				t.Log("Iteration", i)
				b.Reset()
				b.currentVer = uint(i)
				resetMigrations(migs)

				if err := r.SetVer(0); err != nil {
					t.Error(err)
				}
				pattern.Set(0, i-1, downRun, uint(i-1))
				t.Log(pattern)
				verifyRunnerResults(t, b, migs, pattern)
			}
		})
		// t.Run("With fail", func(t *testing.T) {
		// 	for i := 0; i < 10; i++ {
		// 		pattern.Set(0, i, upRun, uint(i))
		// 		b.Reset()
		// 		resetMigrations(migs)
		// 		migs[i].(*migration_mock).fail = errors.New("SomeError")
		// 		if err := r.SetVer(10); !errors.Is(err, gomig.ErrMigrationFailed) {
		// 			t.Errorf("Want: %v, Have: %v", gomig.ErrMigrationFailed, err)
		// 		}
		// 		verifyRunnerResults(t, b, migs, pattern)
		// 	}
		// })

	})
}
