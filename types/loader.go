package types

type Loader struct {
	Packages []*Package
}

func (l *Loader) Load(path string) (*Package, error) {
	return nil, nil
}
