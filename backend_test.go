package gomig_test

type adapter_mock struct {
	currentVer int
	seq        []int
}

func (a *adapter_mock) GetVersion() (int, error) {
	return a.currentVer, nil
}

func (a *adapter_mock) SaveVersion(ver int, messages []string) error {
	a.seq = append(a.seq, ver)
	a.currentVer = ver
	return nil
}

func (a *adapter_mock) reset() {
	a.currentVer = 0
	a.seq = []int{}
}
