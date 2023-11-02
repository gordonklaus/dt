package data

import (
	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
)

type StructValue struct {
	Type             *types.StructType
	Fields           []*StructFieldValue
	HasUnknownFields bool
}

type StructFieldValue struct {
	Type  *types.StructFieldType
	Value Value
}

func NewStructValue(t *types.StructType) *StructValue {
	s := &StructValue{
		Type:   t,
		Fields: make([]*StructFieldValue, len(t.Fields)),
	}
	for i, f := range t.Fields {
		s.Fields[i] = &StructFieldValue{
			Type:  f,
			Value: NewValue(f.Type),
		}
	}
	return s
}

func (s *StructValue) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		for _, f := range s.Fields {
			f.Value.Write(e)
		}
	})
}

func (s *StructValue) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		for _, f := range s.Fields {
			if d.Remaining() == 0 {
				return nil
			}
			if err := f.Value.Read(d); err != nil {
				return err
			}
		}
		if d.Remaining() > 0 {
			s.HasUnknownFields = true
		}
		return nil
	})
}
