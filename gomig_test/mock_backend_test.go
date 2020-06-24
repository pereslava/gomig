package gomig_test

import "context"

type backend_mock struct {
	currentVer uint
	seq        []uint
}

func (a *backend_mock) GetVersion(ctx context.Context) (uint, error) {
	return a.currentVer, nil
}

func (a *backend_mock) SaveVersion(ctx context.Context, ver uint, name string, messages []string) error {
	a.seq = append(a.seq, ver)
	a.currentVer = ver
	return nil
}

func (a *backend_mock) Reset(ctx context.Context) error {
	a.currentVer = 0
	a.seq = []uint{}
	return nil
}
