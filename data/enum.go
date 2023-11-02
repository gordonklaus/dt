package data

import (
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
)

type EnumValue struct {
	Type  *types.EnumType
	Elem  uint64
	Value Value
}

func NewEnumValue(t *types.EnumType) *EnumValue {
	return &EnumValue{
		Type:  t,
		Elem:  0,
		Value: NewValue(t.Elems[0].Type),
	}
}

func (e *EnumValue) Write(enc *bits.Encoder) {
	enc.WriteVarUint_4bit(e.Elem)
	e.Value.Write(enc)
}

func (e *EnumValue) Read(d *bits.Decoder) error {
	if err := d.ReadVarUint_4bit(&e.Elem); err != nil {
		return err
	}
	if e.Elem < uint64(len(e.Type.Elems)) {
		e.Value = NewValue(e.Type.Elems[e.Elem].Type)
	} else {
		e.Value = NewValue(&UnknownEnumElementType)
	}
	return e.Value.Read(d)
}

var UnknownEnumElementType types.StructType
