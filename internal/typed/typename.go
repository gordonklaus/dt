package typed

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type TypeNameEditor struct {
	parent Focuser
	typ    *types.TypeName
	named  editor
	typed  *TypeEditor

	KeyFocus
	focusTyped KeyFocus
}

func NewTypeNameEditor(parent Focuser, typ *types.TypeName, loader *types.Loader) *TypeNameEditor {
	n := &TypeNameEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	n.typed = NewTypeNameTypeEditor(&n.focusTyped, &typ.Type, loader)
	n.named.SetText(typ.Name)
	return n
}

func (n *TypeNameEditor) Layout(gtx C) D {
	for _, e := range n.KeyFocus.Events(gtx) {
		switch e.Name {
		case "←", "↑":
			n.parent.Focus(gtx)
		case "→", "↓":
			n.focusTyped.Focus(gtx)
		case "⏎", "⌤":
			n.named.SetCaret(n.named.Len(), n.named.Len())
			n.named.Focus(gtx)
		case "⌫", "⌦":
			n.named.SetCaret(n.named.Len(), 0)
			n.named.Focus(gtx)
		}
	}

	for {
		e, ok := gtx.Event(key.Filter{Focus: &n.named.Editor, Name: "⎋"})
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.Event:
			if e.Name == "⎋" {
				n.named.SetText(n.typ.Name)
				n.Focus(gtx)
			}
		}
	}

	for e, ok := n.named.Update(gtx); ok; e, ok = n.named.Update(gtx) {
		switch e := e.(type) {
		case widget.SubmitEvent:
			if validName(e.Text) {
				n.typ.Name = e.Text
				n.Focus(gtx)
			}
		}
	}

	for _, e := range n.focusTyped.Events(gtx) {
		switch e.Name {
		case "←", "↑":
			n.Focus(gtx)
		case "→", "↓":
			n.typed.ed.(Focuser).Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			n.typed.Edit(gtx)
		}
	}

	if n.typ.Name == "" {
		n.named.Focus(gtx)
	} else if n.typ.Type == nil {
		n.typed.Edit(gtx)
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

func (ed *editor) Focus(gtx C) {
	gtx.Execute(key.FocusCmd{Tag: &ed.Editor})
}

func (ed *editor) Layout(gtx C) D {
	return material.Editor(theme, &ed.Editor, "").Layout(gtx)
}

func validName(name string) bool {
	return name != ""
}
