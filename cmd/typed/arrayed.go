package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type ArrayTypeEditor struct {
	parent *TypeEditor
	typ    *types.ArrayType
	elem   *TypeEditor

	KeyFocus
}

func NewArrayTypeEditor(parent *TypeEditor, typ *types.ArrayType, loader *types.Loader) *ArrayTypeEditor {
	a := &ArrayTypeEditor{
		parent: parent,
		typ:    typ,
	}
	a.elem = NewTypeEditor(a, &typ.Elem, loader)
	if typ.Elem == nil {
		a.elem.showMenu = true
	}
	return a
}

func (a *ArrayTypeEditor) Type() types.Type { return a.typ }

func (a *ArrayTypeEditor) Layout(gtx C) D {
	for _, e := range a.KeyFocus.Events(gtx, "←|→|↑|↓|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			if ed, ok := a.elem.ed.(interface{ Focus() }); ok {
				ed.Focus()
			}
		case "←":
			switch p := a.parent.parent.(type) {
			case *StructFieldTypeEditor:
				p.focusTyped.Focus()
			case interface{ Focus() }:
				p.Focus()
			}
			// case "↑":
			// 	a.parent.focusNext(false)
			// case "↓":
			// 	a.parent.focusNext(true)
		case "⏎", "⌤", "⌫", "⌦":
			a.elem.showMenu = true
		}
	}

	if a.Focused() && a.typ.Elem == nil {
		*a.parent.typ = nil
		a.parent.ed = nil
		a.parent.showMenu = true
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "[]").Layout),
		layout.Rigid(func(gtx C) D {
			return a.KeyFocus.Layout(gtx, a.elem.Layout)
		}),
	)
}
