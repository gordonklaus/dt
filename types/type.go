package types

import (
	"fmt"
)

type Type interface {
	isType()
}

func (*BasicType) isType()  {}
func (*ArrayType) isType()  {}
func (*EnumType) isType()   {}
func (*StructType) isType() {}
func (*NamedType) isType()  {}

type Kind uint8

const (
	Bool Kind = iota
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
	Enum
	Struct
	Named
)

func (k Kind) String() string {
	switch k {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Int8:
		return "int8"
	case Int16:
		return "int16"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case Uint:
		return "uint"
	case Uint8:
		return "uint8"
	case Uint16:
		return "uint16"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case String:
		return "string"
	case Array:
		return "array"
	case Enum:
		return "enum"
	case Struct:
		return "struct"
	case Named:
		return "named"
	}
	return fmt.Sprintf("Kind(%d)", k)
}

func (e *Encoder) EncodeType(t Type) error {
	if !e.writeBinary(kind(t)) {
		return e.err
	}

	switch t := t.(type) {
	case *BasicType:
		return nil
	case *ArrayType:
		return e.EncodeArrayType(t)
	case *EnumType:
		return e.EncodeEnumType(t)
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
	case *ArrayType:
		return Array
	case *EnumType:
		return Enum
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
		at := &ArrayType{}
		*t = at
		return d.DecodeArrayType(at)
	case Enum:
		et := &EnumType{}
		*t = et
		return d.DecodeEnumType(et)
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
