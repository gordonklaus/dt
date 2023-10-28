package typed

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type BoolTypeEditor struct {
	typ *types.BoolType
}

func NewBoolTypeEditor(typ *types.BoolType) *BoolTypeEditor {
	return &BoolTypeEditor{
		typ: typ,
	}
}

func (b *BoolTypeEditor) Type() types.Type { return b.typ }

func (b *BoolTypeEditor) Layout(gtx C) D {
	return material.Body1(theme, "bool").Layout(gtx)
}
