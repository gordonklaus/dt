package main

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/gordonklaus/data/types"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type TypeEditor struct {
	typ *types.Type

	menuButton widget.Clickable
	showMenu   bool
	menu       component.MenuState
	items      []*typeMenuItem

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

func NewTypeEditor(typ *types.Type) *TypeEditor {
	t := &TypeEditor{
		typ: typ,
		items: []*typeMenuItem{
			{txt: "bool", new: func() types.Type { return &types.BoolType{} }},
			{txt: "int", new: func() types.Type { return &types.IntType{Size: 64} }},
			{txt: "uint", new: func() types.Type { return &types.UintType{Size: 64} }},
			{txt: "float32", new: func() types.Type { return &types.Float32Type{} }},
			{txt: "float64", new: func() types.Type { return &types.Float64Type{} }},
			{txt: "string", new: func() types.Type { return &types.StringType{} }},
			{txt: "struct", new: func() types.Type { return &types.StructType{Fields: []*types.StructFieldType{{}}} }}, // Include a single field because StructTypeEditor has no way yet to add a first field.
			{txt: "enum", new: func() types.Type { return &types.EnumType{} }},
			{txt: "array", new: func() types.Type { return &types.ArrayType{} }},
			{txt: "option", new: func() types.Type { return &types.OptionType{} }},
		},
	}
	t.menu.Options = mapSlice(t.items, func(i *typeMenuItem) func(C) D {
		return component.MenuItem(theme, &i.c, i.txt).Layout
	})
	t.ed = t.newEditor(*typ)
	return t
}

func mapSlice[T, U any](t []T, f func(T) U) []U {
	u := make([]U, len(t))
	for i, t := range t {
		u[i] = f(t)
	}
	return u
}

func (t *TypeEditor) newEditor(typ types.Type) typeEditor {
	switch typ := typ.(type) {
	case nil:
		return nil
	case *types.BoolType, *types.Float32Type, *types.Float64Type, *types.StringType:
		return NewBasicTypeEditor(typ)
	case *types.IntType, *types.UintType:
		return NewIntTypeEditor(typ)
	case *types.StructType:
		return NewStructTypeEditor(typ)
	case *types.EnumType:
		return NewEnumTypeEditor(typ)
	case *types.ArrayType:
		return NewArrayTypeEditor(typ)
	case *types.OptionType:
		return NewOptionTypeEditor(typ)
	}
	panic("unreached")
}

var typeMenuIcon, _ = widget.NewIcon(icons.ImageEdit)

func (t *TypeEditor) Layout(gtx C) D {
	if t.menuButton.Clicked() {
		t.showMenu = !t.showMenu
	}

	if ed := t.itemClicked(); ed != nil {
		t.showMenu = false
		*t.typ = ed.Type()
		t.ed = ed
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			b := material.IconButton(theme, &t.menuButton, typeMenuIcon, "type")
			b.Size = 12
			b.Inset = layout.UniformInset(2)
			b.Background = theme.Bg
			b.Color = theme.Fg
			gtx.Constraints.Min.Y = 0
			return b.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if t.showMenu {
				m := op.Record(gtx.Ops)
				component.Menu(theme, &t.menu).Layout(gtx)
				op.Defer(gtx.Ops, m.Stop())
			}
			return D{}
		}),
		layout.Rigid(func(gtx C) D {
			if t.ed != nil {
				return t.ed.Layout(gtx)
			}
			return D{}
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
