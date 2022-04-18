package data

import (
	"fmt"
)

type Type interface {
	NewValue() Value
}

type Kind uint8

const (
	Invalid Kind = iota

	// Basic kinds
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String

	Array
	Struct
	Enum
	Named
)

func (e *Encoder) EncodeType(t Type) error {
	if !e.writeBinary(kind(t)) {
		return e.err
	}

	switch t := t.(type) {
	case *BasicType:
		return nil
	case *StructType:
		return e.EncodeStructType(t)
	case *NamedType:
		return e.EncodeNamedType(t)
	}

	return e.err
}

func kind(t Type) Kind {
	switch t := t.(type) {
	case *BasicType:
		return t.Kind
	case *StructType:
		return Struct
	case *NamedType:
		return Named
	}
	panic(fmt.Sprintf("no Kind for Type %T", t))
}

func (d *Decoder) DecodeType(t *Type) error {
	var kind Kind
	if !d.readBinary(&kind) {
		return d.err
	}

	if kind >= Bool && kind <= String {
		*t = &BasicType{Kind: kind}
		return nil
	}
	switch kind {
	case Array:
	case Struct:
		st := &StructType{}
		*t = st
		return d.DecodeStructType(st)
	case Named:
		nt := &NamedType{}
		*t = nt
		return d.DecodeNamedType(nt)
	}
	panic(fmt.Sprintf("unknown Kind %d", kind))
}
