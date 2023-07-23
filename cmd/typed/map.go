package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
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
	m.value = NewTypeEditor(m, &typ.Value, loader)
	if typ.Key == nil {
		m.key.showMenu = true
	} else if typ.Value == nil {
		m.value.showMenu = true
	}
	return m
}

func (m *MapTypeEditor) Type() types.Type { return m.typ }

func (m *MapTypeEditor) Layout(gtx C) D {
	for _, e := range m.KeyFocus.Events(gtx, "←|→|↑|↓|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			m.focusValue.Focus()
		case "←":
			switch p := m.parent.parent.(type) {
			case *StructFieldTypeEditor:
				p.focusTyped.Focus()
			case interface{ Focus() }:
				p.Focus()
			}
			// case "↑":
			// 	m.parent.focusNext(false)
			// case "↓":
			// 	m.parent.focusNext(true)
		case "⏎", "⌤", "⌫", "⌦":
			m.key.showMenu = true
		}
	}

	for _, e := range m.focusValue.Events(gtx, "←|→|↑|↓|⏎|⌤|⌫|⌦") {
		switch e.Name {
		case "→":
			if ed, ok := m.value.ed.(interface{ Focus() }); ok {
				ed.Focus()
			}
		case "←":
			m.Focus()
			// case "↑":
			// 	a.parent.focusNext(false)
			// case "↓":
			// 	a.parent.focusNext(true)
		case "⏎", "⌤", "⌫", "⌦":
			m.value.showMenu = true
		}
	}

	if m.Focused() {
		if m.typ.Key == nil {
			*m.parent.typ = nil
			m.parent.ed = nil
			m.parent.showMenu = true
		} else if m.typ.Value == nil {
			m.value.showMenu = true
		}
	}

	if m.focusValue.Focused() && m.typ.Value == nil {
		*m.key.typ = nil
		m.key.ed = nil
		m.key.showMenu = true
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
