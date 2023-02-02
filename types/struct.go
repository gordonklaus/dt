package types

import "github.com/gordonklaus/data/bits"

type StructType struct {
	Fields []*StructFieldType
}

func (s *StructType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(uint64(len(s.Fields)))
		for _, el := range s.Fields {
			el.Write(b)
		}
	})
}

func (s *StructType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		var len uint64
		if err := bits.ReadVarUint(b, &len); err != nil {
			return err
		}
		s.Fields = make([]*StructFieldType, len)
		for i := range s.Fields {
			s.Fields[i] = &StructFieldType{}
			if err := s.Fields[i].Read(b); err != nil {
				return err
			}
		}
		return nil
	})
}

type StructFieldType struct {
	Name, Doc string
	Type      Type // not *EnumType or *StructType
}

func (s *StructFieldType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(s.Name)
		b.WriteString(s.Doc)
		WriteType(b, s.Type)
	})
}

func (s *StructFieldType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&s.Name); err != nil {
			return err
		}
		if err := b.ReadString(&s.Doc); err != nil {
			return err
		}
		return ReadType(b, &s.Type)
	})
}
