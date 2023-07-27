package types

import (
	"reflect"
	"testing"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types/internal/types"
)

func TestType(t *testing.T) {
	testType(t,
		newIntType(),
	)
	testType(t,
		&ArrayType{Elem: newStringType()},
	)
	testType(t,
		&MapType{Key: newStringType(), Value: newFloat64Type()},
	)
	testType(t,
		&EnumType{Elems: []*EnumElemType{
			{Name: "x", Type: &StructType{Fields: []*StructFieldType{{Type: newFloat64Type()}}}},
			{Name: "y", Type: &StructType{Fields: []*StructFieldType{{Type: newBoolType()}}}},
		}},
	)
	testType(t,
		&StructType{Fields: []*StructFieldType{
			{Name: "x", Type: newIntType()},
			{Name: "y", Type: newIntType()},
		}},
	)
	testType(t,
		&NamedType{
			Package: &PackageID_Current{},
			Name:    "Bob",
		},
	)
}

func newBoolType() *BoolType     { return &BoolType{} }
func newIntType() *IntType       { return &IntType{} }
func newFloat64Type() *FloatType { return &FloatType{Size: 64} }
func newStringType() *StringType { return &StringType{} }

func testType(t *testing.T, src Type) {
	srctyp := typeToData(src)
	b := bits.NewBuffer()
	srctyp.Write(b)
	var dsttyp types.Type
	if err := dsttyp.Read(b); err != nil {
		t.Fatal(err)
	}
	if b.Remaining() > 0 {
		t.Errorf("%d bytes remaining", b.Remaining())
	}
	dst := typeFromData(dsttyp)
	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("Types are not equal:\nsrc = %#v\ndst = %#v", src, dst)
	}
}
