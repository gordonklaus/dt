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

type EnumTypeEditor struct {
	typ   *types.EnumType
	elems []*EnumElemTypeEditor
}

func NewEnumTypeEditor(typ *types.EnumType) *EnumTypeEditor {
	s := &EnumTypeEditor{
		typ:   typ,
		elems: make([]*EnumElemTypeEditor, len(typ.Elems)),
	}
	for i, f := range typ.Elems {
		s.elems[i] = NewEnumElemTypeEditor(s, f)
	}
	return s
}

func (e *EnumTypeEditor) Type() types.Type { return e.typ }

func (e *EnumTypeEditor) insertElem(f *EnumElemTypeEditor) {
	for i, f2 := range e.elems {
		if f2 == f {
			elem := &types.EnumElemType{}
			e.typ.Elems = slices.Insert(e.typ.Elems, i+1, elem)
			e.elems = slices.Insert(e.elems, i+1, NewEnumElemTypeEditor(e, elem))
			break
		}
	}
}

func (e *EnumTypeEditor) deleteElem(f *EnumElemTypeEditor) {
	for i, f2 := range e.elems {
		if f2 == f {
			e.typ.Elems = slices.Delete(e.typ.Elems, i, i+1)
			e.elems = slices.Delete(e.elems, i, i+1)
			break
		}
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

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "struct ").Layout),
		layout.Rigid(func(gtx C) D {
			width := gtx.Dp(16)
			height := elemsRec.Dims.Size.Y + gtx.Dp(8)
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
		layout.Rigid(elemsRec.Layout),
	)
}

type EnumElemTypeEditor struct {
	parent *EnumTypeEditor
	typ    *types.EnumElemType
	named  widget.Editor
	typed  *TypeEditor

	nameRec Recording
}

func NewEnumElemTypeEditor(parent *EnumTypeEditor, typ *types.EnumElemType) *EnumElemTypeEditor {
	f := &EnumElemTypeEditor{
		parent: parent,
		typ:    typ,
		named: widget.Editor{
			Alignment:  text.End,
			SingleLine: true,
			Submit:     true,
		},
		typed: NewTypeEditor(&typ.Type),
	}
	f.named.SetText(typ.Name)
	return f
}

func (e *EnumElemTypeEditor) LayoutName(gtx C) int {
	e.nameRec = Record(gtx, material.Editor(theme, &e.named, "").Layout)
	return e.nameRec.Dims.Size.X
}

func (e *EnumElemTypeEditor) Layout(gtx C, nameWidth int) D {
	for _, ev := range e.named.Events() {
		switch ev := ev.(type) {
		case widget.ChangeEvent:
			e.typ.Name = e.named.Text()
		case widget.SubmitEvent:
			if ev.Text == "" {
				e.parent.deleteElem(e)
			} else {
				e.parent.insertElem(e)
			}
		}
	}

	indent := unit.Dp(float32(nameWidth-e.nameRec.Dims.Size.X) / gtx.Metric.PxPerDp)
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: indent}.Layout),
		layout.Rigid(e.nameRec.Layout),
		layout.Rigid(layout.Spacer{Width: 8}.Layout),
		layout.Rigid(e.typed.Layout),
	)
}
