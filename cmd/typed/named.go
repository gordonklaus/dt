package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type NamedTypeEditor struct {
	typ *types.NamedType
}

func NewNamedTypeEditor(typ *types.NamedType) *NamedTypeEditor {
	return &NamedTypeEditor{
		typ: typ,
	}
}

func (n *NamedTypeEditor) Type() types.Type { return n.typ }

func (n *NamedTypeEditor) Layout(gtx C) D {
	return material.Body1(theme, n.typ.TypeName.Name).Layout(gtx)
}
