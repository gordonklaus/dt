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

func NewMapTypeEditor(parent *TypeEditor, typ *types.MapType, core *Core) *MapTypeEditor {
	m := &MapTypeEditor{
		parent: parent,
		typ:    typ,
	}
	m.key = NewMapKeyTypeEditor(m, &typ.Key, core)
	m.value = NewTypeEditor(m, &typ.Value, core)
	return m
}

func (m *MapTypeEditor) Type() types.Type { return m.typ }

func (m *MapTypeEditor) CreateNext(gtx C, after *TypeEditor) {
	if !m.parent.creating {
		return
	}

	if after == nil {
		m.key.Edit(gtx)
	} else if after == m.key {
		m.value.Edit(gtx)
	} else {
		m.parent.CreateNext(gtx, m)
	}
}

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
		}
	}

vevents:
	for {
		var e key.Event
		switch {
		default:
			break vevents
		case m.value.FocusEvent(gtx):
		case m.value.Event(gtx, &e, 0, 0, "←"):
			m.key.Focus(gtx)
		}
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(material.Body1(theme, "map[").Layout),
		layout.Rigid(m.key.Layout),
		layout.Rigid(material.Body1(theme, "]").Layout),
		layout.Rigid(m.value.Layout),
	)
}
