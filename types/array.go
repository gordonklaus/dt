package types

import "github.com/gordonklaus/data/bits"

type ArrayType struct {
	Elem Type // not *EnumType or *StructType
}

func (a *ArrayType) Write(b *bits.Buffer) {
	b.WriteSize(func() { WriteType(b, a.Elem) })
}

func (a *ArrayType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return ReadType(b, &a.Elem) })
}
