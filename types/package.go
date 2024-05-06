package types

import (
	"github.com/gordonklaus/dt/types/internal/types"
)

type Package struct {
	Name, Doc string
	Types     []*TypeName
	TypesByID map[uint64]*TypeName
}

type PackageID interface{ isPackageID() }

type PackageID_Current struct{}

func (PackageID_Current) isPackageID() {}

func (l *packageLoader) packageFromData(p types.Package) *Package {
	l.pkg = &Package{
		Name:      p.Name,
		Doc:       p.Doc,
		Types:     make([]*TypeName, len(p.Types)),
		TypesByID: make(map[uint64]*TypeName, len(p.Types)),
	}
	for i, t := range p.Types {
		l.pkg.Types[i] = l.typeNameFromData(t, l.pkg)
	}
	return l.pkg
}

func (l *Loader) packageToData(p *Package) types.Package {
	pkg := types.Package{
		Name:  p.Name,
		Doc:   p.Doc,
		Types: make([]types.TypeName, len(p.Types)),
	}
	for i, t := range p.Types {
		pkg.Types[i] = *l.typeNameToData(t)
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
