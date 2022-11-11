package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type OptionValue struct {
	*types.OptionType
	Value Value
}

func (o *OptionValue) Write(b *bits.Buffer) {
	b.WriteBool(o.Value != nil)
	if o.Value != nil {
		o.Value.Write(b)
	}
}

func (o *OptionValue) Read(b *bits.Buffer) error {
	if ok, err := b.ReadBool(); err != nil {
		return err
	} else if ok {
		o.Value = NewValue(o.ValueType)
		return o.Value.Read(b)
	}
	return nil
}
