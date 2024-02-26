package typed

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
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
	return o
}

func (o *OptionTypeEditor) Type() types.Type { return o.typ }

func (o *OptionTypeEditor) Layout(gtx C) D {
	for _, e := range o.KeyFocus.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := o.elem.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			o.parent.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			o.elem.Edit(gtx)
		}
	}

	if o.typ.Elem == nil {
		if o.Focused(gtx) {
			*o.parent.typ = nil
			o.parent.ed = nil
			o.parent.Edit(gtx)
		} else {
			o.elem.Edit(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "?").Layout),
		layout.Rigid(func(gtx C) D {
			return o.KeyFocus.Layout(gtx, o.elem.Layout)
		}),
	)
}
