package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type OptionValue struct {
	Type *types.OptionType
	Elem Value
}

func (o *OptionValue) Write(b *bits.Buffer) {
	b.WriteBool(o.Elem != nil)
	if o.Elem != nil {
		o.Elem.Write(b)
	}
}

func (o *OptionValue) Read(b *bits.Buffer) error {
	var ok bool
	if err := b.ReadBool(&ok); err != nil {
		return err
	} else if ok {
		o.Elem = NewValue(o.Type.Elem)
		return o.Elem.Read(b)
	}
	return nil
}
