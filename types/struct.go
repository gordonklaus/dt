package types

import "github.com/gordonklaus/dt/types/internal/types"

type StructType struct {
	Fields []*StructFieldType
}

type StructFieldType = TypeName

func (l *packageLoader) structTypeFromData(t *types.Type_Struct, parent any) *StructType {
	typ := &StructType{Fields: make([]*StructFieldType, len(t.Fields))}
	for i, f := range t.Fields {
		typ.Fields[i] = l.typeNameFromData(f, parent)
	}
	return typ
}

func (l *Loader) structTypeToData(t *StructType) *types.Type_Struct {
	typ := &types.Type_Struct{Fields: make([]types.TypeName, len(t.Fields))}
	for i, f := range t.Fields {
		typ.Fields[i] = *l.typeNameToData(f)
	}
	return typ
}
