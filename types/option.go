package types

import "github.com/gordonklaus/data/bits"

type OptionType struct {
	Elem Type // not *EnumType or *StructType
}

func (o *OptionType) Write(b *bits.Buffer) {
	b.WriteSize(func() { WriteType(b, o.Elem) })
}

func (o *OptionType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return ReadType(b, &o.Elem) })
}
