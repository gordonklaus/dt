package data

import (
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
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

func (a *ArrayValue) Write(e *bits.Encoder) {
	e.WriteVarUint(uint64(len(a.Elems)))
	for _, v := range a.Elems {
		v.Write(e)
	}
}

func (a *ArrayValue) Read(d *bits.Decoder) error {
	var len uint64
	if err := d.ReadVarUint(&len); err != nil {
		return err
	}
	a.Elems = make([]Value, len)
	for i := range a.Elems {
		a.Elems[i] = NewValue(a.Type.Elem)
		if err := a.Elems[i].Read(d); err != nil {
			return err
		}
	}
	return nil
}
