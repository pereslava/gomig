package gomig_test

type backend_mock struct {
	currentVer uint
	seq        []uint
}

func (a *backend_mock) GetVersion() (uint, error) {
	return a.currentVer, nil
}

func (a *backend_mock) SaveVersion(ver uint, messages []string) error {
	a.seq = append(a.seq, ver)
	a.currentVer = ver
	return nil
}

func (a *backend_mock) Reset() error {
	a.currentVer = 0
	a.seq = []uint{}
	return nil
}
