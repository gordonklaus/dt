package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data"
)

type StructTypeEditor struct {
	typ *data.StructType
}

func NewStructTypeEditor(typ *data.StructType) *StructTypeEditor {
	return &StructTypeEditor{
		typ: typ,
	}
}

func (n *StructTypeEditor) Type() data.Type { return n.typ }

func (n *StructTypeEditor) Layout(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "struct").Layout),
	)
}
