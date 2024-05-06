package types

import (
	"github.com/gordonklaus/dt/types/internal/types"
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

func (l *packageLoader) typeFromData(t types.Type, parent any) Type {
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
		nt := &NamedType{Package: packageIDFromData(t.Package)}
		l.namedIDs[nt] = t.ID
		return nt
	case *types.Type_Option:
		return &OptionType{Elem: l.typeFromData(t.Element, nil)}
	case *types.Type_Array:
		return &ArrayType{Elem: l.typeFromData(t.Element, nil)}
	case *types.Type_Map:
		return &MapType{
			Key:   l.typeFromData(t.Key, nil),
			Value: l.typeFromData(t.Value, nil),
		}
	case *types.Type_Enum:
		return l.enumTypeFromData(t, parent)
	case *types.Type_Struct:
		return l.structTypeFromData(t, parent)
	}
	panic("unreached")
}

func (l *Loader) typeToData(t Type) types.Type {
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
			ID:      t.TypeName.ID,
		}}
	case *OptionType:
		return types.Type{Type: &types.Type_Option{Element: l.typeToData(t.Elem)}}
	case *ArrayType:
		return types.Type{Type: &types.Type_Array{Element: l.typeToData(t.Elem)}}
	case *MapType:
		return types.Type{Type: &types.Type_Map{
			Key:   l.typeToData(t.Key),
			Value: l.typeToData(t.Value),
		}}
	case *EnumType:
		return types.Type{Type: l.enumTypeToData(t)}
	case *StructType:
		return types.Type{Type: l.structTypeToData(t)}
	}
	panic("unreached")
}
