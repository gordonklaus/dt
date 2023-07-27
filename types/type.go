package types

import (
	"github.com/gordonklaus/data/types/internal/types"
)

type Type interface {
	isType()
}

func (*BoolType) isType()   {}
func (*IntType) isType()    {}
func (*FloatType) isType()  {}
func (*StringType) isType() {}
func (*NamedType) isType()  {}
func (*OptionType) isType() {}
func (*ArrayType) isType()  {}
func (*MapType) isType()    {}
func (*EnumType) isType()   {}
func (*StructType) isType() {}

func typeFromData(t types.Type) Type {
	switch t := t.Type.(type) {
	case *types.Type_Bool:
		return &BoolType{}
	case *types.Type_Int:
		return &IntType{Unsigned: t.Unsigned}
	case *types.Type_Float:
		return &FloatType{Size: t.Size}
	case *types.Type_String:
		return &StringType{}
	case *types.Type_Named:
		return &NamedType{
			Package: packageIDFromData(t.Package),
			Name:    t.Name,
		}
	case *types.Type_Option:
		return &OptionType{Elem: typeFromData(t.Element)}
	case *types.Type_Array:
		return &ArrayType{Elem: typeFromData(t.Element)}
	case *types.Type_Map:
		return &MapType{
			Key:   typeFromData(t.Key),
			Value: typeFromData(t.Value),
		}
	case *types.Type_Enum:
		return enumTypeFromData(t)
	case *types.Type_Struct:
		return structTypeFromData(t)
	}
	panic("unreached")
}

func typeToData(t Type) types.Type {
	switch t := t.(type) {
	case *BoolType:
		return types.Type{Type: &types.Type_Bool{}}
	case *IntType:
		return types.Type{Type: &types.Type_Int{Unsigned: t.Unsigned}}
	case *FloatType:
		return types.Type{Type: &types.Type_Float{Size: t.Size}}
	case *StringType:
		return types.Type{Type: &types.Type_String{}}
	case *NamedType:
		return types.Type{Type: &types.Type_Named{
			Package: packageIDToData(t.Package),
			Name:    t.Name,
		}}
	case *OptionType:
		return types.Type{Type: &types.Type_Option{Element: typeToData(t.Elem)}}
	case *ArrayType:
		return types.Type{Type: &types.Type_Array{Element: typeToData(t.Elem)}}
	case *MapType:
		return types.Type{Type: &types.Type_Map{
			Key:   typeToData(t.Key),
			Value: typeToData(t.Value),
		}}
	case *EnumType:
		return types.Type{Type: enumTypeToData(t)}
	case *StructType:
		return types.Type{Type: structTypeToData(t)}
	}
	panic("unreached")
}
