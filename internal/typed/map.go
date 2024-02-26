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

	KeyFocus
	focusValue KeyFocus
}

func NewMapTypeEditor(parent *TypeEditor, typ *types.MapType, loader *types.Loader) *MapTypeEditor {
	m := &MapTypeEditor{
		parent: parent,
		typ:    typ,
	}
	m.key = NewMapKeyTypeEditor(m, &typ.Key, loader)
	m.value = NewTypeEditor(&m.focusValue, &typ.Value, loader)
	return m
}

func (m *MapTypeEditor) Type() types.Type { return m.typ }

func (m *MapTypeEditor) Layout(gtx C) D {
	for _, e := range m.KeyFocus.Events(gtx) {
		switch e.Name {
		case "→":
			m.focusValue.Focus(gtx)
		case "←":
			m.parent.parent.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			m.key.Edit(gtx)
		}
	}

	for _, e := range m.focusValue.Events(gtx) {
		switch e.Name {
		case "→":
			if ed, ok := m.value.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case "←":
			m.Focus(gtx)
		case "⏎", "⌤", "⌫", "⌦":
			m.value.Edit(gtx)
		}
	}

	if m.typ.Key == nil {
		if m.Focused(gtx) {
			*m.parent.typ = nil
			m.parent.ed = nil
			m.parent.Edit(gtx)
		} else {
			m.key.Edit(gtx)
		}
	} else if m.typ.Value == nil && m.Focused(gtx) {
		m.value.Edit(gtx)
	}

	if m.focusValue.Focused(gtx) && m.typ.Value == nil {
		*m.key.typ = nil
		m.key.ed = nil
		m.key.Edit(gtx)
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "map[").Layout),
		layout.Rigid(func(gtx C) D {
			return m.KeyFocus.Layout(gtx, m.key.Layout)
		}),
		layout.Rigid(material.Body1(theme, "]").Layout),
		layout.Rigid(func(gtx C) D {
			return m.focusValue.Layout(gtx, m.value.Layout)
		}),
	)
}
