package main

import (
	"fmt"

	"gioui.org/io/key"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

type FloatTypeEditor struct {
	typ *types.FloatType
}

func NewFloatTypeEditor(typ *types.FloatType) *FloatTypeEditor {
	return &FloatTypeEditor{
		typ: typ,
	}
}

func (b *FloatTypeEditor) Type() types.Type { return b.typ }

func (b *FloatTypeEditor) Layout(gtx C) D {
	for _, e := range gtx.Events(b) {
		if e, ok := e.(key.Event); ok && e.State == key.Press {
			switch e.Name {
			case "S":
				if b.typ.Size == 32 {
					b.typ.Size = 64
				} else {
					b.typ.Size = 32
				}
			}
		}
	}

	key.InputOp{
		Tag:  b,
		Keys: "S",
	}.Add(gtx.Ops)

	return material.Body1(theme, fmt.Sprintf("float%d", b.typ.Size)).Layout(gtx)
}
