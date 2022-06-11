package types

import (
	"bytes"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	testType(t,
		NewBasicType(Int64),
		&BasicType{},
	)
	testType(t,
		&ArrayType{Elem: NewBasicType(String)},
		&ArrayType{},
	)
	testType(t,
		&EnumType{Elems: []*EnumElemType{
			{Name: "x", Type: NewBasicType(Float64)},
			{Name: "y", Type: NewBasicType(Bool)},
		}},
		&EnumType{},
	)
	testType(t,
		&StructType{Fields: []*StructFieldType{
			{Name: "x", Type: NewBasicType(Int64)},
			{Name: "y", Type: NewBasicType(Int64)},
		}},
		&StructType{},
	)
	testType(t,
		&NamedType{
			Name: "Bob",
			Type: NewBasicType(Uint16),
		},
		&NamedType{},
	)
}

func testType(t *testing.T, src, dst Type) {
	var b bytes.Buffer
	if err := NewEncoder(&b).EncodeType(src); err != nil {
		t.Fatal(err)
	}
	if err := NewDecoder(&b).DecodeType(&dst); err != nil {
		t.Fatal(err)
	}
	if b.Len() > 0 {
		t.Errorf("%d bytes left over", b.Len())
	}
	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("Types are not equal:\nsrc = %#v\ndst = %#v", src, dst)
	}
}
