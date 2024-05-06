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
	named  nameEditor
	typed  *TypeEditor
}

func NewTypeNameEditor(parent Focuser, typ *types.TypeName, core *Core) *TypeNameEditor {
	n := &TypeNameEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	n.typed = NewTypeNameTypeEditor(typ, core)
	n.named.SetText(typ.Name)
	return n
}

func (n *TypeNameEditor) Focus(gtx C) {
	n.named.Focus(gtx)
}

func (n *TypeNameEditor) Layout(gtx C) D {
nevents:
	for {
		var e key.Event
		switch {
		default:
			break nevents
		case n.named.FocusEvent(gtx):
		case n.named.Event(gtx, &e, 0, 0, "←", "↑"):
			n.parent.Focus(gtx)
		case n.named.Event(gtx, &e, 0, 0, "→", "↓"):
			n.typed.Focus(gtx)
		}
	}

	for {
		e, ok := gtx.Event(key.Filter{Focus: &n.named.Editor, Name: "⎋"})
		if !ok {
			break
		}
		if e.(key.Event).State == key.Press {
			n.named.SetText(n.typ.Name)
			n.Focus(gtx)
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

tevents:
	for {
		var e key.Event
		switch {
		default:
			break tevents
		case n.typed.FocusEvent(gtx):
		case n.typed.Event(gtx, &e, 0, 0, "←", "↑"):
			n.Focus(gtx)
		case n.typed.Event(gtx, &e, 0, 0, "→", "↓"):
			n.typed.ed.(Focuser).Focus(gtx)
		}
	}

	if n.typ.Name == "" {
		n.named.Edit(gtx)
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
		layout.Rigid(n.named.Layout),
		layout.Rigid(layout.Spacer{Width: 4, Height: 4}.Layout),
		layout.Rigid(n.typed.Layout),
	)
}

type nameEditor struct {
	KeyFocus
	widget.Editor
}

func newEditor() nameEditor {
	return nameEditor{
		Editor: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
	}
}

func (ed *nameEditor) Edit(gtx C) {
	gtx.Execute(key.FocusCmd{Tag: &ed.Editor})
}

func (ed *nameEditor) Layout(gtx C) D {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case ed.Event(gtx, &e, 0, 0, "⏎", "⌤"):
			ed.SetCaret(ed.Len(), ed.Len())
			ed.Edit(gtx)
		case ed.Event(gtx, &e, 0, 0, "⌫", "⌦"):
			ed.SetCaret(ed.Len(), 0)
			ed.Edit(gtx)
		}
	}
	return ed.KeyFocus.Layout(gtx, material.Editor(theme, &ed.Editor, "").Layout)
}

func validName(name string) bool {
	return name != ""
}
