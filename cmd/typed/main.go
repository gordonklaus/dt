package main

import (
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/gordonklaus/data/types"
)

func main() {
	go Main()
	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

var theme = material.NewTheme(gofont.Collection())

func Main() {
	w := app.NewWindow(app.Title("typEd"))

	ed := NewNamedTypeEditor(&types.NamedType{
		Name: "person",
		Type: &types.StructType{
			Fields: []*types.StructFieldType{
				{Name: "id", Type: &types.UintType{Size: 64}},
				{Name: "name", Type: &types.StringType{}},
				{Name: "age", Type: &types.IntType{Size: 64}},
			},
		},
	})

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			ops.Reset()
			gtx := layout.NewContext(&ops, e)
			layout.Center.Layout(gtx, ed.Layout)
			e.Frame(&ops)
		case system.DestroyEvent:
			if e.Err != nil {
				log.Print(e.Err)
			}
			return
		}
	}
}
