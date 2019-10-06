package gomig_test

import (
	"errors"

	"github.com/pereslava/gomig"
)

type migration_mock struct {
	index      int
	calledUp   bool
	calledDown bool
	fail       error
}

func (mig *migration_mock) Up() (messages []string, err error) {
	if mig.fail != nil {
		return nil, mig.fail
	}
	mig.calledUp = true
	return []string{}, nil
}

func (mig *migration_mock) Down() (messages []string, err error) {
	if mig.fail != nil {
		return nil, errors.New("Some error")
	}
	mig.calledDown = true
	return []string{}, nil
}

func (mig *migration_mock) reset(index int) {
	mig.index = index
	mig.calledUp = false
	mig.calledDown = false
	mig.fail = nil
}

func resetMigrations(migs []gomig.Migration) {
	for i, m := range migs {
		m.(*migration_mock).reset(i)
	}
}
