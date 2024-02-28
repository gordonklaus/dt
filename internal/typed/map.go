package typed

import (
	"gioui.org/io/key"
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
kevents:
	for {
		var e key.Event
		switch {
		default:
			break kevents
		case m.key.FocusEvent(gtx):
		case m.key.Event(gtx, &e, 0, 0, "→"):
			m.value.Focus(gtx)
		case m.key.Event(gtx, &e, 0, 0, "←"):
			m.parent.Focus(gtx)
		case m.key.Event(gtx, &e, 0, 0, "⏎", "⌤", "⌫", "⌦"):
			m.key.Edit(gtx)
		}
	}

vevents:
	for {
		var e key.Event
		switch {
		default:
			break vevents
		case m.value.FocusEvent(gtx):
		case m.value.Event(gtx, &e, 0, 0, "→"):
			if ed, ok := m.value.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case m.value.Event(gtx, &e, 0, 0, "←"):
			m.key.Focus(gtx)
		case m.value.Event(gtx, &e, 0, 0, "⏎", "⌤", "⌫", "⌦"):
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
