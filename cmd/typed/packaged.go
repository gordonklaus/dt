package main

import (
	"gioui.org/io/key"
	"github.com/gordonklaus/data/types"
)

type PackageEditor struct {
	pkg    *types.Package
	loader *types.Loader
	ed     *TypeNameEditor
}

func NewPackageEditor(pkg *types.Package, typ *types.TypeName, loader *types.Loader) *PackageEditor {
	ed := &PackageEditor{
		pkg:    pkg,
		loader: loader,
		ed:     NewTypeNameEditor(typ, loader),
	}
	return ed
}

func (ed *PackageEditor) Layout(gtx C) D {
	for _, e := range gtx.Events(ed) {
		switch e := e.(type) {
		case key.Event:
			if e.State == key.Press {
				ed.loader.Store(&types.PackageID_Current{})
			}
		}
	}

	key.InputOp{
		Tag:  ed,
		Keys: "Short-S",
	}.Add(gtx.Ops)

	return ed.ed.Layout(gtx)
}
