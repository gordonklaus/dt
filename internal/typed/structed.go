package typed

import (
	"image"
	"slices"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/gordonklaus/dt/types"
)

type StructTypeEditor struct {
	*Core
	parent Focuser
	typ    *types.StructType
	fields []*StructFieldTypeEditor
}

func NewStructTypeEditor(parent Focuser, typ *types.StructType, core *Core) *StructTypeEditor {
	s := &StructTypeEditor{
		parent: parent,
		typ:    typ,
		Core:   core,
		fields: make([]*StructFieldTypeEditor, len(typ.Fields)),
	}
	for i, f := range typ.Fields {
		s.fields[i] = NewStructFieldTypeEditor(s, f, core)
	}
	return s
}

func (s *StructTypeEditor) Type() types.Type { return s.typ }

func (s *StructTypeEditor) CreateNext(gtx C, after *TypeEditor) { s.Focus(gtx) }

func (s *StructTypeEditor) Focus(gtx C) {
	if len(s.fields) == 0 {
		s.insertField(gtx, nil, false)
	} else {
		s.fields[(len(s.fields)-1)/2].Focus(gtx)
	}
}

func (s *StructTypeEditor) focusNext(gtx C, f *StructFieldTypeEditor, next bool) {
	i := slices.Index(s.fields, f) - 1
	if next {
		i += 2
	}
	if i < 0 {
		s.parent.Focus(gtx)
	} else if i < len(s.fields) {
		s.fields[i].Focus(gtx)
	} else {
		switch p := s.parent.(type) {
		case *EnumElemTypeEditor:
			p.parent.focusNext(gtx, p, next)
		}
	}
}

func (s *StructTypeEditor) insertField(gtx C, f *StructFieldTypeEditor, before bool) {
	i := slices.Index(s.fields, f)
	if !before {
		i++
	}
	field := &types.StructFieldType{ID: nextID(s.Pkg)}
	switch p := s.parent.(type) {
	case *TypeEditor:
		field.Parent = p.typeName
	case *EnumElemTypeEditor:
		field.Parent = p.typ
	}
	s.typ.Fields = slices.Insert(s.typ.Fields, i, field)
	s.Pkg.TypesByID[field.ID] = field
	s.fields = slices.Insert(s.fields, i, NewStructFieldTypeEditor(s, field, s.Core))
	s.fields[i].named.Edit(gtx)
}

func (s *StructTypeEditor) swap(f *StructFieldTypeEditor, next bool) {
	i := slices.Index(s.fields, f)
	if next && i == len(s.fields)-1 || !next && i == 0 {
		return
	}
	if !next {
		i--
	}
	s.typ.Fields[i], s.typ.Fields[i+1] = s.typ.Fields[i+1], s.typ.Fields[i]
	s.fields[i], s.fields[i+1] = s.fields[i+1], s.fields[i]
}

func (s *StructTypeEditor) deleteField(gtx C, f *StructFieldTypeEditor, back bool) {
	i := slices.Index(s.fields, f)
	s.typ.Fields = slices.Delete(s.typ.Fields, i, i+1)
	delete(s.Pkg.TypesByID, f.typ.ID)
	s.fields = slices.Delete(s.fields, i, i+1)
	if i > 0 && (back || i >= len(s.fields)) {
		i--
	}
	if i < len(s.fields) {
		s.fields[i].Focus(gtx)
	} else {
		s.parent.Focus(gtx)
	}
}

func (s *StructTypeEditor) Layout(gtx C) D {
	// Iterate over a copy because s.fields may mutate during iteration.
	for _, f := range slices.Clone(s.fields) {
		f.Update(gtx)
	}

	maxFieldNameWidth := 0
	for _, f := range s.fields {
		if x := f.LayoutName(gtx); x > maxFieldNameWidth {
			maxFieldNameWidth = x
		}
	}
	fields := make([]layout.FlexChild, len(s.typ.Fields))
	for i, f := range s.fields {
		f := f
		fields[i] = layout.Rigid(func(gtx C) D {
			return f.Layout(gtx, maxFieldNameWidth)
		})
	}

	fieldsRec := Record(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, fields...)
	})

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if _, ok := s.parent.(*EnumElemTypeEditor); ok && len(fields) == 0 {
				return D{}
			}
			width := gtx.Dp(12)
			height := fieldsRec.Dims.Size.Y + gtx.Dp(8)
			w := float32(width)
			h2 := float32(height) / 2
			d := float32(gtx.Dp(4))
			path := clip.Path{}
			path.Begin(gtx.Ops)
			path.Move(f32.Pt(w, 0))
			path.Cube(f32.Pt(-w-d, 0), f32.Pt(d, h2), f32.Pt(-w, h2))
			path.Cube(f32.Pt(w+d, 0), f32.Pt(-d, h2), f32.Pt(w, h2))
			paint.FillShape(gtx.Ops, theme.Fg, clip.Stroke{
				Path:  path.End(),
				Width: float32(gtx.Dp(1)),
			}.Op())
			return D{Size: image.Pt(width, height)}
		}),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(fieldsRec.Layout),
	)
}

