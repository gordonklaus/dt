package data

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type IntValue struct {
	Type *types.IntType
	i    int64
}

func NewIntValue(t *types.IntType) *IntValue {
	return &IntValue{Type: t}
}

func (i *IntValue) GetInt() int64 {
	if i.Type.Unsigned {
		panic(fmt.Sprintf("int is unsigned"))
	}
	return i.i
}
func (i *IntValue) GetUint() uint64 {
	if !i.Type.Unsigned {
		panic(fmt.Sprintf("int is signed"))
	}
	return uint64(i.i)
}

func (i *IntValue) SetInt(x int64) {
	if i.Type.Unsigned {
		panic(fmt.Sprintf("int is unsigned"))
	}
	if y := x >> i.Type.Size; y > 0 || y < -1 {
		panic(fmt.Sprintf("%d overflows %d-bit int", x, i.Type.Size))
	}
	i.i = x
}

func (i *IntValue) SetUint(x uint64) {
	if !i.Type.Unsigned {
		panic(fmt.Sprintf("int is signed"))
	}
	if x>>i.Type.Size > 0 {
		panic(fmt.Sprintf("%d overflows %d-bit uint", x, i.Type.Size))
	}
	i.i = int64(x)
}

func (i *IntValue) Write(b *bits.Buffer) {
	if i.Type.Unsigned {
		b.WriteVarUint(uint64(i.i))
	} else {
		b.WriteVarInt(i.i)
	}
}

func (i *IntValue) Read(b *bits.Buffer) error {
	if i.Type.Unsigned {
		var x uint64
		err := b.ReadVarUint(&x, int(i.Type.Size))
		i.i = int64(x)
		return err
	} else {
		return b.ReadVarInt(&i.i, int(i.Type.Size))
	}
}
