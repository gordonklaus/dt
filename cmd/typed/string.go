package main

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type StringTypeEditor struct {
	typ types.Type
}

func NewStringTypeEditor(typ types.Type) *StringTypeEditor {
	return &StringTypeEditor{
		typ: typ,
	}
}

func (s *StringTypeEditor) Type() types.Type { return s.typ }

func (s *StringTypeEditor) Layout(gtx C) D {
	return material.Body1(theme, "string").Layout(gtx)
}
