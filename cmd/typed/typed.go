package main

import (
	"math"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/gordonklaus/data"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type TypeEditor struct {
	typ *data.Type

	menuButton widget.Clickable
	showMenu   bool
	menu       component.MenuState
	items      struct {
		int, strct widget.Clickable
	}

	ed typeEditor
}

type typeEditor interface {
	Type() data.Type
	Layout(gtx C) D
}

func NewTypeEditor(typ *data.Type) *TypeEditor {
	t := &TypeEditor{
		typ: typ,
	}
	t.menu.Options = []func(C) D{
		component.MenuItem(theme, &t.items.int, "int").Layout,
		component.MenuItem(theme, &t.items.strct, "struct").Layout,
	}
	return t
}

var typeMenuIcon, _ = widget.NewIcon(icons.ActionSettings)

func (t *TypeEditor) Layout(gtx C) D {
	if t.menuButton.Clicked() {
		t.showMenu = !t.showMenu
	}

	if ed := t.itemClicked(); ed != nil {
		t.showMenu = false
		*t.typ = ed.Type()
		t.ed = ed
	}

	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(64))
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.IconButton(theme, &t.menuButton, typeMenuIcon, "type").Layout),
		layout.Rigid(func(gtx C) D {
			if t.showMenu {
				m := op.Record(gtx.Ops)
				gtx.Constraints.Max.Y = math.MaxInt
				component.Menu(theme, &t.menu).Layout(gtx)
				op.Defer(gtx.Ops, m.Stop())
			}
			return D{}
		}),
		layout.Flexed(1, func(gtx C) D {
			if t.ed != nil {
				return t.ed.Layout(gtx)
			}
			return D{}
		}),
	)
}

func (t *TypeEditor) itemClicked() typeEditor {
	switch {
	case t.items.int.Clicked():
		return NewTypeName(&data.NamedType{Name: "int"})
	case t.items.strct.Clicked():
		return NewStructTypeEditor(&data.StructType{})
	}
	return nil
}
