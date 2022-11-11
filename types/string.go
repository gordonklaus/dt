package types

import "github.com/gordonklaus/data/bits"

type StringType struct{}

func (*StringType) Write(b *bits.Buffer) {
	b.WriteSize(func() {})
}

func (*StringType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return nil })
}
