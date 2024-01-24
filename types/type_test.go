package types

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types/internal/types"
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
}

func newBoolType() *BoolType     { return &BoolType{} }
func newIntType() *IntType       { return &IntType{} }
func newFloat64Type() *FloatType { return &FloatType{Size: 64} }
func newStringType() *StringType { return &StringType{} }

func testType(t *testing.T, src Type) {
	l := Loader{}
	srctyp := l.typeToData(src)
	e := bits.NewEncoder()
	srctyp.Write(e)
	var dsttyp types.Type
	d := bits.NewDecoder(bytes.NewBuffer(e.Bytes()))
	d.SetLimit(e.Size())
	if err := dsttyp.Read(d); err != nil {
		t.Fatal(err)
	}
	if d.Remaining() > 0 {
		t.Errorf("%d bits remaining", d.Remaining())
	}
	dst := l.typeFromData(dsttyp, nil)
	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("Types are not equal:\nsrc = %#v\ndst = %#v", src, dst)
	}
}
