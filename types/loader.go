package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
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

	var p Package
	b := bits.NewReadBuffer(buf)
	if err := p.Read(b); err != nil {
		return nil, err
	}
	if b.Remaining() > 7 {
		return nil, fmt.Errorf("%d bits remaining after reading package", b.Remaining())
	}

	l.Packages[id] = &p

	for _, t := range p.Types {
		if err := l.setNamedTypes(t.Type); err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (l *Loader) setNamedTypes(t Type) error {
	switch t := t.(type) {
	case *EnumType:
		for _, e := range t.Elems {
			if err := l.setNamedTypes(&e.Type); err != nil {
				return err
			}
		}
	case *StructType:
		for _, f := range t.Fields {
			if err := l.setNamedTypes(f.Type); err != nil {
				return err
			}
		}
	case *ArrayType:
		return l.setNamedTypes(t.Elem)
	case *OptionType:
		return l.setNamedTypes(t.Elem)
	case *NamedType:
		p, err := l.Load(t.Package)
		if err != nil {
			return err
		}
		tn := p.Type(t.Name)
		if tn == nil {
			return fmt.Errorf("package %s has no type %s", p.Name, t.Name)
		}
		t.Type = tn.Type
	}
	return nil
}

func (l *Loader) Store(id PackageID) error {
	p, ok := l.Packages[id]
	if !ok {
		return fmt.Errorf("package %#v not yet loaded", id)
	}

	b := bits.NewBuffer()
	p.Write(b)
	return l.storage.Store(id, b.Bytes())
}
