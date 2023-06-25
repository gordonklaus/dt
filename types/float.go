package types

import "github.com/gordonklaus/data/bits"

type FloatType struct {
	Size uint64 // 32 or 64
}

func (f *FloatType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(f.Size)
	})
}

func (f *FloatType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return bits.ReadVarUint(b, &f.Size)
	})
}
