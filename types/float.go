package types

import "github.com/gordonklaus/data/bits"

type Float64Type struct{}

func (*Float64Type) Write(b *bits.Buffer) {
	b.WriteSize(func() {})
}

func (*Float64Type) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return nil })
}

type Float32Type struct{}

func (*Float32Type) Write(b *bits.Buffer) {
	b.WriteSize(func() {})
}

func (*Float32Type) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return nil })
}
