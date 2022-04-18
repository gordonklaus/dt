package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data"
)

type NamedTypeEditor struct {
	typ   *data.NamedType
	named widget.Editor
	typed *TypeEditor
}

func NewNamedTypeEditor(typ *data.NamedType) *NamedTypeEditor {
	return &NamedTypeEditor{
		typ: typ,
		named: widget.Editor{
			SingleLine: true,
		},
		typed: NewTypeEditor(&typ.Type),
	}
}

func (n *NamedTypeEditor) Layout(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, material.Editor(theme, &n.named, "name").Layout),
		layout.Flexed(1, n.typed.Layout),
	)
}
