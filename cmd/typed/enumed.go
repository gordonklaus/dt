package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/gordonklaus/data/types"
	"golang.org/x/exp/slices"
)

type EnumTypeEditor struct {
	parent Focuser
	typ    *types.EnumType
	loader *types.Loader
	elems  []*EnumElemTypeEditor
}

func NewEnumTypeEditor(parent Focuser, typ *types.EnumType, loader *types.Loader) *EnumTypeEditor {
	s := &EnumTypeEditor{
		parent: parent,
		typ:    typ,
		loader: loader,
		elems:  make([]*EnumElemTypeEditor, len(typ.Elems)),
	}
	for i, f := range typ.Elems {
		s.elems[i] = NewEnumElemTypeEditor(s, f, loader)
	}
	return s
}

func (e *EnumTypeEditor) Type() types.Type { return e.typ }

func (e *EnumTypeEditor) Focus() {
	if len(e.elems) == 0 {
		e.insertElem(nil, false)
	} else {
		e.elems[0].Focus()
	}
}

func (e *EnumTypeEditor) focusNext(el *EnumElemTypeEditor, next bool) {
	i := slices.Index(e.elems, el) - 1
	if next {
		i += 2
	}
	if i < 0 {
		e.parent.Focus()
	} else if i < len(e.elems) {
		e.elems[i].Focus()
	}
}

func (e *EnumTypeEditor) insertElem(el *EnumElemTypeEditor, before bool) {
	i := slices.Index(e.elems, el)
	if !before {
		i++
	}
	elem := &types.EnumElemType{Type: &types.StructType{}}
	e.typ.Elems = slices.Insert(e.typ.Elems, i, elem)
	e.elems = slices.Insert(e.elems, i, NewEnumElemTypeEditor(e, elem, e.loader))
	e.elems[i].named.Focus()
}

func (e *EnumTypeEditor) deleteElem(el *EnumElemTypeEditor, back bool) {
	i := slices.Index(e.elems, el)
	e.typ.Elems = slices.Delete(e.typ.Elems, i, i+1)
	e.elems = slices.Delete(e.elems, i, i+1)
	if i > 0 && (back || i >= len(e.elems)) {
		i--
	}
	if i < len(e.elems) {
		e.elems[i].Focus()
	} else {
		e.parent.Focus()
	}
}

func (e *EnumTypeEditor) Layout(gtx C) D {
	maxElemNameWidth := 0
	for _, f := range e.elems {
		if x := f.LayoutName(gtx); x > maxElemNameWidth {
			maxElemNameWidth = x
		}
	}
	elems := make([]layout.FlexChild, len(e.typ.Elems))
	for i, f := range e.elems {
		f := f
		elems[i] = layout.Rigid(func(gtx C) D {
			return f.Layout(gtx, maxElemNameWidth)
		})
	}

	elemsRec := Record(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, elems...)
	})

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			width := elemsRec.Dims.Size.X + gtx.Dp(8)
			height := gtx.Dp(16)
			w2 := float32(width) / 2
			h := float32(height)
			d := float32(gtx.Dp(4))
			path := clip.Path{}
			path.Begin(gtx.Ops)
			path.Move(f32.Pt(0, h))
			path.Cube(f32.Pt(0, -h-d), f32.Pt(w2, d), f32.Pt(w2, -h))
			path.Cube(f32.Pt(0, h+d), f32.Pt(w2, -d), f32.Pt(w2, h))
			paint.FillShape(gtx.Ops, theme.Fg, clip.Stroke{
				Path:  path.End(),
				Width: float32(gtx.Dp(1)),
			}.Op())
			return D{Size: image.Pt(width, height)}
		}),
		layout.Rigid(elemsRec.Layout),
	)
}

type EnumElemTypeEditor struct {
	parent *EnumTypeEditor
	typ    *types.EnumElemType
	named  editor
	typed  *StructTypeEditor

	nameRec Recording

	KeyFocus
	focusNamed KeyFocus
}

func NewEnumElemTypeEditor(parent *EnumTypeEditor, typ *types.EnumElemType, loader *types.Loader) *EnumElemTypeEditor {
	f := &EnumElemTypeEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	f.typed = NewStructTypeEditor(f, typ.Type.(*types.StructType), loader)
	f.named.SetText(typ.Name)
	return f
}

func (e *EnumElemTypeEditor) LayoutName(gtx C) int {
	e.nameRec = Record(gtx, e.named.Layout)
	return e.nameRec.Dims.Size.X
}

func (e *EnumElemTypeEditor) Layout(gtx C, nameWidth int) D {
	for _, ev := range e.KeyFocus.Events(gtx, "←|→|↑|↓|(Shift)-[⏎,⌤,⌫,⌦]") {
		switch ev.Name {
		case "←":
			e.focusNamed.Focus()
		case "→":
			e.typed.Focus()
		case "↑":
			e.parent.focusNext(e, false)
		case "↓":
			e.parent.focusNext(e, true)
		case "⏎", "⌤":
			e.parent.insertElem(e, ev.Modifiers == key.ModShift)
		case "⌫", "⌦":
			e.parent.deleteElem(e, (ev.Name == "⌦") == (ev.Modifiers == key.ModShift))
		}
	}

	for _, ev := range e.focusNamed.Events(gtx, "←|→|⏎|⌤|⌫|⌦|⎋") {
		switch ev.Name {
		case "→":
			e.Focus()
		case "←":
			e.parent.focusNext(e, false)
		case "⏎", "⌤", "⌫", "⌦":
			e.named.SetCaret(e.named.Len(), e.named.Len())
			e.named.Focus()
		case "⎋":
			e.named.SetText(e.typ.Name)
			e.Focus()
		}
	}

	for _, ev := range e.named.Events() {
		switch ev := ev.(type) {
		case widget.SubmitEvent:
			if validName(ev.Text) {
				e.typ.Name = ev.Text
				e.Focus()
			}
		}
	}

	if e.Focused() && e.typ.Name == "" {
		e.parent.deleteElem(e, true)
	}

	indent := unit.Dp(float32(nameWidth-e.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 4}.Layout),
		layout.Rigid(func(gtx C) D {
			return e.KeyFocus.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(layout.Spacer{Width: indent}.Layout),
					layout.Rigid(func(gtx C) D {
						return e.focusNamed.Layout(gtx, e.nameRec.Layout)
					}),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(e.typed.Layout),
				)
			})
		}),
	)
}
