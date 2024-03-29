package data

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gordonklaus/dt/bits"
	"github.com/gordonklaus/dt/types"
)

func TestBasic(t *testing.T) {
	testValue(t, newBool(false), newBool(false))
	testValue(t, newBool(true), newBool(false))

	testValue(t, newUint64(0), newUint64(0))
	testValue(t, newUint64(42), newUint64(0))
	testValue(t, newUint64(59837), newUint64(0))

	testValue(t, newInt64(0), newInt64(0))
	testValue(t, newInt64(42), newInt64(0))
	testValue(t, newInt64(-123619), newInt64(0))

	testValue(t, newFloat32(0), newFloat32(0))
	testValue(t, newFloat32(439.23422132), newFloat32(0))
	testValue(t, newFloat32(-1.236197e23), newFloat32(0))

	testValue(t, newFloat64(0), newFloat64(0))
	testValue(t, newFloat64(439.23422132), newFloat64(0))
	testValue(t, newFloat64(-1.236197e23), newFloat64(0))

	testValue(t, newString(""), newString(""))
	testValue(t, newString("x"), newString(""))
	testValue(t, newString("abc"), newString(""))
	testValue(t, newString("123456"), newString(""))
	testValue(t, newString("1234567"), newString(""))
	testValue(t, newString("12345678"), newString(""))
	testValue(t, newString("1234567890123456789012345678901234567890"), newString(""))
}

func TestEnum(t *testing.T) {
	et := &types.EnumType{Elems: []*types.EnumElemType{
		{Type: &types.StructType{Fields: []*types.StructFieldType{{Type: intType}}}},
		{Type: &types.StructType{Fields: []*types.StructFieldType{{Type: stringType}}}},
		{Type: &types.StructType{Fields: []*types.StructFieldType{{Type: boolType}}}},
	}}

	ev := NewValue(et).(*EnumValue)
	ev.Elem = 1
	ev.Value = NewValue(et.Elems[1].Type)
	ev.Value.(*StructValue).Fields[0].Value = newString("abc")
	testValue(t, ev, NewValue(et))

	ev.Elem = 4
	ev.Value.(*StructValue).Fields[0].Value = newUint64(42)
	expect := NewValue(et).(*EnumValue)
	expect.Elem = 4
	expect.Value = NewStructValue(&UnknownEnumElementType)
	expect.Value.(*StructValue).HasUnknownFields = true
	testValueExpect(t, ev, NewValue(et), expect)
}

func TestStruct(t *testing.T) {
	st := &types.StructType{Fields: []*types.StructFieldType{
		{Type: intType},
		{Type: stringType},
	}}

	sv := NewValue(st).(*StructValue)
	sv.Fields[0].Value = newInt64(3)
	sv.Fields[1].Value = newString("hello")
	testValue(t, sv, NewValue(st))

	sv = NewValue(st).(*StructValue)
	sv.Fields[1].Value = newString("world")
	testValue(t, sv, NewValue(st))

	testValue(t, NewValue(st), NewValue(st))

	sv = NewValue(st).(*StructValue)
	sv.Fields[1].Value = newString("!#%^")
	sv.Fields = append(sv.Fields, &StructFieldValue{
		Type:  &types.StructFieldType{Type: intType},
		Value: newInt64(9),
	})
	expect := NewValue(st).(*StructValue)
	expect.Fields[1].Value = sv.Fields[1].Value
	expect.HasUnknownFields = true
	testValueExpect(t, sv, NewValue(st), expect)
}

func TestArray(t *testing.T) {
	at := &types.ArrayType{Elem: intType}
	testValue(t,
		&ArrayValue{Type: at, Elems: []Value{newInt64(1), newInt64(2), newInt64(3)}},
		NewValue(at),
	)
}

func TestMap(t *testing.T) {
	at := &types.MapType{Key: intType, Value: stringType}
	testValue(t,
		&MapValue{Type: at, Elems: []MapElem{
			{Key: newInt64(1), Value: newString("abc")},
			{Key: newInt64(2), Value: newString("def")},
			{Key: newInt64(3), Value: newString("ghi")},
		}},
		NewValue(at),
	)
}

func newBool(b bool) *BoolValue {
	x := NewBoolValue(boolType)
	x.X = b
	return x
}

func newUint64(i uint64) *IntValue {
	x := NewIntValue(uintType)
	x.SetUint(i)
	return x
}

func newInt64(i int64) *IntValue {
	x := NewIntValue(intType)
	x.SetInt(i)
	return x
}

func newFloat32(f float32) *FloatValue {
	x := NewFloatValue(float32Type)
	x.SetFloat32(f)
	return x
}

func newFloat64(f float64) *FloatValue {
	x := NewFloatValue(float64Type)
	x.SetFloat64(f)
	return x
}

func newString(f string) *StringValue {
	x := NewStringValue(stringType)
	x.X = f
	return x
}

var boolType = &types.BoolType{}
var intType = &types.IntType{}
var uintType = &types.IntType{Unsigned: true}
var float32Type = &types.FloatType{Size: 32}
var float64Type = &types.FloatType{Size: 64}
var stringType = &types.StringType{}

func testValue(t *testing.T, src, dst Value) {
	testValueExpect(t, src, dst, src)
}

func testValueExpect(t *testing.T, src, dst, expect Value) {
	e := bits.NewEncoder()
	src.Write(e)
	d := bits.NewDecoder(bytes.NewBuffer(e.Bytes()))
	d.SetLimit(e.Size())
	if err := dst.Read(d); err != nil {
		t.Fatal(err)
	}
	if d.Remaining() > 0 {
		t.Errorf("%d bits remaining", d.Remaining())
	}
	if !reflect.DeepEqual(expect, dst) {
		t.Fatalf("Values are not equal:\nexpect: %#v\ngot:    %#v", expect, dst)
	}
}
