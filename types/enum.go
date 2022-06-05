package types

type EnumType struct {
	Elems []*EnumElemType
}

func (e *Encoder) EncodeEnumType(en *EnumType) error {
	writeSlice(e, e.EncodeEnumElemType, en.Elems)
	return e.err
}

func (d *Decoder) DecodeEnumType(e *EnumType) error {
	readSlice(d, d.DecodeEnumElemType, &e.Elems)
	return d.err
}

type EnumElemType struct {
	Name string
	Type Type
}

func (e *Encoder) EncodeEnumElemType(f *EnumElemType) error {
	_ = e.writeString(f.Name) && e.encodeType(f.Type)
	return e.err
}

func (d *Decoder) DecodeEnumElemType(f *EnumElemType) error {
	_ = d.readString(&f.Name) && d.decodeType(&f.Type)
	return d.err
}
