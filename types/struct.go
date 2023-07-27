package types

import "github.com/gordonklaus/data/types/internal/types"

type StructType struct {
	Fields []*StructFieldType
}

type StructFieldType struct {
	Name, Doc string
	Type      Type // not *EnumType or *StructType
}

func structTypeFromData(t *types.Type_Struct) *StructType {
	typ := &StructType{Fields: make([]*StructFieldType, len(t.Fields))}
	for i, f := range t.Fields {
		typ.Fields[i] = &StructFieldType{
			Name: f.Name,
			Doc:  f.Doc,
			Type: typeFromData(f.Type),
		}
	}
	return typ
}

func structTypeToData(t *StructType) *types.Type_Struct {
	typ := &types.Type_Struct{Fields: make([]types.StructField, len(t.Fields))}
	for i, e := range t.Fields {
		typ.Fields[i] = types.StructField{
			Name: e.Name,
			Doc:  e.Doc,
			Type: typeToData(e.Type),
		}
	}
	return typ
}
