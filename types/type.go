package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
)

type Type interface {
	bits.ReadWriter
}

type Kind uint64

const (
	Bool Kind = iota
	Uint
	Int
	Float32
	Float64

	Enum
	Struct
	Array

	Option
	String
	Named
)

func (k Kind) String() string {
	switch k {
	case Bool:
		return "bool"
	case Uint:
		return "uint"
	case Int:
		return "int"
	case Float32:
		return "float32"
	case Float64:
		return "float64"

	case Enum:
		return "enum"
	case Struct:
		return "struct"
	case Array:
		return "array"

	case Option:
		return "option"
	case String:
		return "string"
	case Named:
		return "named"
	}
	return fmt.Sprintf("Kind(%d)", k)
}

func WriteType(b *bits.Buffer, t Type) {
	b.WriteVarUint_4bit(uint64(kind(t)))
	t.Write(b)
}

func kind(t Type) Kind {
	switch t.(type) {
	case *BoolType:
		return Bool
	case *UintType:
		return Uint
	case *IntType:
		return Int
	case *Float32Type:
		return Float32
	case *Float64Type:
		return Float64

	case *EnumType:
		return Enum
	case *StructType:
		return Struct
	case *ArrayType:
		return Array

	case *OptionType:
		return Option
	case *StringType:
		return String
	case *NamedType:
		return Named
	}
	panic(fmt.Sprintf("no Kind for Type %T", t))
}

func ReadType(b *bits.Buffer, t *Type) error {
	var k uint64
	if err := b.ReadVarUint_4bit(&k); err != nil {
		return err
	}

	*t = NewType(Kind(k))
	return (*t).Read(b)
}

func NewType(k Kind) Type {
	switch k {
	case Bool:
		return &BoolType{}
	case Uint:
		return &UintType{}
	case Int:
		return &IntType{}
	case Float32:
		return &Float32Type{}
	case Float64:
		return &Float64Type{}

	case Enum:
		return &EnumType{}
	case Struct:
		return &StructType{}
	case Array:
		return &ArrayType{}

	case Option:
		return &OptionType{}
	case String:
		return &StringType{}
	case Named:
		return &NamedType{}
	}
	panic(fmt.Sprintf("unknown Kind %d", k))
}
