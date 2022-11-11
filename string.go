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

func (i *StringValue) Write(b *bits.Buffer) {
	b.WriteString(i.X)
}

func (i *StringValue) Read(b *bits.Buffer) error {
	return b.ReadString(&i.X)
}
