package types

type NamedType struct {
	Name string
	Type Type
}

func (e *Encoder) EncodeNamedType(n *NamedType) error {
	_ = e.writeString(n.Name) && e.encodeType(n.Type)
	return e.err
}

func (d *Decoder) DecodeNamedType(n *NamedType) error {
	_ = d.readString(&n.Name) && d.decodeType(&n.Type)
	return d.err
}
