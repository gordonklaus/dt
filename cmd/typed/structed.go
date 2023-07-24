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

func (s *StructTypeEditor) Focus() {
	if len(s.fields) == 0 {
		s.insertField(nil, false)
	} else {
		s.fields[(len(s.fields)-1)/2].Focus()
	}
}

func (s *StructTypeEditor) focusNext(f *StructFieldTypeEditor, next bool) {
	i := slices.Index(s.fields, f) - 1
	if next {
		i += 2
	}
	if i < 0 {
		s.parent.Focus()
	} else if i < len(s.fields) {
		s.fields[i].Focus()
	} else {
		switch p := s.parent.(type) {
		case *EnumElemTypeEditor:
			p.parent.focusNext(p, next)
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

func (s *StructTypeEditor) deleteField(f *StructFieldTypeEditor, back bool) {
	i := slices.Index(s.fields, f)
	s.typ.Fields = slices.Delete(s.typ.Fields, i, i+1)
	s.fields = slices.Delete(s.fields, i, i+1)
	if i > 0 && (back || i >= len(s.fields)) {
		i--
	}
	if i < len(s.fields) {
		s.fields[i].Focus()
	} else {
		s.parent.Focus()
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
	for _, e := range f.KeyFocus.Events(gtx, "←|→|↑|↓|(Shift)-[⏎,⌤,⌫,⌦]") {
		switch e.Name {
		case "←":
			f.focusNamed.Focus()
		case "→":
			f.focusTyped.Focus()
		case "↑":
			f.parent.focusNext(f, false)
		case "↓":
			f.parent.focusNext(f, true)
		case "⏎", "⌤":
			f.parent.insertField(f, e.Modifiers == key.ModShift)
		case "⌫", "⌦":
			f.parent.deleteField(f, (e.Name == "⌦") == (e.Modifiers == key.ModShift))
		}
	}

	for _, e := range f.focusNamed.Events(gtx, "←|→|⏎|⌤|⌫|⌦|⎋") {
		switch e.Name {
		case "→":
			f.Focus()
		case "←":
			f.parent.parent.Focus()
		case "⏎", "⌤", "⌫", "⌦":
			f.named.SetCaret(f.named.Len(), f.named.Len())
			f.named.Focus()
		case "⎋":
			f.named.SetText(f.typ.Name)
			f.Focus()
		}
	}

	for _, e := range f.focusTyped.Events(gtx, "←|→|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			if ed, ok := f.typed.ed.(Focuser); ok {
				ed.Focus()
			}
		case "←":
			f.Focus()
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
					f.Focus()
				}
			}
		}
	}

	if f.Focused() && (f.typ.Name == "" || f.typ.Type == nil) {
		f.parent.deleteField(f, true)
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
