package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type StringValue struct {
	Type *types.StringType
	X    string
}

func NewStringValue(t *types.StringType) *StringValue {
	return &StringValue{Type: t}
}

func (i *StringValue) Write(e *bits.Encoder) {
	e.WriteString(i.X)
}

func (i *StringValue) Read(d *bits.Decoder) error {
	return d.ReadString(&i.X)
}
