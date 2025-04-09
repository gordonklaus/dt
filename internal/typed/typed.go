package typed

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type TypeEditor struct {
	*Core
	parent       Editor
	typ          *types.Type
	typeName     *types.TypeName
	mapKeyEditor bool
	creating     bool

	KeyFocus
	menu *Menu

	ed typeEditor
}

type Editor interface {
	CreateNext(gtx C, after *TypeEditor)
}

type typeEditor interface {
	Type() types.Type
	Layout(gtx C) D
}

func NewTypeNameTypeEditor(parent Editor, typ *types.TypeName, core *Core) *TypeEditor {
	t := NewTypeEditor(parent, &typ.Type, core)
	t.typeName = typ
	return t
}

func NewMapKeyTypeEditor(parent Editor, typ *types.Type, core *Core) *TypeEditor {
	t := NewTypeEditor(parent, typ, core)
	t.mapKeyEditor = true
	return t
}

func NewTypeEditor(parent Editor, typ *types.Type, core *Core) *TypeEditor {
	t := &TypeEditor{
		parent: parent,
		typ:    typ,
		Core:   core,
		menu:   NewMenu(core),
	}
	t.ed = t.newEditor(*typ)
	return t
}

func (t *TypeEditor) newEditor(typ types.Type) typeEditor {
	switch typ := typ.(type) {
	case nil:
		return nil
	case *types.BoolType:
		return NewBoolTypeEditor(typ)
	case *types.IntType:
		return NewIntTypeEditor(typ)
	case *types.FloatType:
		return NewFloatTypeEditor(typ)
	case *types.StringType:
		return NewStringTypeEditor(typ)
	case *types.StructType:
		return NewStructTypeEditor(t, typ, t.Core)
	case *types.EnumType:
		return NewEnumTypeEditor(t, typ, t.Core)
	case *types.ArrayType:
		return NewArrayTypeEditor(t, typ, t.Core)
	case *types.MapType:
		return NewMapTypeEditor(t, typ, t.Core)
	case *types.OptionType:
		return NewOptionTypeEditor(t, typ, t.Core)
	case *types.NamedType:
		return NewNamedTypeEditor(typ)
	}
	panic("unreached")
}

func (t *TypeEditor) Edit(gtx C) {
	t.menu.Items = nil
	if t.typeName != nil {
		t.menu.Add("struct", &types.StructType{})
		t.menu.Add("enum", &types.EnumType{})
	} else if t.mapKeyEditor {
		t.menu.Add("int", &types.IntType{})
		t.menu.Add("uint", &types.IntType{Unsigned: true})
		t.menu.Add("float32", &types.FloatType{Size: 32})
		t.menu.Add("float64", &types.FloatType{Size: 64})
		t.menu.Add("string", &types.StringType{})
	} else {
		t.menu.Add("bool", &types.BoolType{})
		t.menu.Add("int", &types.IntType{})
		t.menu.Add("uint", &types.IntType{Unsigned: true})
		t.menu.Add("float32", &types.FloatType{Size: 32})
		t.menu.Add("float64", &types.FloatType{Size: 64})
		t.menu.Add("string", &types.StringType{})
		t.menu.Add("array", &types.ArrayType{})
		t.menu.Add("map", &types.MapType{})
		t.menu.Add("option", &types.OptionType{})
		for _, n := range t.Pkg.Types {
			t.menu.Add(n.Name, &types.NamedType{Package: types.PackageID_Current{}, TypeName: n})
			if e, ok := n.Type.(*types.EnumType); ok {
				for _, el := range e.Elems {
					t.menu.Add(n.Name+"."+el.Name, &types.NamedType{Package: types.PackageID_Current{}, TypeName: el})
				}
			}
		}
	}
	t.menu.Focus(gtx)
}

func (t *TypeEditor) Layout(gtx C) D {
events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case t.FocusEvent(gtx):
		case t.Event(gtx, &e, 0, 0, "→"):
			if ed, ok := t.ed.(Focuser); ok {
				ed.Focus(gtx)
			}
		case t.Event(gtx, &e, 0, 0, "⏎", "⌤", "⌫", "⌦"):
			t.Edit(gtx)
		}
	}

	switch e := t.menu.Update(gtx).(type) {
	case MenuSubmitEvent:
		t.ed = t.newEditor(e.Item.Value.(types.Type))
		*t.typ = t.ed.Type()
		t.Focus(gtx)
		t.CreateNext(gtx, nil)
	case MenuCancelEvent:
		t.Focus(gtx)
	}

	// Both the menu and the editor must be laid out, or they will never take the focus.
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(t.menu.Layout),
		layout.Stacked(func(gtx C) D {
			if t.menu.Focused(gtx) {
				gtx.Constraints.Max = image.Point{}
			}
			if t.ed != nil {
				return t.KeyFocus.Layout(gtx, t.ed.Layout)
			}
			lbl := material.Body1(theme, "type")
			lbl.Font.Style = font.Italic
			return lbl.Layout(gtx)
		}),
	)
}

func (t *TypeEditor) CreateNext(gtx C, after typeEditor) {
	if after == nil {
		t.creating = true
	}
	if !t.creating {
		return
	}

	if ed, ok := t.ed.(Editor); ok && after == nil {
		ed.CreateNext(gtx, nil)
	} else {
		t.creating = false
		t.parent.CreateNext(gtx, t)
	}
}
