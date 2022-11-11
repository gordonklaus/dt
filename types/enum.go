package types

import "github.com/gordonklaus/data/bits"

type EnumType struct {
	Elems []*EnumElemType
}

func (e *EnumType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(uint64(len(e.Elems)))
		for _, el := range e.Elems {
			el.Write(b)
		}
	})
}

func (e *EnumType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		len, err := b.ReadVarUint()
		if err != nil {
			return err
		}
		e.Elems = make([]*EnumElemType, len)
		for i := range e.Elems {
			e.Elems[i] = &EnumElemType{}
			if err := e.Elems[i].Read(b); err != nil {
				return err
			}
		}
		return nil
	})
}

type EnumElemType struct {
	Name string
	Type Type
}

func (e *EnumElemType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(e.Name)
		WriteType(b, e.Type)
	})
}

func (e *EnumElemType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&e.Name); err != nil {
			return err
		}
		return ReadType(b, &e.Type)
	})
}
