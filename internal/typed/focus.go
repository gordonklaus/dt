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

func (f *KeyFocus) FocusEvent(gtx C) bool {
	event.Op(gtx.Ops, f)
	e, _ := gtx.Event(key.FocusFilter{Target: f})
	_, ok := e.(key.FocusEvent)
	return ok
}

func (f *KeyFocus) Event(gtx C, e *key.Event, required, optional key.Modifiers, names ...key.Name) bool {
	event.Op(gtx.Ops, f)
	filters := make([]event.Filter, len(names))
	for i, name := range names {
		filters[i] = key.Filter{Focus: f, Required: required, Optional: optional, Name: name}
	}
	ev, _ := gtx.Event(filters...)
	if ev, ok := ev.(key.Event); ok && ev.State == key.Press {
		*e = ev
		return true
	}
	return false
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
