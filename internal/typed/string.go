package typed

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type StringTypeEditor struct {
	typ *types.StringType
}

func NewStringTypeEditor(typ *types.StringType) *StringTypeEditor {
	return &StringTypeEditor{
		typ: typ,
	}
}

func (b *StringTypeEditor) Type() types.Type { return b.typ }

func (b *StringTypeEditor) Layout(gtx C) D {
	return material.Body1(theme, "string").Layout(gtx)
}
