package main

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type TypeNameEditor struct {
	typ   *types.TypeName
	named widget.Editor
	typed *TypeEditor
}

func NewTypeNameEditor(typ *types.TypeName, loader *types.Loader) *TypeNameEditor {
	n := &TypeNameEditor{
		typ: typ,
		named: widget.Editor{
			SingleLine: true,
		},
		typed: NewTypeNameTypeEditor(&typ.Type, loader),
	}
	n.named.SetText(typ.Name)
	return n
}

func (n *TypeNameEditor) Layout(gtx C) D {
	for _, e := range n.named.Events() {
		switch e.(type) {
		case widget.ChangeEvent:
			n.typ.Name = n.named.Text()
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Editor(theme, &n.named, "").Layout),
		layout.Rigid(n.typed.Layout),
	)
}
