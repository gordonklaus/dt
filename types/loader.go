package types

import (
	"bytes"
	"fmt"

	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types/internal/types"
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
	bb := bytes.NewBuffer(buf)
	if err := bits.Read(bb, &pkg); err != nil {
		return nil, err
	}
	if bb.Len() > 0 {
		return nil, fmt.Errorf("%d bytes remaining after reading package", bb.Len())
	}

	ldr := packageLoader{l, nil, map[*NamedType]uint64{}}
	p := ldr.packageFromData(pkg)
	if err := ValidatePackage(p); err != nil {
		return nil, err
	}
	l.Packages[id] = p
	for nt, id := range ldr.namedIDs {
		p, err := l.Load(nt.Package)
		if err != nil {
			return nil, err
		}
		nt.TypeName = p.TypesByID[id]
		if nt.TypeName == nil {
			return nil, fmt.Errorf("package %s has no type with ID %d", p.Name, id)
		}
	}

	return p, nil
}

type packageLoader struct {
	*Loader
	pkg      *Package
	namedIDs map[*NamedType]uint64
}

func (l *Loader) Store(id PackageID) error {
	p, ok := l.Packages[id]
	if !ok {
		return fmt.Errorf("package %#v not yet loaded", id)
	}

	if err := ValidatePackage(p); err != nil {
		return err
	}

	pkg := l.packageToData(p)

	enc := bits.NewEncoder()
	pkg.Write(enc)
	return l.storage.Store(id, enc.Bytes())
}
