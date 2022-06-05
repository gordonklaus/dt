package data

import "github.com/gordonklaus/data/types"

type StructValue struct {
	Fields []*StructFieldValue
}

type StructFieldValue struct {
	Name  string
	Value Value
}

func NewStructValue(t *types.StructType) *StructValue {
	s := &StructValue{
		Fields: make([]*StructFieldValue, len(t.Fields)),
	}
	for i, f := range t.Fields {
		s.Fields[i] = &StructFieldValue{
			Name:  f.Name,
			Value: NewValue(f.Type),
		}
	}
	return s
}

func (e *Encoder) EncodeStructValue(s *StructValue) error {
	for _, f := range s.Fields {
		if !e.encodeValue(f.Value) {
			return e.err
		}
	}
	return e.err
}

func (d *Decoder) DecodeStructValue(s *StructValue) error {
	for _, f := range s.Fields {
		if !d.decodeValue(f.Value) {
			return d.err
		}
	}
	return d.err
}
