package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type StructValue struct {
	Type             *types.StructType
	Fields           []*StructFieldValue
	HasUnknownFields bool
}

type StructFieldValue struct {
	*types.StructFieldType
	Value Value
}

func NewStructValue(t *types.StructType) *StructValue {
	s := &StructValue{
		Type:   t,
		Fields: make([]*StructFieldValue, len(t.Fields)),
	}
	for i, f := range t.Fields {
		s.Fields[i] = &StructFieldValue{
			StructFieldType: f,
			Value:           NewValue(f.Type),
		}
	}
	return s
}

func (s *StructValue) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		for _, f := range s.Fields {
			f.Value.Write(b)
		}
	})
}

func (s *StructValue) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		for _, f := range s.Fields {
			if b.Remaining() == 0 {
				return nil
			}
			if err := f.Value.Read(b); err != nil {
				return err
			}
		}
		if b.Remaining() > 0 {
			s.HasUnknownFields = true
		}
		return nil
	})
}
