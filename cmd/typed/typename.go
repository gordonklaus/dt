package main

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type TypeNameEditor struct {
	typ   *types.TypeName
	named editor
	typed *TypeEditor

	KeyFocus
	focusTyped KeyFocus
}

func NewTypeNameEditor(typ *types.TypeName, loader *types.Loader) *TypeNameEditor {
	n := &TypeNameEditor{
		typ:   typ,
		named: newEditor(),
	}
	n.typed = NewTypeNameTypeEditor(&n.focusTyped, &typ.Type, loader)
	n.named.SetText(typ.Name)
	return n
}

func (n *TypeNameEditor) Layout(gtx C) D {
	for _, e := range n.KeyFocus.Events(gtx, "→|↓|⏎|⌤|⌫|⌦|⎋") {
		switch e.Name {
		case "→", "↓":
			n.focusTyped.Focus()
		case "⏎", "⌤", "⌫", "⌦":
			n.named.SetCaret(n.named.Len(), n.named.Len())
			n.named.Focus()
		case "⎋":
			if n.named.Focused() {
				n.named.SetText(n.typ.Name)
				n.Focus()
			}
		}
	}

	for _, e := range n.focusTyped.Events(gtx, "←|→|↑|↓|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "←", "↑":
			n.Focus()
		case "→", "↓":
			n.typed.ed.(Focuser).Focus()
		case "⏎", "⌤", "⌫", "⌦":
			n.typed.Edit()
		}
	}

	for _, e := range n.named.Events() {
		switch e := e.(type) {
		case widget.SubmitEvent:
			if validName(e.Text) {
				n.typ.Name = e.Text
				n.Focus()
			}
		}
	}

	if n.typ.Name == "" {
		n.named.Focus()
	} else if n.typ.Type == nil {
		n.typed.Edit()
	}

	axis := layout.Vertical
	if _, ok := n.typ.Type.(*types.StructType); ok {
		axis = layout.Horizontal
	}
	return layout.Flex{
		Axis:      axis,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return n.KeyFocus.Layout(gtx, n.named.Layout)
		}),
		layout.Rigid(layout.Spacer{Width: 4, Height: 4}.Layout),
		layout.Rigid(func(gtx C) D {
			return n.focusTyped.Layout(gtx, n.typed.Layout)
		}),
	)
}

type editor struct {
	widget.Editor
}

func newEditor() editor {
	return editor{
		Editor: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
	}
}

func (ed *editor) Layout(gtx C) D {
	if ed.Focused() {
		key.InputOp{
			Tag:  ed,
			Keys: "←|→|↑|↓",
		}.Add(gtx.Ops)
	}
	return material.Editor(theme, &ed.Editor, "").Layout(gtx)
}

func validName(name string) bool {
	return name != ""
}
