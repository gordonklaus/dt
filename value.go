package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type Value interface {
	bits.ReadWriter
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
	case *types.NamedType:
		return NewValue(t.Type)
	}
	return nil
}
