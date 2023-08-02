package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	_ = fmt.Print
	_ = bits.NewBuffer
	_ = maps.Keys[map[int]int]
	_ = slices.Sort[[]int]
)

type Type struct{ Type Type__Enum }
type Type__Enum interface {
	isType()
	bits.ReadWriter
}

func (*Type_Bool) isType()   {}
func (*Type_Int) isType()    {}
func (*Type_Float) isType()  {}
func (*Type_Enum) isType()   {}
func (*Type_Struct) isType() {}
func (*Type_Array) isType()  {}
func (*Type_Map) isType()    {}
func (*Type_Option) isType() {}
func (*Type_String) isType() {}
func (*Type_Named) isType()  {}

func (x *Type) Write(b *bits.Buffer) {
	switch x.Type.(type) {
	case *Type_Bool:
		b.WriteVarUint_4bit(0)
	case *Type_Int:
		b.WriteVarUint_4bit(1)
	case *Type_Float:
		b.WriteVarUint_4bit(2)
	case *Type_Enum:
		b.WriteVarUint_4bit(3)
	case *Type_Struct:
		b.WriteVarUint_4bit(4)
	case *Type_Array:
		b.WriteVarUint_4bit(5)
	case *Type_Map:
		b.WriteVarUint_4bit(6)
	case *Type_Option:
		b.WriteVarUint_4bit(7)
	case *Type_String:
		b.WriteVarUint_4bit(8)
	case *Type_Named:
		b.WriteVarUint_4bit(9)
	default:
		panic(fmt.Sprintf("invalid Type enum value %T", x))
	}
	x.Type.Write(b)
}

func (x *Type) Read(b *bits.Buffer) error {
	var i uint64
	if err := b.ReadVarUint_4bit(&i); err != nil {
		return err
	}
	switch i {
	case 0:
		x.Type = new(Type_Bool)
	case 1:
		x.Type = new(Type_Int)
	case 2:
		x.Type = new(Type_Float)
	case 3:
		x.Type = new(Type_Enum)
	case 4:
		x.Type = new(Type_Struct)
	case 5:
		x.Type = new(Type_Array)
	case 6:
		x.Type = new(Type_Map)
	case 7:
		x.Type = new(Type_Option)
	case 8:
		x.Type = new(Type_String)
	case 9:
		x.Type = new(Type_Named)
	default:
		x.Type = nil // TODO: &Type__Unknown{i}
	}
	return x.Type.Read(b)
}

type Type_Bool struct{}

func (x *Type_Bool) Write(b *bits.Buffer) {
	b.WriteSize(func() {
	})
}

func (x *Type_Bool) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return nil
	})
}

type Type_Int struct {
	Unsigned bool
}

func (x *Type_Int) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteBool(x.Unsigned)
	})
}

func (x *Type_Int) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadBool(&x.Unsigned); err != nil {
			return err
		}
		return nil
	})
}

type Type_Float struct {
	Size uint64
}

func (x *Type_Float) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(x.Size)
	})
}

func (x *Type_Float) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadVarUint(&x.Size); err != nil {
			return err
		}
		return nil
	})
}

type Type_Enum struct {
	Elements []EnumElement
}

func (x *Type_Enum) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(uint64(len(x.Elements)))
		for _, x := range x.Elements {
			(x).Write(b)
		}
	})
}

func (x *Type_Enum) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		{
			var len uint64
			if err := b.ReadVarUint(&len); err != nil {
				return err
			}
			x.Elements = make([]EnumElement, len)
			for i := range x.Elements {
				if err := (&(x.Elements)[i]).Read(b); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type Type_Struct struct {
	Fields []StructField
}

func (x *Type_Struct) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint(uint64(len(x.Fields)))
		for _, x := range x.Fields {
			(x).Write(b)
		}
	})
}

