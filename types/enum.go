package types

import "github.com/gordonklaus/dt/types/internal/types"

type EnumType struct {
	Elems []*EnumElemType
}

type EnumElemType = TypeName

func (l *Loader) enumTypeFromData(t *types.Type_Enum, namedIDs map[*NamedType]uint64) *EnumType {
	typ := &EnumType{Elems: make([]*EnumElemType, len(t.Elements))}
	for i, e := range t.Elements {
		typ.Elems[i] = l.typeNameFromData(e, namedIDs)
	}
	return typ
}

func (l *Loader) enumTypeToData(t *EnumType) *types.Type_Enum {
	typ := &types.Type_Enum{Elements: make([]types.TypeName, len(t.Elems))}
	for i, e := range t.Elems {
		typ.Elements[i] = *l.typeNameToData(e)
	}
	return typ
}
