package types

import (
	"github.com/gordonklaus/dt/types/internal/types"
)

type Package struct {
	Name, Doc string
	Types     []*TypeName
}

func (p *Package) Type(id uint64) *TypeName {
	for _, t := range p.Types {
		if t.ID == id {
			return t
		}
	}
	return nil
}

type PackageID interface{ isPackageID() }

type PackageID_Current struct{}

func (PackageID_Current) isPackageID() {}

func (l *Loader) packageFromData(p types.Package, namedIDs map[*NamedType]uint64) *Package {
	pkg := &Package{
		Name:  p.Name,
		Doc:   p.Doc,
		Types: make([]*TypeName, len(p.Types)),
	}
	for i, t := range p.Types {
		pkg.Types[i] = &TypeName{
			ID:   t.ID,
			Name: t.Name,
			Doc:  t.Doc,
			Type: l.typeFromData(t.Type, namedIDs),
		}
	}
	return pkg
}

func (l *Loader) packageToData(p *Package) types.Package {
	pkg := types.Package{
		Name:  p.Name,
		Doc:   p.Doc,
		Types: make([]types.TypeName, len(p.Types)),
	}
	for i, t := range p.Types {
		pkg.Types[i] = types.TypeName{
			ID:   t.ID,
			Name: t.Name,
			Doc:  t.Doc,
			Type: l.typeToData(t.Type),
		}
	}
	return pkg
}

func packageIDFromData(p types.PackageID) PackageID {
	switch p.PackageID.(type) {
	case *types.PackageID_Current:
		return PackageID_Current{}
	}
	panic("unreached")
}

func packageIDToData(p PackageID) types.PackageID {
	switch p.(type) {
	case PackageID_Current:
		return types.PackageID{PackageID: &types.PackageID_Current{}}
	}
	panic("unreached")
}
