package data

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
	"golang.org/x/exp/constraints"
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
	slices.SortFunc(m.Elems, func(a, b MapElem) int { return cmp(a.Key, b.Key) })

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

func cmp(a, b Value) int {
	switch a := a.(type) {
	case *IntValue:
		b := b.(*IntValue)
		if a.Type.Unsigned && b.Type.Unsigned {
			return cmpCompare(a.GetUint(), b.GetUint())
		}
		if !a.Type.Unsigned && !b.Type.Unsigned {
			return cmpCompare(a.GetInt(), b.GetInt())
		}
		panic("cannot compare int and uint")
	case *FloatValue:
		b := b.(*FloatValue)
		return cmpCompare(a.x, b.x)
	case *StringValue:
		b := b.(*StringValue)
		return cmpCompare(a.X, b.X)
	}
	panic(fmt.Sprintf("cannot compare %T and %T", a, b))
}

// cmpCompare is a copy of cmp.Compare from the Go 1.21 release.
func cmpCompare[T constraints.Ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}
func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}
