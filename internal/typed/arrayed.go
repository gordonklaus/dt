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
}

func NewArrayTypeEditor(parent *TypeEditor, typ *types.ArrayType, loader *types.Loader) *ArrayTypeEditor {
	a := &ArrayTypeEditor{
		parent: parent,
		typ:    typ,
	}
	a.elem = NewTypeEditor(&typ.Elem, loader)
	return a
}

func (a *ArrayTypeEditor) Type() types.Type { return a.typ }

func (a *ArrayTypeEditor) Focus(gtx C) {
	a.elem.Focus(gtx)
}

func (a *ArrayTypeEditor) Layout(gtx C) D {
	for _, e := range a.elem.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := a.elem.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			a.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			a.elem.Edit(gtx)
		}
	}

	if a.typ.Elem == nil {
		if a.elem.Focused(gtx) {
			*a.parent.typ = nil
			a.parent.ed = nil
			a.parent.Edit(gtx)
		} else {
			a.elem.Edit(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "[]").Layout),
		layout.Rigid(a.elem.Layout),
	)
}
