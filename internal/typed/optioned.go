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
	for _, e := range o.elem.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := o.elem.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			o.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			o.elem.Edit(gtx)
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
