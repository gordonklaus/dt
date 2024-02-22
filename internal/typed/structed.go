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

type StructTypeEditor struct {
	parent Focuser
	typ    *types.StructType
	loader *types.Loader
	fields []*StructFieldTypeEditor
}

func NewStructTypeEditor(parent Focuser, typ *types.StructType, loader *types.Loader) *StructTypeEditor {
	s := &StructTypeEditor{
		parent: parent,
		typ:    typ,
		loader: loader,
		fields: make([]*StructFieldTypeEditor, len(typ.Fields)),
	}
	for i, f := range typ.Fields {
		s.fields[i] = NewStructFieldTypeEditor(s, f, loader)
	}
	return s
}

func (s *StructTypeEditor) Type() types.Type { return s.typ }

func (s *StructTypeEditor) Focus(gtx C) {
	if len(s.fields) == 0 {
		s.insertField(nil, false)
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

func (s *StructTypeEditor) insertField(f *StructFieldTypeEditor, before bool) {
	i := slices.Index(s.fields, f)
	if !before {
		i++
	}
	field := &types.StructFieldType{}
	s.typ.Fields = slices.Insert(s.typ.Fields, i, field)
	s.fields = slices.Insert(s.fields, i, NewStructFieldTypeEditor(s, field, s.loader))
	s.fields[i].named.Focus()
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
	named  editor
	typed  *TypeEditor

	nameRec Recording

	KeyFocus
	focusNamed, focusTyped KeyFocus
}

func NewStructFieldTypeEditor(parent *StructTypeEditor, typ *types.StructFieldType, loader *types.Loader) *StructFieldTypeEditor {
	f := &StructFieldTypeEditor{
		parent: parent,
		typ:    typ,
		named:  newEditor(),
	}
	f.typed = NewTypeEditor(f, &typ.Type, loader)
	f.named.SetText(typ.Name)
	return f
}

func (f *StructFieldTypeEditor) LayoutName(gtx C) int {
	f.nameRec = Record(gtx, f.named.Layout)
	return f.nameRec.Dims.Size.X
}

func (f *StructFieldTypeEditor) Layout(gtx C, nameWidth int) D {
	for _, e := range f.KeyFocus.Events(gtx, "←|→|(Shift)-[↑,↓]|(Shift)-[⏎,⌤,⌫,⌦]") {
		switch e.Name {
		case "←":
			f.focusNamed.Focus(gtx)
		case "→":
			f.focusTyped.Focus(gtx)
		case "↑":
			if e.Modifiers == key.ModShift {
				f.parent.swap(f, false)
			} else {
				f.parent.focusNext(gtx, f, false)
			}
		case "↓":
			if e.Modifiers == key.ModShift {
				f.parent.swap(f, true)
			} else {
				f.parent.focusNext(gtx, f, true)
			}
		case "⏎", "⌤":
			f.parent.insertField(f, e.Modifiers == key.ModShift)
		case "⌫", "⌦":
			f.parent.deleteField(gtx, f, e.Name == "⌫" && e.Modifiers == 0)
		}
	}

	for _, e := range f.focusNamed.Events(gtx, "←|→|⏎|⌤|⌫|⌦|⎋") {
		switch e.Name {
		case "→":
			f.Focus(gtx)
		case "←":
			f.parent.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			f.named.SetCaret(f.named.Len(), f.named.Len())
			f.named.Focus()
		case "⎋":
			f.named.SetText(f.typ.Name)
			f.Focus(gtx)
		}
	}

	for _, e := range f.focusTyped.Events(gtx, "←|→|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			if ed, ok := f.typed.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			f.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			f.typed.Edit()
		}
	}

	for _, e := range f.named.Events() {
		switch e := e.(type) {
		case widget.SubmitEvent:
			if validName(e.Text) {
				f.typ.Name = e.Text
				if f.typ.Type == nil {
					f.typed.Edit()
				} else {
					f.Focus(gtx)
				}
			}
		}
	}

	if f.Focused() && (f.typ.Name == "" || f.typ.Type == nil) {
		f.parent.deleteField(gtx, f, true)
	}

	indent := unit.Dp(float32(nameWidth-f.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return f.KeyFocus.Layout(gtx, func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(layout.Spacer{Width: indent}.Layout),
			layout.Rigid(func(gtx C) D {
				return f.focusNamed.Layout(gtx, f.nameRec.Layout)
			}),
			layout.Rigid(layout.Spacer{Width: 8}.Layout),
			layout.Rigid(func(gtx C) D {
				return f.focusTyped.Layout(gtx, f.typed.Layout)
			}),
		)
	})
}
