package typed

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type TypeEditor struct {
	typ    *types.Type
	loader *types.Loader

	KeyFocus
	menuFocus   KeyFocus
	menu        layout.List
	items       []*typeMenuItem
	focusedItem int

	ed typeEditor
}

type typeMenuItem struct {
	txt string
	new func() types.Type
}

type typeEditor interface {
	Type() types.Type
	Layout(gtx C) D
}

func NewTypeNameTypeEditor(typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(typ, loader)
	t.items = []*typeMenuItem{
		{txt: "struct", new: func() types.Type { return &types.StructType{} }},
		{txt: "enum", new: func() types.Type { return &types.EnumType{} }},
	}
	return t
}

func NewMapKeyTypeEditor(typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(typ, loader)
	t.items = []*typeMenuItem{
		{txt: "int", new: func() types.Type { return &types.IntType{} }},
		{txt: "uint", new: func() types.Type { return &types.IntType{Unsigned: true} }},
		{txt: "float32", new: func() types.Type { return &types.FloatType{Size: 32} }},
		{txt: "float64", new: func() types.Type { return &types.FloatType{Size: 64} }},
		{txt: "string", new: func() types.Type { return &types.StringType{} }},
	}
	return t
}

func NewTypeEditor(typ *types.Type, loader *types.Loader) *TypeEditor {
	t := newTypeEditor(typ, loader)
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
			return &types.NamedType{Package: &types.PackageID_Current{}, TypeName: n}
		}})
	}
	return t
}

func newTypeEditor(typ *types.Type, loader *types.Loader) *TypeEditor {
	t := &TypeEditor{
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
		return NewStructTypeEditor(t, typ, t.loader)
	case *types.EnumType:
		return NewEnumTypeEditor(t, typ, t.loader)
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

func (t *TypeEditor) Edit(gtx C) {
	t.menuFocus.Focus(gtx)
}

func (t *TypeEditor) Layout(gtx C) D {
	t.updateMenu(gtx)

	return layout.Stack{Alignment: layout.SE}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			if t.ed != nil {
				return t.KeyFocus.Layout(gtx, t.ed.Layout)
			}
			lbl := material.Body1(theme, "type")
			lbl.Font.Style = font.Italic
			return lbl.Layout(gtx)
		}),
		layout.Stacked(t.layoutMenu),
	)
}

func (t *TypeEditor) updateMenu(gtx C) {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case t.menuFocus.FocusEvent(gtx):
		case t.menuFocus.Event(gtx, &e, 0, 0, "↑"):
			if t.focusedItem > 0 {
				t.focusedItem--
			}
		case t.menuFocus.Event(gtx, &e, 0, 0, "↓"):
			if t.focusedItem < len(t.items)-1 {
				t.focusedItem++
			}
		case t.menuFocus.Event(gtx, &e, 0, 0, "⏎", "⌤"):
			t.ed = t.newEditor(t.items[t.focusedItem].new())
			*t.typ = t.ed.Type()
			if ed, ok := t.ed.(Focuser); ok {
				ed.Focus(gtx)
			} else {
				t.Focus(gtx)
			}
		case t.menuFocus.Event(gtx, &e, 0, 0, "⎋"):
			t.Focus(gtx)
		}
	}
}

func (t *TypeEditor) layoutMenu(gtx C) D {
	if t.menuFocus.Focused(gtx) {
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
					return t.menu.Layout(gtx, len(t.items), func(gtx C, i int) D {
						return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
							return t.layoutMenuItem(gtx, i)
						})
					})
				})
			}),
		)
		op.Defer(gtx.Ops, m.Stop())
	}
	return D{}
}

func (t *TypeEditor) layoutMenuItem(gtx C, i int) D {
	w := material.Body1(theme, t.items[i].txt).Layout
	if i == t.focusedItem {
		return t.menuFocus.Layout(gtx, w)
	}
	return w(gtx)
}
