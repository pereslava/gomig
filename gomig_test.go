package gomig_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pereslava/gomig"
)

func TestNoBackend(t *testing.T) {
	r := gomig.NewRunner(nil, nil)
	t.Run("Auto fails if no backend set", func(t *testing.T) {
		if err := r.Auto(); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})
	t.Run("SetVer fails if no backend set", func(t *testing.T) {
		if err := r.SetVer(0); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})
	t.Run("ForceVer fails if no backend set", func(t *testing.T) {
		if err := r.ForceVer(0); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})
}

func TestAuto(t *testing.T) {
	backend := &adapter_mock{}
	migs := make([]gomig.Migration, 10)
	for i, _ := range migs {
		migs[i] = &migMock{index: i}
	}
	r := gomig.NewRunner(migs, backend)

	t.Run("Run all migrations", func(t *testing.T) {
		backend.reset()
		resetMigrations(migs)
		if err := r.Auto(); err != nil {
			t.Errorf("Auto failed: %v", err)
		}
		s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		if !reflect.DeepEqual(s, backend.seq) {
			t.Errorf("Want: %v, Have: %v", s, backend.seq)
		}
		for _, m := range migs {
			mig := m.(*migMock)
			if mig.calledUp != 1 {
				t.Errorf("Migration %d not run", mig.index)
			}
		}

		if curVer, _ := backend.GetVersion(); curVer != 10 {
			t.Errorf("Check Current Version failed, Want: 10, Have: %d", curVer)
		}
	})

	resetMigrations(migs)
	backend.reset()
	migs[5].(*migMock).fail = errors.New("SomeError")

	t.Run("Version must be one before failed migration", func(t *testing.T) {
		err := r.Auto()
		if !errors.Is(err, gomig.ErrMigrationFailed) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrMigrationFailed, err)
		}
		curVer, _ := backend.GetVersion()
		if curVer != 5 {
			t.Errorf("CurrentVersion: Want: 5, Have: %d", curVer)
		}
		s := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(s, backend.seq) {
			t.Errorf("Want: %v, Have: %v", s, backend.seq)
		}
	})

	resetMigrations(migs)
	backend.reset()

	t.Run("From current to latest", func(t *testing.T) {
		backend.currentVer = 5
		err := r.Auto()
		if err != nil {
			t.Errorf("Auto failed: %v", err)
		}

		curVer, _ := backend.GetVersion()
		if curVer != 10 {
			t.Errorf("CurrentVersion: Want: 10, Have: %d", curVer)
		}
		s := []int{6, 7, 8, 9, 10}
		if !reflect.DeepEqual(s, backend.seq) {
			t.Errorf("Want: %v, Have: %v", s, backend.seq)
		}
	})
}
