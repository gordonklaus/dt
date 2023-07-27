package main

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type TypeEditor struct {
	parent Focuser
	typ    *types.Type
	loader *types.Loader

	showMenu    bool
	menu        layout.List
	items       []*typeMenuItem
	focusedItem int

	ed typeEditor
}

type typeMenuItem struct {
	c   widget.Clickable
	txt string
	new func() types.Type
}

type typeEditor interface {
	Type() types.Type
	Layout(gtx C) D
}

func NewTypeNameTypeEditor(parent Focuser, typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(parent, typ, loader)
	t.items = []*typeMenuItem{
		{txt: "struct", new: func() types.Type { return &types.StructType{} }},
		{txt: "enum", new: func() types.Type { return &types.EnumType{} }},
	}
	return t
}

func NewMapKeyTypeEditor(parent Focuser, typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(parent, typ, loader)
	t.items = []*typeMenuItem{
		{txt: "int", new: func() types.Type { return &types.IntType{} }},
		{txt: "uint", new: func() types.Type { return &types.IntType{Unsigned: true} }},
		{txt: "float32", new: func() types.Type { return &types.FloatType{Size: 32} }},
		{txt: "float64", new: func() types.Type { return &types.FloatType{Size: 64} }},
		{txt: "string", new: func() types.Type { return &types.StringType{} }},
	}
	return t
}

func NewTypeEditor(parent Focuser, typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(parent, typ, loader)
	t.items = []*typeMenuItem{
		{txt: "bool", new: func() types.Type { return &types.BoolType{} }},
		{txt: "int", new: func() types.Type { return &types.IntType{} }},
		{txt: "uint", new: func() types.Type { return &types.IntType{Unsigned: true} }},
		{txt: "float32", new: func() types.Type { return &types.FloatType{Size: 32} }},
		{txt: "float64", new: func() types.Type { return &types.FloatType{Size: 64} }},
		{txt: "string", new: func() types.Type { return &types.StringType{} }},
		{txt: "array", new: func() types.Type { return &types.ArrayType{} }},
		{txt: "map", new: func() types.Type { return &types.MapType{} }},
		{txt: "option", new: func() types.Type { return &types.OptionType{} }},
	}
	pkg, _ := loader.Load(&types.PackageID_Current{})
	for _, n := range pkg.Types {
		n := n
		t.items = append(t.items, &typeMenuItem{txt: n.Name, new: func() types.Type {
			return &types.NamedType{Package: &types.PackageID_Current{}, Name: n.Name, Type: n.Type}
		}})
	}
	return t
}

func newTypeEditor(parent Focuser, typ *types.Type, loader *types.Loader) *TypeEditor {
	t := &TypeEditor{
		parent: parent,
		typ:    typ,
		loader: loader,
		menu:   layout.List{Axis: layout.Vertical},
	}
	t.ed = t.newEditor(*typ)
	return t
}

func (t *TypeEditor) newEditor(typ types.Type) typeEditor {
	switch typ := typ.(type) {
	case nil:
		return nil
	case *types.BoolType:
		return NewBoolTypeEditor(typ)
	case *types.IntType:
		return NewIntTypeEditor(typ)
	case *types.FloatType:
		return NewFloatTypeEditor(typ)
	case *types.StringType:
		return NewStringTypeEditor(typ)
	case *types.StructType:
		return NewStructTypeEditor(t.parent, typ, t.loader)
	case *types.EnumType:
		return NewEnumTypeEditor(t.parent, typ, t.loader)
	case *types.ArrayType:
		return NewArrayTypeEditor(t, typ, t.loader)
	case *types.MapType:
		return NewMapTypeEditor(t, typ, t.loader)
	case *types.OptionType:
		return NewOptionTypeEditor(t, typ, t.loader)
	case *types.NamedType:
		return NewNamedTypeEditor(typ)
	}
	panic("unreached")
}

func (t *TypeEditor) Edit() { t.showMenu = true }

func (t *TypeEditor) Layout(gtx C) D {
	if ed := t.itemClicked(); ed != nil {
		t.showMenu = false
		*t.typ = ed.Type()
		t.ed = ed
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	return layout.Stack{Alignment: layout.SE}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			if t.ed != nil {
				return t.ed.Layout(gtx)
			}
			lbl := material.Body1(theme, "type")
			lbl.Font.Style = font.Italic
			return lbl.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			if t.showMenu {
				t.layoutMenu(gtx)
			}
			return D{}
		}),
	)
}

func (t *TypeEditor) layoutMenu(gtx C) {
	m := op.Record(gtx.Ops)
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			r := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(4))
			paint.FillShape(gtx.Ops, theme.Bg, r.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, theme.Fg,
				clip.Stroke{
					Path:  r.Path(gtx.Ops),
					Width: 1,
				}.Op(),
			)
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
				return t.menu.Layout(gtx, len(t.items), t.layoutMenuItem)
			})
		}),
	)
	op.Defer(gtx.Ops, m.Stop())
}

func (t *TypeEditor) layoutMenuItem(gtx C, i int) D {
	it := t.items[i]

	if i == t.focusedItem {
		for _, e := range gtx.Events(t) {
			switch e := e.(type) {
			case key.Event:
				if e.State == key.Press {
					switch e.Name {
					case "↑":
						if t.focusedItem > 0 {
							t.focusedItem--
						}
					case "↓":
						if t.focusedItem < len(t.items)-1 {
							t.focusedItem++
						}
					case "⏎", "⌤":
						t.showMenu = false
						t.ed = t.newEditor(it.new())
						*t.typ = t.ed.Type()
						if ed, ok := t.ed.(Focuser); ok {
							ed.Focus()
						} else {
							t.parent.Focus()
						}
					case "⎋":
						t.showMenu = false
						t.parent.Focus()
					}
				}
			}
		}

		key.FocusOp{Tag: t}.Add(gtx.Ops)
		key.InputOp{
			Tag:  t,
			Keys: "↑|↓|⏎|⌤|⎋",
		}.Add(gtx.Ops)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if i != t.focusedItem {
				return D{}
			}
			paint.FillShape(gtx.Ops, color.NRGBA{A: 64},
				clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(4)).Op(gtx.Ops))
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
				return material.Body1(theme, it.txt).Layout(gtx)
			})
		}),
	)
}

func (t *TypeEditor) itemClicked() typeEditor {
	for _, i := range t.items {
		if i.c.Clicked() {
			return t.newEditor(i.new())
		}
	}
	return nil
}
