package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data"
)

type TypeName struct {
	typ *data.NamedType
}

func NewTypeName(typ *data.NamedType) *TypeName {
	return &TypeName{
		typ: typ,
	}
}

func (n *TypeName) Type() data.Type { return n.typ }

func (n *TypeName) Layout(gtx C) D {
	return material.Body1(theme, n.typ.Name).Layout(gtx)
}