func (x *Type_Struct) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		{
			var len uint64
			if err := b.ReadVarUint(&len); err != nil {
				return err
			}
			x.Fields = make([]StructField, len)
			for i := range x.Fields {
				if err := (&(x.Fields)[i]).Read(b); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type Type_Array struct {
	Element Type
}

func (x *Type_Array) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		(x.Element).Write(b)
	})
}

func (x *Type_Array) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := (&x.Element).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type Type_Map struct {
	Key   Type
	Value Type
}

func (x *Type_Map) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		(x.Key).Write(b)
		(x.Value).Write(b)
	})
}

func (x *Type_Map) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := (&x.Key).Read(b); err != nil {
			return err
		}
		if err := (&x.Value).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type Type_Option struct {
	Element Type
}

func (x *Type_Option) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		(x.Element).Write(b)
	})
}

func (x *Type_Option) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := (&x.Element).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type Type_String struct{}

func (x *Type_String) Write(b *bits.Buffer) {
	b.WriteSize(func() {
	})
}

func (x *Type_String) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return nil
	})
}

type Type_Named struct {
	Package PackageId
	Name    string
}

func (x *Type_Named) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		(x.Package).Write(b)
		b.WriteString(x.Name)
	})
}

func (x *Type_Named) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := (&x.Package).Read(b); err != nil {
			return err
		}
		if err := b.ReadString(&x.Name); err != nil {
			return err
		}
		return nil
	})
}

type EnumElement struct {
	Name string
	Doc  string
	Type Type
}

func (x *EnumElement) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(x.Name)
		b.WriteString(x.Doc)
		(x.Type).Write(b)
	})
}

func (x *EnumElement) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&x.Name); err != nil {
			return err
		}
		if err := b.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type StructField struct {
	Name string
	Doc  string
	Type Type
}

func (x *StructField) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(x.Name)
		b.WriteString(x.Doc)
		(x.Type).Write(b)
	})
}

func (x *StructField) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&x.Name); err != nil {
			return err
		}
		if err := b.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type Package struct {
	Name  string
	Doc   string
	Types []TypeName
}

func (x *Package) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(x.Name)
		b.WriteString(x.Doc)
		b.WriteVarUint(uint64(len(x.Types)))
		for _, x := range x.Types {
			(x).Write(b)
		}
	})
}

func (x *Package) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&x.Name); err != nil {
			return err
		}
		if err := b.ReadString(&x.Doc); err != nil {
			return err
		}
		{
			var len uint64
			if err := b.ReadVarUint(&len); err != nil {
				return err
			}
			x.Types = make([]TypeName, len)
			for i := range x.Types {
				if err := (&(x.Types)[i]).Read(b); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type TypeName struct {
	Name string
	Doc  string
	Type Type
}

func (x *TypeName) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(x.Name)
		b.WriteString(x.Doc)
		(x.Type).Write(b)
	})
}

func (x *TypeName) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&x.Name); err != nil {
			return err
		}
		if err := b.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(b); err != nil {
			return err
		}
		return nil
	})
}

type PackageId struct{ PackageId PackageId__Enum }
type PackageId__Enum interface {
	isPackageId()
	bits.ReadWriter
}

func (*PackageId_Current) isPackageId() {}

func (x *PackageId) Write(b *bits.Buffer) {
	switch x.PackageId.(type) {
	case *PackageId_Current:
		b.WriteVarUint_4bit(0)
	default:
		panic(fmt.Sprintf("invalid PackageId enum value %T", x))
	}
	x.PackageId.Write(b)
}

func (x *PackageId) Read(b *bits.Buffer) error {
	var i uint64
	if err := b.ReadVarUint_4bit(&i); err != nil {
		return err
	}
	switch i {
	case 0:
		x.PackageId = new(PackageId_Current)
	default:
		x.PackageId = nil // TODO: &PackageId__Unknown{i}
	}
	return x.PackageId.Read(b)
}

type PackageId_Current struct{}

func (x *PackageId_Current) Write(b *bits.Buffer) {
	b.WriteSize(func() {
	})
}

func (x *PackageId_Current) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		return nil
	})
}
