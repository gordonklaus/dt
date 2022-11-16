package main

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type IntTypeEditor struct {
	typ  types.Type
	size widget.Editor
}

func NewIntTypeEditor(typ types.Type) *IntTypeEditor {
	return &IntTypeEditor{
		typ: typ,
		size: widget.Editor{
			SingleLine: true,
			// MaxLen:     2,
			// Filter:     "0123456789",
		},
	}
}

func (i *IntTypeEditor) Type() types.Type { return i.typ }

func (i *IntTypeEditor) Layout(gtx C) D {
	var size uint64
	switch t := i.typ.(type) {
	case *types.UintType:
		size = t.Size
	case *types.IntType:
		size = t.Size
	}

	for _, e := range i.size.Events() {
		if _, ok := e.(widget.ChangeEvent); ok {
			if x, err := strconv.ParseUint(i.size.Text(), 10, 64); err == nil {
				size = x
				if size == 0 || size > 64 {
					size = 64
				}
				switch t := i.typ.(type) {
				case *types.UintType:
					t.Size = size
				case *types.IntType:
					t.Size = size
				}
			}
		}
	}

	_, caret := i.size.CaretPos()
	i.size.SetText(strconv.FormatUint(size, 10))
	i.size.SetCaret(caret, caret)

	s := "uint"
	if _, ok := i.typ.(*types.IntType); ok {
		s = "int"
	}
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(material.Body1(theme, s).Layout),
		layout.Rigid(material.Editor(theme, &i.size, "").Layout),
	)
}
