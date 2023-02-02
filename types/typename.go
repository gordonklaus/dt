package types

import "github.com/gordonklaus/data/bits"

type TypeName struct {
	Name, Doc string
	Type      Type // *EnumType or *StructType
}

func (t *TypeName) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(t.Name)
		b.WriteString(t.Doc)
		WriteType(b, t.Type)
	})
}

func (t *TypeName) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&t.Name); err != nil {
			return err
		}
		if err := b.ReadString(&t.Doc); err != nil {
			return err
		}
		return ReadType(b, &t.Type)
	})
}
