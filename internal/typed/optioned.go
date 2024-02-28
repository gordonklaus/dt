package typed

import (
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type OptionTypeEditor struct {
	parent *TypeEditor
	typ    *types.OptionType
	elem   *TypeEditor
}

func NewOptionTypeEditor(parent *TypeEditor, typ *types.OptionType, loader *types.Loader) *OptionTypeEditor {
	o := &OptionTypeEditor{
		parent: parent,
		typ:    typ,
	}
	o.elem = NewTypeEditor(&typ.Elem, loader)
	return o
}

func (o *OptionTypeEditor) Type() types.Type { return o.typ }

func (o *OptionTypeEditor) Focus(gtx C) {
	o.elem.Focus(gtx)
}

func (o *OptionTypeEditor) Layout(gtx C) D {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case o.elem.FocusEvent(gtx):
		case o.elem.Event(gtx, &e, 0, 0, "←"):
			o.parent.Focus(gtx)
		}
	}

	if o.typ.Elem == nil {
		if o.elem.Focused(gtx) {
			*o.parent.typ = nil
			o.parent.ed = nil
			o.parent.Edit(gtx)
		} else {
			o.elem.Edit(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "?").Layout),
		layout.Rigid(o.elem.Layout),
	)
}
