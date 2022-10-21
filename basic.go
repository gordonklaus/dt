package data

import (
	"github.com/gordonklaus/data/types"
)

type BasicValue struct{ X any }

func NewBasicValue(t *types.BasicType) *BasicValue {
	switch t.Kind {
	case types.Bool:
		return &BasicValue{X: new(bool)}
	case types.Int:
		return &BasicValue{X: new(int)}
	case types.Uint:
		return &BasicValue{X: new(uint)}
	case types.Float32:
		return &BasicValue{X: new(float32)}
	case types.Float64:
		return &BasicValue{X: new(float64)}
	case types.String:
		return &BasicValue{X: new(string)}
	}
	return nil
}

func (e *Encoder) EncodeBasicValue(v *BasicValue) error {
	e.writeBinary(v.X)
	return e.err
}

func (d *Decoder) DecodeBasicValue(v *BasicValue) error {
	d.readBinary(v.X)
	return d.err
}
