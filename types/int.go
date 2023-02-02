package types

import "github.com/gordonklaus/data/bits"

type UintType struct {
	Size uint64 // 0 < Size <= 64
}

func (i *UintType) Write(b *bits.Buffer) {
	b.WriteSize(func() { b.WriteVarUint(i.Size) })
}

func (i *UintType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return bits.ReadVarUint(b, &i.Size)
	})
}

type IntType struct {
	Size uint64 // 0 < Size <= 64
}

func (i *IntType) Write(b *bits.Buffer) {
	b.WriteSize(func() { b.WriteVarUint(i.Size) })
}

func (i *IntType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return bits.ReadVarUint(b, &i.Size)
	})
}
