package data

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type UintValue struct {
	Type *types.UintType
	i    uint64
}

func NewUintValue(t *types.UintType) *UintValue {
	return &UintValue{Type: t}
}

func (i *UintValue) Get() uint64 { return i.i }

func (i *UintValue) Set(x uint64) {
	if x>>i.Type.Size > 0 {
		panic(fmt.Sprintf("%d overflows %d-bit uint", x, i.Type.Size))
	}
	i.i = x
}

func (i *UintValue) Write(b *bits.Buffer) {
	b.WriteVarUint(i.i)
}

func (i *UintValue) Read(b *bits.Buffer) error {
	return b.ReadVarUint(&i.i, int(i.Type.Size))
}

type IntValue struct {
	Type *types.IntType
	i    int64
}

func NewIntValue(t *types.IntType) *IntValue {
	return &IntValue{Type: t}
}

func (i *IntValue) Get() int64 { return i.i }

func (i *IntValue) Set(x int64) {
	if y := x >> i.Type.Size; y > 0 || y < -1 {
		panic(fmt.Sprintf("%d overflows %d-bit int", x, i.Type.Size))
	}
	i.i = x
}

func (i *IntValue) Write(b *bits.Buffer) {
	b.WriteVarInt(i.i)
}

func (i *IntValue) Read(b *bits.Buffer) error {
	return b.ReadVarInt(&i.i, int(i.Type.Size))
}
