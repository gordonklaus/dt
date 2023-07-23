package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type OptionTypeEditor struct {
	parent *TypeEditor
	typ    *types.OptionType
	elem   *TypeEditor

	KeyFocus
}

func NewOptionTypeEditor(parent *TypeEditor, typ *types.OptionType, loader *types.Loader) *OptionTypeEditor {
	o := &OptionTypeEditor{
		parent: parent,
		typ:    typ,
	}
	o.elem = NewTypeEditor(o, &typ.Elem, loader)
	if typ.Elem == nil {
		o.elem.showMenu = true
	}
	return o
}

func (o *OptionTypeEditor) Type() types.Type { return o.typ }

func (o *OptionTypeEditor) Layout(gtx C) D {
	for _, e := range o.KeyFocus.Events(gtx, "←|→|↑|↓|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			if ed, ok := o.elem.ed.(interface{ Focus() }); ok {
				ed.Focus()
			}
		case "←":
			switch p := o.parent.parent.(type) {
			case *StructFieldTypeEditor:
				p.focusTyped.Focus()
			case interface{ Focus() }:
				p.Focus()
			}
			// case "↑":
			// 	o.parent.focusNext(false)
			// case "↓":
			// 	o.parent.focusNext(true)
		case "⏎", "⌤", "⌫", "⌦":
			o.elem.showMenu = true
		}
	}

	if o.Focused() && o.typ.Elem == nil {
		*o.parent.typ = nil
		o.parent.ed = nil
		o.parent.showMenu = true
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "?").Layout),
		layout.Rigid(func(gtx C) D {
			return o.KeyFocus.Layout(gtx, o.elem.Layout)
		}),
	)
}
