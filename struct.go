package data

type StructValue struct {
	Fields []*StructFieldValue
}

type StructFieldValue struct {
	Name  string
	Value Value
}

func (e *Encoder) EncodeStruct(s *StructValue) error {
	for _, f := range s.Fields {
		if !e.encodeValue(f.Value) {
			return e.err
		}
	}
	return e.err
}

func (d *Decoder) DecodeStruct(s *StructValue) error {
	for _, f := range s.Fields {
		if !d.decodeValue(f.Value) {
			return d.err
		}
	}
	return d.err
}
