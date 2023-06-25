package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type ArrayValue struct {
	Type  *types.ArrayType
	Elems []Value
}

func NewArrayValue(t *types.ArrayType) *ArrayValue {
	return &ArrayValue{
		Type: t,
	}
}

func (a *ArrayValue) Write(b *bits.Buffer) {
	b.WriteVarUint(uint64(len(a.Elems)))
	for _, v := range a.Elems {
		v.Write(b)
	}
}

func (a *ArrayValue) Read(b *bits.Buffer) error {
	var len uint64
	if err := b.ReadVarUint(&len); err != nil {
		return err
	}
	a.Elems = make([]Value, len)
	for i := range a.Elems {
		a.Elems[i] = NewValue(a.Type.Elem)
		if err := a.Elems[i].Read(b); err != nil {
			return err
		}
	}
	return nil
}
