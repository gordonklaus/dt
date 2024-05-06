package types

import "github.com/gordonklaus/dt/types/internal/types"

type EnumType struct {
	Elems []*EnumElemType
}

type EnumElemType = TypeName

func (l *packageLoader) enumTypeFromData(t *types.Type_Enum, parent any) *EnumType {
	typ := &EnumType{Elems: make([]*EnumElemType, len(t.Elements))}
	for i, e := range t.Elements {
		typ.Elems[i] = l.typeNameFromData(e, parent)
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
