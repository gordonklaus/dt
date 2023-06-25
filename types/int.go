package types

import "github.com/gordonklaus/data/bits"

type IntType struct {
	Size     uint64 // 0 < Size <= 64
	Unsigned bool
}

func (i *IntType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(i.Size)
		b.WriteBool(i.Unsigned)
	})
}

func (i *IntType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := bits.ReadVarUint(b, &i.Size); err != nil {
			return nil
		}
		return b.ReadBool(&i.Unsigned)
	})
}
