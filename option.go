package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type OptionValue struct {
	Type *types.OptionType
	Elem Value
}

func (o *OptionValue) Write(e *bits.Encoder) {
	e.WriteBool(o.Elem != nil)
	if o.Elem != nil {
		o.Elem.Write(e)
	}
}

func (o *OptionValue) Read(d *bits.Decoder) error {
	var ok bool
	if err := d.ReadBool(&ok); err != nil {
		return err
	} else if ok {
		o.Elem = NewValue(o.Type.Elem)
		return o.Elem.Read(d)
	}
	return nil
}
