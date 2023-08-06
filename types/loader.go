package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types/internal/types"
)

type Loader struct {
	storage  *Storage
	Packages map[PackageID]*Package
}

func NewLoader(storage *Storage) *Loader {
	return &Loader{
		storage:  storage,
		Packages: map[PackageID]*Package{},
	}
}

func (l *Loader) Load(id PackageID) (*Package, error) {
	if p, ok := l.Packages[id]; ok {
		return p, nil
	}

	buf, err := l.storage.Load(id)
	if err != nil {
		return nil, err
	}

	var pkg types.Package
	b := bits.NewReadBuffer(buf)
	if err := pkg.Read(b); err != nil {
		return nil, err
	}
	if b.Remaining() > 7 {
		return nil, fmt.Errorf("%d bits remaining after reading package", b.Remaining())
	}

	namedTypes := map[*NamedType]string{}

	p := packageFromData(pkg, namedTypes)

	if err := ValidatePackage(p); err != nil {
		return nil, err
	}

	l.Packages[id] = p

	for nt, name := range namedTypes {
		p, err := l.Load(nt.Package)
		if err != nil {
			return nil, err
		}
		tn := p.Type(name)
		if tn == nil {
			return nil, fmt.Errorf("package %s has no type %s", p.Name, name)
		}
		nt.TypeName = tn
	}

	return p, nil
}

func (l *Loader) Store(id PackageID) error {
	p, ok := l.Packages[id]
	if !ok {
		return fmt.Errorf("package %#v not yet loaded", id)
	}

	if err := ValidatePackage(p); err != nil {
		return err
	}

	pkg := packageToData(p)

	b := bits.NewBuffer()
	pkg.Write(b)
	return l.storage.Store(id, b.Bytes())
}
