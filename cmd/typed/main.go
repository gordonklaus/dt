package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
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
	typeName := flag.String("type", "", "")
	flag.Parse()
	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	if *typeName == "" {
		log.Fatal("type argument is required (for now)")
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	loader := types.NewLoader(types.NewStorage(dir))
	pkg, err := loader.Load(&types.PackageID_Current{}) // TODO: Resolve current package ID based on current directory and source control/module configuration.
	if errors.Is(err, fs.ErrNotExist) {
		pkg = &types.Package{Name: filepath.Base(dir)}
		loader.Packages[&types.PackageID_Current{}] = pkg
	} else if err != nil {
		log.Fatal(err)
	}

	typ := pkg.Type(*typeName)
	if typ == nil {
		typ = &types.TypeName{
			Name: *typeName,
		}
		pkg.Types = append(pkg.Types, typ)
	}

	w := app.NewWindow(app.Title("typEd"))
	w.Perform(system.ActionMaximize)

	ed := NewPackageEditor(pkg, typ, loader)

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			ops.Reset()
			gtx := layout.NewContext(&ops, e)

			key.InputOp{Tag: w, Keys: "Tab"}.Add(gtx.Ops) // Disable tab navigation globally.

			layout.Center.Layout(gtx, ed.Layout)
			e.Frame(&ops)
		case system.DestroyEvent:
			if e.Err != nil {
				log.Print(e.Err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}
