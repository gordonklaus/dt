package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type NamedTypeEditor struct {
	typ   *types.NamedType
	named widget.Editor
	typed *TypeEditor
}

func NewNamedTypeEditor(typ *types.NamedType) *NamedTypeEditor {
	n := &NamedTypeEditor{
		typ: typ,
		named: widget.Editor{
			SingleLine: true,
		},
		typed: NewTypeEditor(&typ.Type),
	}
	n.named.SetText(typ.Name)
	return n
}

func (n *NamedTypeEditor) Layout(gtx C) D {
	for _, e := range n.named.Events() {
		switch e := e.(type) {
		case widget.SubmitEvent:
			n.typ.Name = e.Text
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Editor(theme, &n.named, "").Layout),
		layout.Rigid(n.typed.Layout),
	)
}
