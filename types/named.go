package types

import "github.com/gordonklaus/data/bits"

type NamedType struct {
	Name string
	Type Type
}

func (n *NamedType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(n.Name)
		WriteType(b, n.Type)
	})
}

func (n *NamedType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&n.Name); err != nil {
			return err
		}
		return ReadType(b, &n.Type)
	})
}
