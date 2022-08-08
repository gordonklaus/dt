package data

import (
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

func (e *Encoder) EncodeArrayValue(a *ArrayValue) error {
	if !e.writeBinary(uint(len(a.Elems))) {
		return e.err
	}
	for _, v := range a.Elems {
		if !e.encodeValue(v) {
			return e.err
		}
	}
	return e.err
}

func (d *Decoder) DecodeArrayValue(a *ArrayValue) error {
	var len uint
	if !d.readBinary(&len) {
		return d.err
	}
	a.Elems = make([]Value, len)
	for i := range a.Elems {
		a.Elems[i] = NewValue(a.Type.Elem)
		if !d.decodeValue(a.Elems[i]) {
			return d.err
		}
	}
	return d.err
}
