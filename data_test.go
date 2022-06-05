package data

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gordonklaus/data/types"
)

func TestValue(t *testing.T) {
	testValue(t,
		newInt64(42),
		newInt64(0),
	)
	at := &types.ArrayType{Elem: &types.BasicType{Kind: types.Int64}}
	testValue(t,
		&ArrayValue{Type: at, Elems: []Value{
			newInt64(1),
			newInt64(2),
			newInt64(3),
		}},
		&ArrayValue{Type: at},
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

func newInt64(i int64) *BasicValue { return &BasicValue{X: &i} }

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
