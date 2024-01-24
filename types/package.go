package types

import (
	"github.com/gordonklaus/dt/types/internal/types"
)

type Package struct {
	Name, Doc string
	Types     []*TypeName
}

func (p *Package) Type(name string) *TypeName {
	for _, t := range p.Types {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type PackageID interface{ isPackageID() }

type PackageID_Current struct{}

func (*PackageID_Current) isPackageID() {}

func (l *Loader) packageFromData(p types.Package, namedTypes map[*NamedType]uint64) *Package {
	pkg := &Package{
		Name:  p.Name,
		Doc:   p.Doc,
		Types: make([]*TypeName, len(p.Types)),
	}
	for i, t := range p.Types {
		pkg.Types[i] = &TypeName{
			Name: t.Name,
			Doc:  t.Doc,
			Type: l.typeFromData(t.Type, namedTypes),
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
			Name: t.Name,
			Doc:  t.Doc,
			Type: l.typeToData(t.Type),
		}
	}
	return pkg
}

func packageIDFromData(p types.PackageId) PackageID {
	switch p.PackageId.(type) {
	case *types.PackageId_Current:
		return &PackageID_Current{}
	}
	panic("unreached")
}

func packageIDToData(p PackageID) types.PackageId {
	switch p.(type) {
	case *PackageID_Current:
		return types.PackageId{PackageId: &types.PackageId_Current{}}
	}
	panic("unreached")
}
