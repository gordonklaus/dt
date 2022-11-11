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
		len, err := b.ReadVarUint()
		if err != nil {
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
	Name string
	Type Type
}

func (s *StructFieldType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(s.Name)
		WriteType(b, s.Type)
	})
}

func (s *StructFieldType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&s.Name); err != nil {
			return err
		}
		return ReadType(b, &s.Type)
	})
}
