package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type MapTypeEditor struct {
	typ        *types.MapType
	key, value *TypeEditor
}

func NewMapTypeEditor(typ *types.MapType, loader *types.Loader) *MapTypeEditor {
	return &MapTypeEditor{
		typ:   typ,
		key:   NewMapKeyTypeEditor(&typ.Key, loader),
		value: NewTypeEditor(&typ.Value, loader),
	}
}

func (a *MapTypeEditor) Type() types.Type { return a.typ }

func (a *MapTypeEditor) Layout(gtx C) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "map[").Layout),
		layout.Rigid(a.key.Layout),
		layout.Rigid(material.Body1(theme, "]").Layout),
		layout.Rigid(a.value.Layout),
	)
}
