package typed

import (
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type NamedTypeEditor struct {
	typ *types.NamedType
}

func NewNamedTypeEditor(typ *types.NamedType) *NamedTypeEditor {
	return &NamedTypeEditor{
		typ: typ,
	}
}

func (n *NamedTypeEditor) Type() types.Type { return n.typ }

func (n *NamedTypeEditor) Layout(gtx C) D {
	txt := n.typ.TypeName.Name
	if tn, ok := n.typ.TypeName.Parent.(*types.TypeName); ok {
		txt = tn.Name + "." + txt
	}
	return material.Body1(theme, txt).Layout(gtx)
}
