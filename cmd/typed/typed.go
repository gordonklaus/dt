package main

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
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
	items      struct {
		int, strct widget.Clickable
	}

	ed typeEditor
}

type typeEditor interface {
	Type() types.Type
	Layout(gtx C) D
}

func NewTypeEditor(typ *types.Type) *TypeEditor {
	t := &TypeEditor{
		typ: typ,
	}
	t.menu.Options = []func(C) D{
		component.MenuItem(theme, &t.items.int, "int").Layout,
		component.MenuItem(theme, &t.items.strct, "struct").Layout,
	}

	switch typ := (*typ).(type) {
	case *types.IntType, *types.UintType:
		t.ed = NewIntTypeEditor(typ)
	case *types.StringType:
		t.ed = NewStringTypeEditor(typ)
	case *types.StructType:
		t.ed = NewStructTypeEditor(typ)
	}

	return t
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
			b.Size = unit.Dp(12)
			b.Inset = layout.UniformInset(unit.Dp(2))
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
	switch {
	case t.items.int.Clicked():
		return NewIntTypeEditor(&types.IntType{Size: 64})
	case t.items.strct.Clicked():
		return NewStructTypeEditor(&types.StructType{})
	}
	return nil
}
