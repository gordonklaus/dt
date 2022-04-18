package data

type StructType struct {
	Fields []*StructFieldType
}

func (e *Encoder) EncodeStructType(s *StructType) error {
	writeSlice(e, e.EncodeStructFieldType, s.Fields)
	return e.err
}

func (d *Decoder) DecodeStructType(s *StructType) error {
	readSlice(d, d.DecodeStructFieldType, &s.Fields)
	return d.err
}

func (s *StructType) NewValue() Value {
	return &StructValue{}
}

type StructFieldType struct {
	Name string
	Type Type
}

func (e *Encoder) EncodeStructFieldType(f *StructFieldType) error {
	_ = e.writeString(f.Name) && e.encodeType(f.Type)
	return e.err
}

func (d *Decoder) DecodeStructFieldType(f *StructFieldType) error {
	_ = d.readString(&f.Name) && d.decodeType(&f.Type)
	return d.err
}
