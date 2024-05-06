package types

import "github.com/gordonklaus/dt/types/internal/types"

type TypeName struct {
	ID        uint64
	Name, Doc string
	Type      Type // *EnumType or *StructType
	Parent    any  // *Package or *TypeName or *EnumElemType
}

func (l *packageLoader) typeNameFromData(t types.TypeName, parent any) *TypeName {
	n := &TypeName{
		ID:     t.ID,
		Name:   t.Name,
		Doc:    t.Doc,
		Parent: parent,
	}
	n.Type = l.typeFromData(t.Type, n)
	l.pkg.TypesByID[n.ID] = n
	return n
}

func (l *Loader) typeNameToData(t *TypeName) *types.TypeName {
	return &types.TypeName{
		ID:   t.ID,
		Name: t.Name,
		Doc:  t.Doc,
		Type: l.typeToData(t.Type),
	}
}
