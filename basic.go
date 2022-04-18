package data

type BasicValue struct{ X any }

func NewBasicValue[
	T interface {
		bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string
	},
](x *T) *BasicValue {
	return &BasicValue{X: x}
}

func (e *Encoder) EncodeBasicValue(v *BasicValue) error {
	e.writeBinary(v.X)
	return e.err
}

func (d *Decoder) DecodeBasicValue(v *BasicValue) error {
	d.readBinary(v.X)
	return d.err
}
