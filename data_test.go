package data

import (
	"bytes"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	testValue(t,
		&StructType{Fields: []*StructFieldType{
			{Name: "x", Type: NewBasicType(Int64)},
			{Name: "y", Type: NewBasicType(Int64)},
		}},
		&StructType{},
	)
}

func TestValue(t *testing.T) {
	testValue(t,
		newInt64(42),
		newInt64(0),
	)
	testValue(t,
		&StructValue{Fields: []*StructFieldValue{
			{Value: newInt64(3)},
			{Value: newInt64(7)},
		}},
		&StructValue{Fields: []*StructFieldValue{
			{Value: newInt64(0)},
			{Value: newInt64(0)},
		}},
	)
}

func newInt64(i int64) *BasicValue { return NewBasicValue(&i) }

func testValue(t *testing.T, src, dst Value) {
	var b bytes.Buffer
	if err := NewEncoder(&b).EncodeValue(src); err != nil {
		t.Fatal(err)
	}
	if err := NewDecoder(&b).DecodeValue(dst); err != nil {
		t.Fatal(err)
	}
	if b.Len() > 0 {
		t.Errorf("%d bytes left over", b.Len())
	}
	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("Values are not equal:\nsrc = %#v\ndst = %#v", src, dst)
	}
}
