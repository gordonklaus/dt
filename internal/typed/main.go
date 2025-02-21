package typed

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type C = layout.Context
type D = layout.Dimensions

var theme = material.NewTheme()

func init() {
	theme.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
}

func Edit(loader *types.Loader, pkg *types.Package) {
	go edit(loader, pkg)
	app.Main()
}

func edit(loader *types.Loader, pkg *types.Package) {
	var w app.Window
	w.Option(app.Title("typEd"))
	w.Perform(system.ActionMaximize)

	ed := NewPackageEditor(&Core{
		Loader: loader,
		Pkg:    pkg,
	})

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, e)

			// Disable tab navigation globally.
			for ok := true; ok; _, ok = gtx.Event(key.Filter{Name: key.NameTab, Optional: key.ModShift}) {
			}

			layout.Center.Layout(gtx, ed.Layout)
			e.Frame(&ops)
		case app.DestroyEvent:
			if e.Err != nil {
				fmt.Println(e.Err)
			}
			os.Exit(0)
		}
	}
}
