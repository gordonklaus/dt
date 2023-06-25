package main

import (
	"gioui.org/io/key"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type IntTypeEditor struct {
	typ *types.IntType
}

func NewIntTypeEditor(typ *types.IntType) *IntTypeEditor {
	return &IntTypeEditor{
		typ: typ,
	}
}

func (i *IntTypeEditor) Type() types.Type { return i.typ }

func (i *IntTypeEditor) Layout(gtx C) D {
	for _, e := range gtx.Events(i) {
		if e, ok := e.(key.Event); ok && e.State == key.Press {
			switch e.Name {
			case "U":
				i.typ.Unsigned = !i.typ.Unsigned
			}
		}
	}

	key.InputOp{
		Tag:  i,
		Keys: "U",
	}.Add(gtx.Ops)

	s := "int"
	if i.typ.Unsigned {
		s = "uint"
	}
	return material.Body1(theme, s).Layout(gtx)
}
