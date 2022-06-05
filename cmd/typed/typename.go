package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type TypeName struct {
	typ *types.NamedType
}

func NewTypeName(typ *types.NamedType) *TypeName {
	return &TypeName{
		typ: typ,
	}
}

func (n *TypeName) Type() types.Type { return n.typ }

func (n *TypeName) Layout(gtx C) D {
	return material.Body1(theme, n.typ.Name).Layout(gtx)
}
