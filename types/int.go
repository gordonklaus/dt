package types

import "github.com/gordonklaus/data/bits"

type IntType struct {
	Unsigned bool
}

func (i *IntType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteBool(i.Unsigned)
	})
}

func (i *IntType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return b.ReadBool(&i.Unsigned)
	})
}
