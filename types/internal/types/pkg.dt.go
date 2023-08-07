package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	_ = fmt.Print
	_ = bits.NewEncoder
	_ = maps.Keys[map[int]int]
	_ = slices.Sort[[]int]
)

type Type struct{ Type Type__Enum }
type Type__Enum interface {
	isType()
	bits.Value
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

func (x *Type) Write(e *bits.Encoder) {
	switch x.Type.(type) {
	case *Type_Bool:
		e.WriteVarUint_4bit(0)
	case *Type_Int:
		e.WriteVarUint_4bit(1)
	case *Type_Float:
		e.WriteVarUint_4bit(2)
	case *Type_Enum:
		e.WriteVarUint_4bit(3)
	case *Type_Struct:
		e.WriteVarUint_4bit(4)
	case *Type_Array:
		e.WriteVarUint_4bit(5)
	case *Type_Map:
		e.WriteVarUint_4bit(6)
	case *Type_Option:
		e.WriteVarUint_4bit(7)
	case *Type_String:
		e.WriteVarUint_4bit(8)
	case *Type_Named:
		e.WriteVarUint_4bit(9)
	default:
		panic(fmt.Sprintf("invalid Type enum value %T", x))
	}
	x.Type.Write(e)
}

func (x *Type) Read(d *bits.Decoder) error {
	var i uint64
	if err := d.ReadVarUint_4bit(&i); err != nil {
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
	return x.Type.Read(d)
}

type Type_Bool struct{}

func (x *Type_Bool) Write(e *bits.Encoder) {
	e.WriteSize(func() {
	})
}

func (x *Type_Bool) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		return nil
	})
}

type Type_Int struct {
	Unsigned bool
}

func (x *Type_Int) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteBool(x.Unsigned)
	})
}

func (x *Type_Int) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadBool(&x.Unsigned); err != nil {
			return err
		}
		return nil
	})
}

type Type_Float struct {
	Size uint64
}

func (x *Type_Float) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteVarUint(x.Size)
	})
}

func (x *Type_Float) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadVarUint(&x.Size); err != nil {
			return err
		}
		return nil
	})
}

type Type_Enum struct {
	Elements []EnumElement
}

func (x *Type_Enum) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteVarUint(uint64(len(x.Elements)))
		for _, x := range x.Elements {
			(x).Write(e)
		}
	})
}

func (x *Type_Enum) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		{
			var len uint64
			if err := d.ReadVarUint(&len); err != nil {
				return err
			}
			x.Elements = make([]EnumElement, len)
			for i := range x.Elements {
				if err := (&(x.Elements)[i]).Read(d); err != nil {
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

func (x *Type_Struct) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteVarUint(uint64(len(x.Fields)))
		for _, x := range x.Fields {
			(x).Write(e)
		}
	})
}

func (x *Type_Struct) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		{
			var len uint64
			if err := d.ReadVarUint(&len); err != nil {
				return err
			}
			x.Fields = make([]StructField, len)
			for i := range x.Fields {
				if err := (&(x.Fields)[i]).Read(d); err != nil {
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

func (x *Type_Array) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		(x.Element).Write(e)
	})
}

func (x *Type_Array) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := (&x.Element).Read(d); err != nil {
			return err
		}
		return nil
	})
}

type Type_Map struct {
	Key   Type
	Value Type
}

func (x *Type_Map) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		(x.Key).Write(e)
		(x.Value).Write(e)
	})
}

func (x *Type_Map) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := (&x.Key).Read(d); err != nil {
			return err
		}
		if err := (&x.Value).Read(d); err != nil {
			return err
		}
		return nil
	})
}

type Type_Option struct {
	Element Type
}

func (x *Type_Option) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		(x.Element).Write(e)
	})
}

func (x *Type_Option) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := (&x.Element).Read(d); err != nil {
			return err
		}
		return nil
	})
}

type Type_String struct{}

func (x *Type_String) Write(e *bits.Encoder) {
	e.WriteSize(func() {
	})
}

func (x *Type_String) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		return nil
	})
}

type Type_Named struct {
	Package PackageId
	Name    string
}

func (x *Type_Named) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		(x.Package).Write(e)
		e.WriteString(x.Name)
	})
}

func (x *Type_Named) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := (&x.Package).Read(d); err != nil {
			return err
		}
		if err := d.ReadString(&x.Name); err != nil {
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

func (x *EnumElement) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteString(x.Name)
		e.WriteString(x.Doc)
		(x.Type).Write(e)
	})
}

func (x *EnumElement) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadString(&x.Name); err != nil {
			return err
		}
		if err := d.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(d); err != nil {
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

func (x *StructField) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteString(x.Name)
		e.WriteString(x.Doc)
		(x.Type).Write(e)
	})
}

func (x *StructField) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadString(&x.Name); err != nil {
			return err
		}
		if err := d.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(d); err != nil {
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

func (x *Package) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteString(x.Name)
		e.WriteString(x.Doc)
		e.WriteVarUint(uint64(len(x.Types)))
		for _, x := range x.Types {
			(x).Write(e)
		}
	})
}

func (x *Package) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadString(&x.Name); err != nil {
			return err
		}
		if err := d.ReadString(&x.Doc); err != nil {
			return err
		}
		{
			var len uint64
			if err := d.ReadVarUint(&len); err != nil {
				return err
			}
			x.Types = make([]TypeName, len)
			for i := range x.Types {
				if err := (&(x.Types)[i]).Read(d); err != nil {
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

func (x *TypeName) Write(e *bits.Encoder) {
	e.WriteSize(func() {
		e.WriteString(x.Name)
		e.WriteString(x.Doc)
		(x.Type).Write(e)
	})
}

func (x *TypeName) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		if err := d.ReadString(&x.Name); err != nil {
			return err
		}
		if err := d.ReadString(&x.Doc); err != nil {
			return err
		}
		if err := (&x.Type).Read(d); err != nil {
			return err
		}
		return nil
	})
}

type PackageId struct{ PackageId PackageId__Enum }
type PackageId__Enum interface {
	isPackageId()
	bits.Value
}

func (*PackageId_Current) isPackageId() {}

func (x *PackageId) Write(e *bits.Encoder) {
	switch x.PackageId.(type) {
	case *PackageId_Current:
		e.WriteVarUint_4bit(0)
	default:
		panic(fmt.Sprintf("invalid PackageId enum value %T", x))
	}
	x.PackageId.Write(e)
}

func (x *PackageId) Read(d *bits.Decoder) error {
	var i uint64
	if err := d.ReadVarUint_4bit(&i); err != nil {
		return err
	}
	switch i {
	case 0:
		x.PackageId = new(PackageId_Current)
	default:
		x.PackageId = nil // TODO: &PackageId__Unknown{i}
	}
	return x.PackageId.Read(d)
}

type PackageId_Current struct{}

func (x *PackageId_Current) Write(e *bits.Encoder) {
	e.WriteSize(func() {
	})
}

func (x *PackageId_Current) Read(d *bits.Decoder) error {
	return d.ReadSize(func() error {
		return nil
	})
}
