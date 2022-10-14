package data

import (
	"fmt"

	"github.com/gordonklaus/data/types"
)

type Value interface {
	isValue()
}

func (*OptionValue) isValue() {}
func (*BasicValue) isValue()  {}
func (*ArrayValue) isValue()  {}
func (*StructValue) isValue() {}
func (*EnumValue) isValue()   {}

func NewValue(t types.Type) Value {
	switch t := t.(type) {
	case *types.BasicType:
		return NewBasicValue(t)
	case *types.ArrayType:
		return NewArrayValue(t)
	case *types.StructType:
		return NewStructValue(t)
	case *types.EnumType:
		return NewEnumValue(t)
	case *types.NamedType:
		return NewValue(t.Type)
	}
	return nil
}

func (e *Encoder) EncodeValue(v Value) error {
	switch v := v.(type) {
	case *BasicValue:
		return e.EncodeBasicValue(v)
	case *ArrayValue:
		return e.EncodeArrayValue(v)
	case *StructValue:
		return e.EncodeStructValue(v)
	case *EnumValue:
		return e.EncodeEnumValue(v)
	}
	panic(fmt.Sprintf("invalid Value type %T", v))
}

func (d *Decoder) DecodeValue(v Value) error {
	switch v := v.(type) {
	case *BasicValue:
		return d.DecodeBasicValue(v)
	case *ArrayValue:
		return d.DecodeArrayValue(v)
	case *StructValue:
		return d.DecodeStructValue(v)
	case *EnumValue:
		return d.DecodeEnumValue(v)
	}
	panic(fmt.Sprintf("invalid Value type %T", v))
}
