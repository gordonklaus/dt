package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type BoolValue struct {
	Type *types.BoolType
	X    bool
}

func NewBoolValue(t *types.BoolType) *BoolValue {
	return &BoolValue{Type: t}
}

func (i *BoolValue) Write(b *bits.Buffer) {
	b.WriteBool(i.X)
}

func (i *BoolValue) Read(b *bits.Buffer) error {
	return b.ReadBool(&i.X)
}
