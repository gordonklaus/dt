package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type ArrayTypeEditor struct {
	typ  *types.ArrayType
	elem *TypeEditor
}

func NewArrayTypeEditor(typ *types.ArrayType) *ArrayTypeEditor {
	return &ArrayTypeEditor{
		typ:  typ,
		elem: NewTypeEditor(&typ.Elem),
	}
}

func (a *ArrayTypeEditor) Type() types.Type { return a.typ }

func (a *ArrayTypeEditor) Layout(gtx C) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "[]").Layout),
		layout.Rigid(a.elem.Layout),
	)
}
