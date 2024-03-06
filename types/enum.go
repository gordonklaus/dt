package types

import "github.com/gordonklaus/dt/types/internal/types"

type EnumType struct {
	Elems []*EnumElemType
}

type EnumElemType struct {
	ID        uint64
	Name, Doc string
	Type      Type // *StructType
}

func (l *Loader) enumTypeFromData(t *types.Type_Enum, namedIDs map[*NamedType]uint64) *EnumType {
	typ := &EnumType{Elems: make([]*EnumElemType, len(t.Elements))}
	for i, e := range t.Elements {
		typ.Elems[i] = &EnumElemType{
			ID:   e.ID,
			Name: e.Name,
			Doc:  e.Doc,
			Type: l.typeFromData(e.Type, namedIDs),
		}
	}
	return typ
}

func (l *Loader) enumTypeToData(t *EnumType) *types.Type_Enum {
	typ := &types.Type_Enum{Elements: make([]types.EnumElement, len(t.Elems))}
	for i, e := range t.Elems {
		typ.Elements[i] = types.EnumElement{
			ID:   e.ID,
			Name: e.Name,
			Doc:  e.Doc,
			Type: l.typeToData(e.Type),
		}
	}
	return typ
}
