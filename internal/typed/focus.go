package typed

import (
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Focuser interface {
	Focus(gtx C)
}

type KeyFocus struct {
	focused bool
}

func (f *KeyFocus) Focus(gtx C) {
	key.FocusOp{Tag: f}.Add(gtx.Ops)
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
