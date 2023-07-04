package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type OptionTypeEditor struct {
	typ *types.OptionType
	val *TypeEditor
}

func NewOptionTypeEditor(typ *types.OptionType, loader *types.Loader) *OptionTypeEditor {
	return &OptionTypeEditor{
		typ: typ,
		val: NewTypeEditor(&typ.Elem, loader),
	}
}

func (o *OptionTypeEditor) Type() types.Type { return o.typ }

func (o *OptionTypeEditor) Layout(gtx C) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "?").Layout),
		layout.Rigid(o.val.Layout),
	)
}
