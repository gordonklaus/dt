package main

import (
	"fmt"

	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type FloatTypeEditor struct {
	typ *types.FloatType
}

func NewFloatTypeEditor(typ *types.FloatType) *FloatTypeEditor {
	return &FloatTypeEditor{
		typ: typ,
	}
}

func (b *FloatTypeEditor) Type() types.Type { return b.typ }

func (b *FloatTypeEditor) Layout(gtx C) D {
	return material.Body1(theme, fmt.Sprintf("float%d", b.typ.Size)).Layout(gtx)
}
