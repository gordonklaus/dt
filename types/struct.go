package types

import "github.com/gordonklaus/dt/types/internal/types"

type StructType struct {
	Fields []*StructFieldType
}

type StructFieldType struct {
	ID        uint64
	Name, Doc string
	Type      Type // not *EnumType or *StructType
}

func (l *Loader) structTypeFromData(t *types.Type_Struct, namedIDs map[*NamedType]uint64) *StructType {
	typ := &StructType{Fields: make([]*StructFieldType, len(t.Fields))}
	for i, f := range t.Fields {
		typ.Fields[i] = &StructFieldType{
			ID:   f.ID,
			Name: f.Name,
			Doc:  f.Doc,
			Type: l.typeFromData(f.Type, namedIDs),
		}
	}
	return typ
}

func (l *Loader) structTypeToData(t *StructType) *types.Type_Struct {
	typ := &types.Type_Struct{Fields: make([]types.StructField, len(t.Fields))}
	for i, f := range t.Fields {
		typ.Fields[i] = types.StructField{
			ID:   f.ID,
			Name: f.Name,
			Doc:  f.Doc,
			Type: l.typeToData(f.Type),
		}
	}
	return typ
}
