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
	case types.Int8:
		return &BasicValue{X: new(int8)}
	case types.Int16:
		return &BasicValue{X: new(int16)}
	case types.Int32:
		return &BasicValue{X: new(int32)}
	case types.Int64:
		return &BasicValue{X: new(int64)}
	case types.Uint:
		return &BasicValue{X: new(uint)}
	case types.Uint8:
		return &BasicValue{X: new(uint8)}
	case types.Uint16:
		return &BasicValue{X: new(uint16)}
	case types.Uint32:
		return &BasicValue{X: new(uint32)}
	case types.Uint64:
		return &BasicValue{X: new(uint64)}
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
