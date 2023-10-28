package typed

import (
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
	s := "int"
	if i.typ.Unsigned {
		s = "uint"
	}
	return material.Body1(theme, s).Layout(gtx)
}
