package data

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type MapValue struct {
	Type  *types.MapType
	Elems []MapElem
}

type MapElem struct {
	Key, Value Value // Key is *IntValue, *FloatValue, or *StringValue
}

func NewMapValue(t *types.MapType) *MapValue {
	return &MapValue{
		Type: t,
	}
}

func (m *MapValue) Write(e *bits.Encoder) {
	slices.SortFunc(m.Elems, func(a, b MapElem) int { return compare(a.Key, b.Key) })

	e.WriteVarUint(uint64(len(m.Elems)))
	for _, el := range m.Elems {
		el.Key.Write(e)
		el.Value.Write(e)
	}
}

func (a *MapValue) Read(d *bits.Decoder) error {
	var len uint64
	if err := d.ReadVarUint(&len); err != nil {
		return err
	}
	a.Elems = make([]MapElem, len)
	for i := range a.Elems {
		a.Elems[i].Key = NewValue(a.Type.Key)
		if err := a.Elems[i].Key.Read(d); err != nil {
			return err
		}
		a.Elems[i].Value = NewValue(a.Type.Value)
		if err := a.Elems[i].Value.Read(d); err != nil {
			return err
		}
	}
	return nil
}

func compare(a, b Value) int {
	switch a := a.(type) {
	case *IntValue:
		b := b.(*IntValue)
		if a.Type.Unsigned && b.Type.Unsigned {
			return cmp.Compare(a.GetUint(), b.GetUint())
		}
		if !a.Type.Unsigned && !b.Type.Unsigned {
			return cmp.Compare(a.GetInt(), b.GetInt())
		}
		panic("cannot compare int and uint")
	case *FloatValue:
		b := b.(*FloatValue)
		return cmp.Compare(a.x, b.x)
	case *StringValue:
		b := b.(*StringValue)
		return cmp.Compare(a.X, b.X)
	}
	panic(fmt.Sprintf("cannot compare %T and %T", a, b))
}