type StructFieldTypeEditor struct {
	parent *StructTypeEditor
	typ    *types.StructFieldType
	KeyFocus
	named nameEditor
	typed *TypeEditor

	nameRec Recording
}

func NewStructFieldTypeEditor(parent *StructTypeEditor, typ *types.StructFieldType, core *Core) *StructFieldTypeEditor {
	f := &StructFieldTypeEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	f.typed = NewTypeEditor(f, &typ.Type, core)
	f.named.SetText(typ.Name)
	return f
}

func (f *StructFieldTypeEditor) CreateNext(gtx C, after *TypeEditor) { f.Focus(gtx) }

func (f *StructFieldTypeEditor) Update(gtx C) {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case f.FocusEvent(gtx):
		case f.Event(gtx, &e, 0, 0, "←"):
			f.named.Focus(gtx)
		case f.Event(gtx, &e, 0, 0, "→"):
			f.typed.Focus(gtx)
		case f.Event(gtx, &e, 0, key.ModShift, "↑", "↓"):
			if e.Modifiers == key.ModShift {
				f.parent.swap(f, e.Name == "↓")
			} else {
				f.parent.focusNext(gtx, f, e.Name == "↓")
			}
		case f.Event(gtx, &e, 0, key.ModShift, "⏎", "⌤"):
			f.parent.insertField(gtx, f, e.Modifiers == key.ModShift)
		case f.Event(gtx, &e, 0, key.ModShift, "⌫", "⌦"):
			f.parent.deleteField(gtx, f, e.Name == "⌫" && e.Modifiers == 0)
		}
	}

nevents:
	for {
		var e key.Event
		switch {
		default:
			break nevents
		case f.named.FocusEvent(gtx):
		case f.named.Event(gtx, &e, 0, 0, "→"):
			f.Focus(gtx)
		case f.named.Event(gtx, &e, 0, 0, "←"):
			f.parent.parent.Focus(gtx)
		case f.named.Event(gtx, &e, 0, key.ModShift, "↑", "↓"):
			if e.Modifiers == key.ModShift {
				f.parent.swap(f, e.Name == "↓")
			} else {
				f.parent.focusNext(gtx, f, e.Name == "↓")
			}
		}
	}

	for {
		filters := []event.Filter{key.Filter{Focus: &f.named.Editor, Name: "⎋"}}
		if f.named.Text() == "" && f.typ.Type == nil {
			filters = append(filters, key.Filter{Focus: &f.named.Editor, Name: "←"})
		}
		e, ok := gtx.Event(filters...)
		if !ok {
			break
		}
		if e.(key.Event).State == key.Press {
			if f.typ.Name == "" || f.typ.Type == nil {
				f.parent.deleteField(gtx, f, true)
			} else {
				f.named.SetText(f.typ.Name)
				f.Focus(gtx)
			}
		}
	}

	for {
		e, ok := f.named.Update(gtx)
		if !ok {
			break
		}
		switch e := e.(type) {
		case widget.SubmitEvent:
			if validName(e.Text) {
				f.typ.Name = e.Text
				if f.typ.Type == nil {
					f.typed.Edit(gtx)
				} else {
					f.named.Focus(gtx)
				}
			}
		}
	}

tevents:
	for {
		var e key.Event
		switch {
		default:
			break tevents
		case f.typed.FocusEvent(gtx):
		case f.typed.Event(gtx, &e, 0, 0, "←"):
			f.Focus(gtx)
		case f.typed.Event(gtx, &e, 0, key.ModShift, "↑", "↓"):
			if e.Modifiers == key.ModShift {
				f.parent.swap(f, e.Name == "↓")
			} else {
				f.parent.focusNext(gtx, f, e.Name == "↓")
			}
		}
	}
}

func (f *StructFieldTypeEditor) LayoutName(gtx C) int {
	f.nameRec = Record(gtx, f.named.Layout)
	return f.nameRec.Dims.Size.X
}

func (f *StructFieldTypeEditor) Layout(gtx C, nameWidth int) D {
	f.Update(gtx)

	indent := unit.Dp(float32(nameWidth-f.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return f.KeyFocus.Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(layout.Spacer{Width: indent}.Layout),
			layout.Rigid(f.nameRec.Layout),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(f.typed.Layout),
		)
	})
}
