package types

import "github.com/gordonklaus/dt/types/internal/types"

type TypeName struct {
	ID        uint64
	Name, Doc string
	Type      Type // *EnumType or *StructType
}

func (l *Loader) typeNameFromData(t types.TypeName, namedIDs map[*NamedType]uint64) *TypeName {
	return &TypeName{
		ID:   t.ID,
		Name: t.Name,
		Doc:  t.Doc,
		Type: l.typeFromData(t.Type, namedIDs),
	}
}

func (l *Loader) typeNameToData(t *TypeName) *types.TypeName {
	return &types.TypeName{
		ID:   t.ID,
		Name: t.Name,
		Doc:  t.Doc,
		Type: l.typeToData(t.Type),
	}
}
