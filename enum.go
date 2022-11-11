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
	if s, ok := e.Value.(*StructValue); ok {
		// Struct values already include their size.
		s.Write(b)
	} else {
		b.WriteSize(func() { e.Value.Write(b) })
	}
}

func (e *EnumValue) Read(b *bits.Buffer) error {
	var err error
	if e.Elem, err = b.ReadVarUint_4bit(); err != nil {
		return err
	}
	if e.Elem < uint64(len(e.Type.Elems)) {
		e.Value = NewValue(e.Type.Elems[e.Elem].Type)
	} else {
		e.Value = &UnknownEnumElement{}
	}
	if s, ok := e.Value.(*StructValue); ok {
		return s.Read(b)
	}
	return b.ReadSize(func() error { return e.Value.Read(b) })
}

type UnknownEnumElement struct{}

func (*UnknownEnumElement) Write(b *bits.Buffer)      {}
func (*UnknownEnumElement) Read(b *bits.Buffer) error { return nil }
