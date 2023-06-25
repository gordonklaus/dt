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
	i.i = x
}

func (i *IntValue) SetUint(x uint64) {
	if !i.Type.Unsigned {
		panic(fmt.Sprintf("int is signed"))
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
		err := b.ReadVarUint(&x)
		i.i = int64(x)
		return err
	} else {
		return b.ReadVarInt(&i.i)
	}
}
