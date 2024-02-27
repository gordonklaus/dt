package typed

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type MapTypeEditor struct {
	parent     *TypeEditor
	typ        *types.MapType
	key, value *TypeEditor
}

func NewMapTypeEditor(parent *TypeEditor, typ *types.MapType, loader *types.Loader) *MapTypeEditor {
	m := &MapTypeEditor{
		parent: parent,
		typ:    typ,
	}
	m.key = NewMapKeyTypeEditor(&typ.Key, loader)
	m.value = NewTypeEditor(&typ.Value, loader)
	return m
}

func (m *MapTypeEditor) Type() types.Type { return m.typ }

func (m *MapTypeEditor) Focus(gtx C) {
	m.key.Focus(gtx)
}

func (m *MapTypeEditor) Layout(gtx C) D {
	for _, e := range m.key.Events(gtx) {
		switch e.Name {
		case "→":
			m.value.Focus(gtx)
		case "←":
			m.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			m.key.Edit(gtx)
		}
	}

	for _, e := range m.value.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := m.value.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			m.key.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			m.value.Edit(gtx)
		}
	}

	if m.typ.Key == nil {
		if m.key.Focused(gtx) {
			*m.parent.typ = nil
			m.parent.ed = nil
			m.parent.Edit(gtx)
		} else {
			m.key.Edit(gtx)
		}
	} else if m.typ.Value == nil {
		if m.key.Focused(gtx) {
			m.value.Edit(gtx)
		} else if m.value.Focused(gtx) {
			*m.key.typ = nil
			m.key.ed = nil
			m.key.Edit(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "map[").Layout),
		layout.Rigid(m.key.Layout),
		layout.Rigid(material.Body1(theme, "]").Layout),
		layout.Rigid(m.value.Layout),
	)
}
