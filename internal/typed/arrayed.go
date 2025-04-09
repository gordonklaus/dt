package typed

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type ArrayTypeEditor struct {
	parent *TypeEditor
	typ    *types.ArrayType
	elem   *TypeEditor
}

func NewArrayTypeEditor(parent *TypeEditor, typ *types.ArrayType, core *Core) *ArrayTypeEditor {
	a := &ArrayTypeEditor{
		parent: parent,
		typ:    typ,
	}
	a.elem = NewTypeEditor(a, &typ.Elem, core)
	return a
}

func (a *ArrayTypeEditor) Type() types.Type { return a.typ }

func (a *ArrayTypeEditor) CreateNext(gtx C, after *TypeEditor) {
	if after == nil {
		a.elem.Edit(gtx)
	} else {
		a.parent.CreateNext(gtx, a)
	}
}

func (a *ArrayTypeEditor) Focus(gtx C) {
	a.elem.Focus(gtx)
}

func (a *ArrayTypeEditor) Layout(gtx C) D {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case a.elem.FocusEvent(gtx):
		case a.elem.Event(gtx, &e, 0, 0, "←"):
			a.parent.Focus(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "[]").Layout),
		layout.Rigid(a.elem.Layout),
	)
}
