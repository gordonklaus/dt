package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type BasicType struct {
	typ *types.BasicType
}

func NewBasic(typ *types.BasicType) *BasicType {
	return &BasicType{
		typ: typ,
	}
}

func (n *BasicType) Type() types.Type { return n.typ }

func (n *BasicType) Layout(gtx C) D {
	return material.Body1(theme, n.typ.Kind.String()).Layout(gtx)
}
