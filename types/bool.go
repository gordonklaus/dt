package types

import "github.com/gordonklaus/data/bits"

type BoolType struct{}

func (*BoolType) Write(b *bits.Buffer) {
	b.WriteSize(func() {})
}

func (*BoolType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return nil })
}
