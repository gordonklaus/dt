package typed

import (
	"image"
	"slices"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/gordonklaus/dt/types"
)

type EnumTypeEditor struct {
	*Core
	parent *TypeEditor
	typ    *types.EnumType
	elems  []*EnumElemTypeEditor
}

func NewEnumTypeEditor(parent *TypeEditor, typ *types.EnumType, core *Core) *EnumTypeEditor {
	s := &EnumTypeEditor{
		parent: parent,
		typ:    typ,
		Core:   core,
		elems:  make([]*EnumElemTypeEditor, len(typ.Elems)),
	}
	for i, f := range typ.Elems {
		s.elems[i] = NewEnumElemTypeEditor(s, f, core)
	}
	return s
}

func (e *EnumTypeEditor) Type() types.Type { return e.typ }

func (e *EnumTypeEditor) Focus(gtx C) {
	if len(e.elems) == 0 {
		e.insertElem(gtx, nil, false)
	} else {
		e.elems[0].Focus(gtx)
	}
}

func (e *EnumTypeEditor) focusNext(gtx C, el *EnumElemTypeEditor, next bool) {
	i := slices.Index(e.elems, el) - 1
	if next {
		i += 2
	}
	if i < 0 {
		e.parent.Focus(gtx)
	} else if i < len(e.elems) {
		e.elems[i].Focus(gtx)
	}
}

func (e *EnumTypeEditor) insertElem(gtx C, el *EnumElemTypeEditor, before bool) {
	i := slices.Index(e.elems, el)
	if !before {
		i++
	}
	elem := &types.EnumElemType{
		ID:     nextID(e.Pkg),
		Type:   &types.StructType{},
		Parent: e.parent.typeName,
	}
	e.typ.Elems = slices.Insert(e.typ.Elems, i, elem)
	e.Pkg.TypesByID[elem.ID] = elem
	e.elems = slices.Insert(e.elems, i, NewEnumElemTypeEditor(e, elem, e.Core))
	e.elems[i].named.Edit(gtx)
}

func (e *EnumTypeEditor) swap(el *EnumElemTypeEditor, next bool) {
	i := slices.Index(e.elems, el)
	if next && i == len(e.elems)-1 || !next && i == 0 {
		return
	}
	if !next {
		i--
	}
	e.typ.Elems[i], e.typ.Elems[i+1] = e.typ.Elems[i+1], e.typ.Elems[i]
	e.elems[i], e.elems[i+1] = e.elems[i+1], e.elems[i]
}

func (e *EnumTypeEditor) deleteElem(gtx C, el *EnumElemTypeEditor, back bool) {
	for len(el.typed.fields) > 0 {
		el.typed.deleteField(gtx, el.typed.fields[0], true)
	}
	i := slices.Index(e.elems, el)
	e.typ.Elems = slices.Delete(e.typ.Elems, i, i+1)
	delete(e.Pkg.TypesByID, el.typ.ID)
	e.elems = slices.Delete(e.elems, i, i+1)
	if i > 0 && (back || i >= len(e.elems)) {
		i--
	}
	if i < len(e.elems) {
		e.elems[i].Focus(gtx)
	} else {
		e.parent.Focus(gtx)
	}
}

func (e *EnumTypeEditor) Layout(gtx C) D {
	// Iterate over a copy because e.elems may mutate during iteration.
	for _, f := range slices.Clone(e.elems) {
		f.Update(gtx)
	}

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
	KeyFocus
	named nameEditor
	typed *StructTypeEditor

	nameRec Recording
}

func NewEnumElemTypeEditor(parent *EnumTypeEditor, typ *types.EnumElemType, core *Core) *EnumElemTypeEditor {
	f := &EnumElemTypeEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	f.typed = NewStructTypeEditor(f, typ.Type.(*types.StructType), core)
	f.named.SetText(typ.Name)
	return f
}

func (e *EnumElemTypeEditor) Update(gtx C) {
events:
	for {
		var ev key.Event
		switch {
		default:
			break events
		case e.FocusEvent(gtx):
		case e.Event(gtx, &ev, 0, 0, "←"):
			e.named.Focus(gtx)
		case e.Event(gtx, &ev, 0, 0, "→"):
			e.typed.Focus(gtx)
		case e.Event(gtx, &ev, 0, key.ModShift, "↑", "↓"):
			if ev.Modifiers == key.ModShift {
				e.parent.swap(e, ev.Name == "↓")
			} else {
				e.parent.focusNext(gtx, e, ev.Name == "↓")
			}
		case e.Event(gtx, &ev, 0, key.ModShift, "⏎", "⌤"):
			e.parent.insertElem(gtx, e, ev.Modifiers == key.ModShift)
		case e.Event(gtx, &ev, 0, key.ModShift, "⌫", "⌦"):
			e.parent.deleteElem(gtx, e, ev.Name == "⌫" && ev.Modifiers == 0)
		}
	}

nevents:
	for {
		var ev key.Event
		switch {
		default:
			break nevents
		case e.named.FocusEvent(gtx):
		case e.named.Event(gtx, &ev, 0, 0, "→"):
			e.Focus(gtx)
		case e.named.Event(gtx, &ev, 0, 0, "←"):
			e.parent.focusNext(gtx, e, false)
		case e.named.Event(gtx, &ev, 0, key.ModShift, "↑", "↓"):
			if ev.Modifiers == key.ModShift {
				e.parent.swap(e, ev.Name == "↓")
			} else {
				e.parent.focusNext(gtx, e, ev.Name == "↓")
			}
		}
	}

	for {
		ev, ok := gtx.Event(key.Filter{Focus: &e.named.Editor, Name: "⎋"})
		if !ok {
			break
		}
		if ev.(key.Event).State == key.Press {
			e.named.SetText(e.typ.Name)
			e.Focus(gtx)
		}
	}

	for {
		ev, ok := e.named.Update(gtx)
		if !ok {
			break
		}
		switch ev := ev.(type) {
		case widget.SubmitEvent:
			if validName(ev.Text) {
				e.typ.Name = ev.Text
				e.Focus(gtx)
			}
		}
	}

	if e.Focused(gtx) && e.typ.Name == "" {
		e.parent.deleteElem(gtx, e, true)
	}
}

func (e *EnumElemTypeEditor) LayoutName(gtx C) int {
	e.nameRec = Record(gtx, e.named.Layout)
	return e.nameRec.Dims.Size.X
}

func (e *EnumElemTypeEditor) Layout(gtx C, nameWidth int) D {
	e.Update(gtx)

	indent := unit.Dp(float32(nameWidth-e.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 4}.Layout),
		layout.Rigid(func(gtx C) D {
			return e.KeyFocus.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(layout.Spacer{Width: indent}.Layout),
					layout.Rigid(e.nameRec.Layout),
					layout.Rigid(layout.Spacer{Width: 8}.Layout),
					layout.Rigid(e.typed.Layout),
				)
			})
		}),
	)
}
