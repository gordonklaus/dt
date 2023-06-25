package types

import (
	"reflect"
	"testing"

	"github.com/gordonklaus/data/bits"
)

func TestType(t *testing.T) {
	testType(t,
		newInt64Type(),
		&IntType{},
	)
	testType(t,
		&ArrayType{Elem: newStringType()},
		&ArrayType{},
	)
	testType(t,
		&EnumType{Elems: []*EnumElemType{
			{Name: "x", Type: newFloat64Type()},
			{Name: "y", Type: newBoolType()},
		}},
		&EnumType{},
	)
	testType(t,
		&StructType{Fields: []*StructFieldType{
			{Name: "x", Type: newInt64Type()},
			{Name: "y", Type: newInt64Type()},
		}},
		&StructType{},
	)
	testType(t,
		&NamedType{
			Package: &PackageID_Current{},
			Name:    "Bob",
		},
		&NamedType{},
	)
}

func newBoolType() *BoolType     { return &BoolType{} }
func newInt64Type() *IntType     { return &IntType{Size: 64} }
func newFloat64Type() *FloatType { return &FloatType{Size: 64} }
func newStringType() *StringType { return &StringType{} }

func testType(t *testing.T, src, dst Type) {
	b := bits.NewBuffer()
	WriteType(b, src)
	if err := ReadType(b, &dst); err != nil {
		t.Fatal(err)
	}
	if b.Remaining() > 0 {
		t.Errorf("%d bytes remaining", b.Remaining())
	}
	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("Types are not equal:\nsrc = %#v\ndst = %#v", src, dst)
	}
}
