package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type IntTypeEditor struct {
	typ types.Type
}

func NewIntTypeEditor(typ types.Type) *IntTypeEditor {
	return &IntTypeEditor{
		typ: typ,
	}
}

func (i *IntTypeEditor) Type() types.Type { return i.typ }

func (i *IntTypeEditor) Layout(gtx C) D {
	s := "uint"
	if _, ok := i.typ.(*types.IntType); ok {
		s = "int"
	}
	return material.Body1(theme, s).Layout(gtx)
}
