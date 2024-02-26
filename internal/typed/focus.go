package typed

import (
	"image"
	"image/color"

	"gioui.org/io/event"
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
	_ int // because pointers to empty structs may not be unique
}

func (f *KeyFocus) Focus(gtx C) {
	gtx.Execute(key.FocusCmd{Tag: f})
}

func (f *KeyFocus) Focused(gtx C) bool {
	return gtx.Focused(f)
}

func (f *KeyFocus) Events(gtx C) []key.Event {
	events := []key.Event{}
	for {
		e, ok := gtx.Event(
			key.FocusFilter{Target: f},
			key.Filter{Focus: f, Optional: key.ModCommand | key.ModShift | key.ModAlt | key.ModSuper},
		)
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.FocusEvent:
		case key.Event:
			if e.State == key.Press {
				events = append(events, e)
			}
		}
	}

	event.Op(gtx.Ops, f)

	return events
}

func (f *KeyFocus) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if f.Focused(gtx) {
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
