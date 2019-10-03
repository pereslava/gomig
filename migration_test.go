package gomig_test

import (
	"errors"

	"github.com/pereslava/gomig"
)

type migMock struct {
	index      int
	calledUp   int
	calledDown int
	fail       error
}

func (mig *migMock) Up() (messages []string, err error) {
	if mig.fail != nil {
		return nil, mig.fail
	}
	mig.calledUp++
	return []string{}, nil
}

func (mig *migMock) Down() (messages []string, err error) {
	if mig.fail != nil {
		return nil, errors.New("Some error")
	}
	mig.calledDown++
	return []string{}, nil
}

func (mig *migMock) reset(index int) {
	mig.index = index
	mig.calledUp = 0
	mig.calledDown = 0
	mig.fail = nil
}

func resetMigrations(migs []gomig.Migration) {
	for i, m := range migs {
		m.(*migMock).reset(i)
	}
}
