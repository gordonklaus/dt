package data

import (
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
)

type Value interface {
	bits.Value
}

func NewValue(t types.Type) Value {
	switch t := t.(type) {
	case *types.BoolType:
		return NewBoolValue(t)
	case *types.IntType:
		return NewIntValue(t)
	case *types.FloatType:
		return NewFloatValue(t)
	case *types.StringType:
		return NewStringValue(t)
	case *types.StructType:
		return NewStructValue(t)
	case *types.EnumType:
		return NewEnumValue(t)
	case *types.ArrayType:
		return NewArrayValue(t)
	case *types.MapType:
		return NewMapValue(t)
	case *types.NamedType:
		return NewValue(t.TypeName.Type)
	}
	return nil
}
