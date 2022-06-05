package data

import (
	"fmt"

	"github.com/gordonklaus/data/types"
)

type EnumValue struct {
	Type  *types.EnumType
	Elem  int
	Value Value
}

func NewEnumValue(t *types.EnumType) *EnumValue {
	return &EnumValue{
		Type:  t,
		Elem:  0,
		Value: NewValue(t.Elems[0].Type),
	}
}

func (e *Encoder) EncodeEnumValue(en *EnumValue) error {
	_ = e.writeBinary(en.Elem) || e.encodeValue(en.Value)
	return e.err
}

func (d *Decoder) DecodeEnumValue(e *EnumValue) error {
	if !d.readBinary(&e.Elem) {
		return d.err
	}
	if e.Elem < 0 || e.Elem >= len(e.Type.Elems) {
		return fmt.Errorf("enum index out of range: %d", e.Elem)
	}
	e.Value = NewValue(e.Type.Elems[e.Elem].Type)
	d.decodeValue(e.Value)
	return d.err
}
