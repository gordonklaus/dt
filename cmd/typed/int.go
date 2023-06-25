package main

import (
	"fmt"

	"gioui.org/io/key"
	"gioui.org/layout"
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
			case "S":
				i.typ.Size /= 2
				if i.typ.Size < 8 {
					i.typ.Size = 64
				}
			case "U":
				i.typ.Unsigned = !i.typ.Unsigned
			}
		}
	}

	key.InputOp{
		Tag:  i,
		Keys: "S|U",
	}.Add(gtx.Ops)

	s := fmt.Sprintf("int%d", i.typ.Size)
	if i.typ.Unsigned {
		s = "u" + s
	}
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(theme, s).Layout),
	)
}
