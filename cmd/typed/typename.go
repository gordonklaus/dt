package main

import (
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type TypeNameEditor struct {
	typ   *types.TypeName
	named widget.Editor
	typed *TypeEditor

	KeyFocus
	focusTyped KeyFocus
}

func NewTypeNameEditor(typ *types.TypeName, loader *types.Loader) *TypeNameEditor {
	n := &TypeNameEditor{
		typ: typ,
		named: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
	}
	n.typed = NewTypeNameTypeEditor(n, &typ.Type, loader)
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
			n.typed.ed.(interface{ Focus() }).Focus()
		case "⏎", "⌤", "⌫", "⌦":
			n.typed.showMenu = true
		}
	}

	for _, e := range n.named.Events() {
		switch e := e.(type) {
		case widget.SubmitEvent:
			n.typ.Name = e.Text
			n.Focus()
		}
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
			return n.KeyFocus.Layout(gtx, material.Editor(theme, &n.named, "").Layout)
		}),
		layout.Rigid(layout.Spacer{Height: 4}.Layout),
		layout.Rigid(func(gtx C) D {
			return n.focusTyped.Layout(gtx, n.typed.Layout)
		}),
	)
}

type KeyFocus struct {
	requestFocus, focused bool
}

func (f *KeyFocus) Focus() {
	f.requestFocus = true
}

func (f *KeyFocus) Focused() bool {
	return f.focused
}

func (f *KeyFocus) Events(gtx C, keys key.Set) []key.Event {
	events := []key.Event{}
	for _, e := range gtx.Events(f) {
		switch e := e.(type) {
		case key.FocusEvent:
			f.focused = e.Focus
		case key.Event:
			if e.State == key.Press {
				events = append(events, e)
			}
		}
	}

	if f.requestFocus {
		f.requestFocus = false
		key.FocusOp{Tag: f}.Add(gtx.Ops)
	}
	key.InputOp{
		Tag:  f,
		Keys: keys,
	}.Add(gtx.Ops)

	return events
}

func (f *KeyFocus) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if f.focused {
				m := gtx.Dp(4)
				defer op.Offset(image.Pt(-m, -m)).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, color.NRGBA{A: 64},
					clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min.Add(image.Pt(2*m, 2*m))}, 2*m).Op(gtx.Ops))
			}
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}
