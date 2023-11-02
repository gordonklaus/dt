package data

import (
	"fmt"

	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
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

func (i *IntValue) Write(e *bits.Encoder) {
	if i.Type.Unsigned {
		e.WriteVarUint(uint64(i.i))
	} else {
		e.WriteVarInt(i.i)
	}
}

func (i *IntValue) Read(d *bits.Decoder) error {
	if i.Type.Unsigned {
		var x uint64
		err := d.ReadVarUint(&x)
		i.i = int64(x)
		return err
	} else {
		return d.ReadVarInt(&i.i)
	}
}
