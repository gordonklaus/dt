package main

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
	"golang.org/x/exp/slices"
)

type StructTypeEditor struct {
	typ    *types.StructType
	loader *types.Loader
	fields []*StructFieldTypeEditor
}

func NewStructTypeEditor(typ *types.StructType, loader *types.Loader) *StructTypeEditor {
	s := &StructTypeEditor{
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

func (s *StructTypeEditor) insertField(f *StructFieldTypeEditor) {
	for i, f2 := range s.fields {
		if f2 == f {
			field := &types.StructFieldType{}
			s.typ.Fields = slices.Insert(s.typ.Fields, i+1, field)
			s.fields = slices.Insert(s.fields, i+1, NewStructFieldTypeEditor(s, field, s.loader))
			break
		}
	}
}

func (s *StructTypeEditor) deleteField(f *StructFieldTypeEditor) {
	for i, f2 := range s.fields {
		if f2 == f {
			s.typ.Fields = slices.Delete(s.typ.Fields, i, i+1)
			s.fields = slices.Delete(s.fields, i, i+1)
			break
		}
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
		layout.Rigid(material.Body1(theme, "struct ").Layout),
		layout.Rigid(func(gtx C) D {
			width := gtx.Dp(16)
			height := fieldsRec.Dims.Size.Y + gtx.Dp(8)
			w := float32(width)
			h2 := float32(height) / 2
			path := clip.Path{}
			path.Begin(gtx.Ops)
			path.Move(f32.Pt(w, 0))
			path.Cube(f32.Pt(-w, 0), f32.Pt(0, h2), f32.Pt(-w, h2))
			path.Cube(f32.Pt(w, 0), f32.Pt(0, h2), f32.Pt(w, h2))
			paint.FillShape(gtx.Ops, theme.Fg, clip.Stroke{
				Path:  path.End(),
				Width: float32(gtx.Dp(1)),
			}.Op())
			return D{Size: image.Pt(width, height)}
		}),
		layout.Rigid(fieldsRec.Layout),
	)
}

type StructFieldTypeEditor struct {
	parent *StructTypeEditor
	typ    *types.StructFieldType
	named  widget.Editor
	typed  *TypeEditor

	nameRec Recording
}

func NewStructFieldTypeEditor(parent *StructTypeEditor, typ *types.StructFieldType, loader *types.Loader) *StructFieldTypeEditor {
	f := &StructFieldTypeEditor{
		parent: parent,
		typ:    typ,
		named: widget.Editor{
			Alignment:  text.End,
			SingleLine: true,
			Submit:     true,
		},
		typed: NewTypeEditor(&typ.Type, loader),
	}
	f.named.SetText(typ.Name)
	return f
}

func (f *StructFieldTypeEditor) LayoutName(gtx C) int {
	f.nameRec = Record(gtx, material.Editor(theme, &f.named, "").Layout)
	return f.nameRec.Dims.Size.X
}

func (f *StructFieldTypeEditor) Layout(gtx C, nameWidth int) D {
	for _, e := range f.named.Events() {
		switch e := e.(type) {
		case widget.ChangeEvent:
			f.typ.Name = f.named.Text()
		case widget.SubmitEvent:
			if e.Text == "" {
				f.parent.deleteField(f)
			} else {
				f.parent.insertField(f)
			}
		}
	}

	indent := unit.Dp(float32(nameWidth-f.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: indent}.Layout),
		layout.Rigid(f.nameRec.Layout),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(f.typed.Layout),
	)
}
