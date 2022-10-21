package data

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gordonklaus/data/types"
)

func TestBasic(t *testing.T) {
	testValue(t, newInt(42), newInt(0))
}

func TestArray(t *testing.T) {
	at := &types.ArrayType{Elem: &types.BasicType{Kind: types.Int}}
	testValue(t,
		&ArrayValue{Type: at, Elems: []Value{newInt(1), newInt(2), newInt(3)}},
		NewValue(at),
	)
}

func TestStruct(t *testing.T) {
	st := &types.StructType{Fields: []*types.StructFieldType{
		{Type: &types.BasicType{Kind: types.Int}},
		{Type: &types.BasicType{Kind: types.Int}},
	}}

	sv := NewValue(st).(*StructValue)
	sv.Fields[0].Value = newInt(3)
	sv.Fields[1].Value = newInt(7)
	testValue(t, sv, NewValue(st))

	sv = NewValue(st).(*StructValue)
	sv.Fields[1].Value = newInt(5)
	testValue(t, sv, NewValue(st))

	testValue(t, NewValue(st), NewValue(st))

	sv = NewValue(st).(*StructValue)
	sv.Fields = append(sv.Fields, &StructFieldValue{
		StructFieldType: &types.StructFieldType{Type: &types.BasicType{Kind: types.Int}},
		Value:           newInt(9),
	})
	sv.Fields[1].Value = newInt(5)
	expect := NewValue(st).(*StructValue)
	expect.Fields[1].Value = newInt(5)
	testValueExpect(t, sv, NewValue(st), expect)
}

func newInt(i int) *BasicValue { return &BasicValue{X: &i} }

func testValue(t *testing.T, src, dst Value) {
	testValueExpect(t, src, dst, src)
}

func testValueExpect(t *testing.T, src, dst, expect Value) {
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
	if !reflect.DeepEqual(expect, dst) {
		t.Fatalf("Values are not equal:\nexpect: %#v\ngot:    %#v", expect, dst)
	}
}
