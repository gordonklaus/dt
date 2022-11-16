package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type BasicTypeEditor struct {
	typ types.Type
}

func NewBasicTypeEditor(typ types.Type) *BasicTypeEditor {
	return &BasicTypeEditor{
		typ: typ,
	}
}

func (b *BasicTypeEditor) Type() types.Type { return b.typ }

func (b *BasicTypeEditor) Layout(gtx C) D {
	s := ""
	switch b.typ.(type) {
	case *types.BoolType:
		s = "bool"
	case *types.Float32Type:
		s = "float32"
	case *types.Float64Type:
		s = "float64"
	case *types.StringType:
		s = "string"
	}
	return material.Body1(theme, s).Layout(gtx)
}
