package data

import "fmt"

type Value interface{}

func (e *Encoder) EncodeValue(v Value) error {
	switch v := v.(type) {
	case *BasicValue:
		return e.EncodeBasicValue(v)
	case *StructValue:
		return e.EncodeStruct(v)
	}
	panic(fmt.Sprintf("invalid Value type %T", v))
}

func (d *Decoder) DecodeValue(v Value) error {
	switch v := v.(type) {
	case *BasicValue:
		return d.DecodeBasicValue(v)
	case *StructValue:
		return d.DecodeStruct(v)
	}
	panic(fmt.Sprintf("invalid Value type %T", v))
}

func (e *Encoder) encodeValue(v Value) bool {
	e.err = e.EncodeValue(v)
	return e.err == nil
}

func (d *Decoder) decodeValue(v Value) bool {
	d.err = d.DecodeValue(v)
	return d.err == nil
}
