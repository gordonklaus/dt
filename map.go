package data

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
	"golang.org/x/exp/slices"
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

func (m *MapValue) Write(b *bits.Buffer) {
	slices.SortFunc(m.Elems, func(a, b MapElem) bool { return less(a.Key, b.Key) })

	b.WriteVarUint(uint64(len(m.Elems)))
	for _, e := range m.Elems {
		e.Key.Write(b)
		e.Value.Write(b)
	}
}

func (a *MapValue) Read(b *bits.Buffer) error {
	var len uint64
	if err := b.ReadVarUint(&len); err != nil {
		return err
	}
	a.Elems = make([]MapElem, len)
	for i := range a.Elems {
		a.Elems[i].Key = NewValue(a.Type.Key)
		if err := a.Elems[i].Key.Read(b); err != nil {
			return err
		}
		a.Elems[i].Value = NewValue(a.Type.Value)
		if err := a.Elems[i].Value.Read(b); err != nil {
			return err
		}
	}
	return nil
}

func less(a, b Value) bool {
	switch a := a.(type) {
	case *IntValue:
		b := b.(*IntValue)
		if a.Type.Unsigned && b.Type.Unsigned {
			return a.GetUint() < b.GetUint()
		}
		if !a.Type.Unsigned && !b.Type.Unsigned {
			return a.GetInt() < b.GetInt()
		}
		panic(fmt.Sprintf("cannot compare int and uint"))
	case *FloatValue:
		b := b.(*FloatValue)
		return a.x < b.x
	case *StringValue:
		b := b.(*StringValue)
		return a.X < b.X
	}
	panic(fmt.Sprintf("cannot compare %T and %T", a, b))
}
