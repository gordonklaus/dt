package data

import (
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
)

type BoolValue struct {
	Type *types.BoolType
	X    bool
}

func NewBoolValue(t *types.BoolType) *BoolValue {
	return &BoolValue{Type: t}
}

func (i *BoolValue) Write(e *bits.Encoder) {
	e.WriteBool(i.X)
}

func (i *BoolValue) Read(d *bits.Decoder) error {
	return d.ReadBool(&i.X)
}
