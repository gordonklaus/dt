package types

type ArrayType struct {
	Elem Type
}

func (e *Encoder) EncodeArrayType(a *ArrayType) error {
	return e.EncodeType(a.Elem)
}

func (d *Decoder) DecodeArrayType(a *ArrayType) error {
	return d.DecodeType(&a.Elem)
}
