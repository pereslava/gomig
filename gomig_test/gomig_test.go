package gomig_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pereslava/gomig"
)

type migDir int

const (
	notRun migDir = iota
	upRun
	downRun
)

type verifyPattern struct {
	seq    []uint
	migs   []migDir
	curVer uint
}

func (p *verifyPattern) Set(from, to int, dir migDir, curVer uint) {
	p.migs = make([]migDir, 10)
	if from <= to {
		p.seq = make([]uint, to-from)
		for i, _ := range p.seq {
			p.seq[i] = uint(i + from + 1)
		}
		for i := from; i < to; i++ {
			p.migs[i] = dir
		}
	} else {
		p.seq = make([]uint, from-to)
		for i := 0; i < len(p.seq); i++ {
			p.seq[i] = uint(10 - i - 1)
		}
		for i := from - 1; i >= to; i-- {
			p.migs[i] = dir
		}
	}
	p.curVer = curVer
}

func verifyRunnerResults(t *testing.T, b *backend_mock, migs []gomig.Migration, pattern *verifyPattern) {
	if !(len(pattern.seq) == 0 && len(b.seq) == 0) {
		if !reflect.DeepEqual(pattern.seq, b.seq) {
			t.Errorf("Versions sequence: Want: %v, Have: %v", pattern.seq, b.seq)
		}
	}
	if v, _ := b.GetVersion(); v != pattern.curVer {
		t.Errorf("Check Current Version failed, Want: %d, Have: %d", pattern.curVer, v)
	}

	for i, mig := range migs {
		m := mig.(*migration_mock)
		switch pattern.migs[i] {
		case notRun:
			if m.calledUp {
				t.Errorf("Migration %d, Want: notRun, Have upRun", i)
			}
			if m.calledDown {
				t.Errorf("Migration %d, Want: notRun, Have downRun", i)
			}
		case upRun:
			if !m.calledUp {
				t.Errorf("Migration %d, Not called Up", i)
			}
			if m.calledDown {
				t.Errorf("Migration %d, Want: upRun, Have downRun", i)
			}
		case downRun:
			if !m.calledDown {
				t.Errorf("Migration %d, Not called Down", i)
			}
			if m.calledUp {
				t.Errorf("Migration %d, Want: runDown, Have upRun", i)
			}
		}
	}
}

func setup() (gomig.BackendAdapter, []gomig.Migration, gomig.Runner) {
	backend := &backend_mock{}
	migs := make([]gomig.Migration, 10)
	for i, _ := range migs {
		migs[i] = &migration_mock{index: i}
	}
	r := gomig.NewRunner(migs, backend)
	return backend, migs, *r
}

func TestNoBackend(t *testing.T) {
	r := gomig.NewRunner(nil, nil)
	t.Run("Reset fails if no backend set", func(t *testing.T) {
		if err := r.Reset(); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})

	t.Run("ForceVer fails if no backend set", func(t *testing.T) {
		if err := r.ForceVer(0); !errors.Is(err, gomig.ErrNoBackend) {
			t.Errorf("Want: %v, Have: %v", gomig.ErrNoBackend, err)
		}
	})
}
