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
	Int
	Float

	Enum
	Struct
	Array
	Map

	Option
	String
	Named
)

func (k Kind) String() string {
	switch k {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Float:
		return "float"

	case Enum:
		return "enum"
	case Struct:
		return "struct"
	case Array:
		return "array"
	case Map:
		return "map"

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
	case *IntType:
		return Int
	case *FloatType:
		return Float

	case *EnumType:
		return Enum
	case *StructType:
		return Struct
	case *ArrayType:
		return Array
	case *MapType:
		return Map

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
	case Int:
		return &IntType{}
	case Float:
		return &FloatType{}

	case Enum:
		return &EnumType{}
	case Struct:
		return &StructType{}
	case Array:
		return &ArrayType{}
	case Map:
		return &MapType{}

	case Option:
		return &OptionType{}
	case String:
		return &StringType{}
	case Named:
		return &NamedType{}
	}
	panic(fmt.Sprintf("unknown Kind %d", k))
}
