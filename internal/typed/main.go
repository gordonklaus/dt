package typed

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type C = layout.Context
type D = layout.Dimensions

var theme = material.NewTheme(gofont.Collection())

func Run(loader *types.Loader, pkg *types.Package) {
	go run(loader, pkg)
	app.Main()
}

func run(loader *types.Loader, pkg *types.Package) {
	w := app.NewWindow(app.Title("typEd"))
	w.Perform(system.ActionMaximize)

	ed := NewPackageEditor(pkg, loader)

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			ops.Reset()
			gtx := layout.NewContext(&ops, e)

			key.InputOp{Tag: w, Keys: "(Shift)-Tab"}.Add(gtx.Ops) // Disable tab navigation globally.

			layout.Center.Layout(gtx, ed.Layout)
			e.Frame(&ops)
		case system.DestroyEvent:
			if e.Err != nil {
				fmt.Println(e.Err)
			}
		}
	}
}
