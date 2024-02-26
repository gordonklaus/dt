package typed

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
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
	return a
}

func (a *ArrayTypeEditor) Type() types.Type { return a.typ }

func (a *ArrayTypeEditor) Layout(gtx C) D {
	for _, e := range a.KeyFocus.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := a.elem.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			a.parent.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			a.elem.Edit(gtx)
		}
	}

	if a.typ.Elem == nil {
		if a.Focused(gtx) {
			*a.parent.typ = nil
			a.parent.ed = nil
			a.parent.Edit(gtx)
		} else {
			a.elem.Edit(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "[]").Layout),
		layout.Rigid(func(gtx C) D {
			return a.KeyFocus.Layout(gtx, a.elem.Layout)
		}),
	)
}
