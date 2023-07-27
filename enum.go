package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
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

func (e *EnumValue) Write(b *bits.Buffer) {
	b.WriteVarUint_4bit(e.Elem)
	e.Value.Write(b)
}

func (e *EnumValue) Read(b *bits.Buffer) error {
	if err := b.ReadVarUint_4bit(&e.Elem); err != nil {
		return err
	}
	if e.Elem < uint64(len(e.Type.Elems)) {
		e.Value = NewValue(e.Type.Elems[e.Elem].Type)
	} else {
		e.Value = NewValue(&UnknownEnumElementType)
	}
	return e.Value.Read(b)
}

var UnknownEnumElementType types.StructType
