package typed

import (
	"image"
	"image/color"
	"slices"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
)

type PackageEditor struct {
	pkg    *types.Package
	loader *types.Loader

	KeyFocus
	list        layout.List
	focusedType int
	ed          *TypeNameEditor
}

func NewPackageEditor(pkg *types.Package, loader *types.Loader) *PackageEditor {
	ed := &PackageEditor{
		pkg:    pkg,
		loader: loader,
		list: layout.List{
			Axis: layout.Vertical,
		},
	}
	if len(pkg.Types) == 0 {
		pkg.Types = []*types.TypeName{{}}
	}
	ed.ed = NewTypeNameEditor(pkg.Types[0], loader)
	ed.ed.Focus()
	return ed
}

func (ed *PackageEditor) Layout(gtx C) D {
	for _, e := range ed.Events(gtx, "←|→|↑|↓|(Shift)-[⏎,⌤]|Short-S") {
		switch e.Name {
		case "←":
			ed.Focus()
		case "→":
			ed.ed.Focus()
		case "↑":
			if ed.focusedType > 0 {
				ed.focusedType--
				ed.ed = NewTypeNameEditor(ed.pkg.Types[ed.focusedType], ed.loader)
			}
		case "↓":
			if ed.focusedType < len(ed.pkg.Types)-1 {
				ed.focusedType++
				ed.ed = NewTypeNameEditor(ed.pkg.Types[ed.focusedType], ed.loader)
			}
		case "⏎", "⌤":
			n := &types.TypeName{}
			if e.Modifiers != key.ModShift && len(ed.pkg.Types) > 0 {
				ed.focusedType++
			}
			ed.pkg.Types = slices.Insert(ed.pkg.Types, ed.focusedType, n)
			ed.ed = NewTypeNameEditor(n, ed.loader)
			ed.ed.Focus()
		case "S":
			ed.loader.Store(&types.PackageID_Current{})
		}
	}

	listRec := Record(gtx, func(gtx C) D {
		return ed.list.Layout(gtx, len(ed.pkg.Types), ed.layoutTypeName)
	})
	var edRec Recording
	if ed.ed != nil {
		edRec = Record(gtx, ed.ed.Layout)
	}

	w2 := gtx.Metric.PxToDp(gtx.Constraints.Max.X / 2)
	l2 := gtx.Metric.PxToDp(listRec.Dims.Size.X / 2)
	e2 := gtx.Metric.PxToDp(edRec.Dims.Size.X / 2)
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(layout.Spacer{Width: w2 - 256 - l2}.Layout),
		layout.Rigid(listRec.Layout),
		layout.Rigid(layout.Spacer{Width: 256 - l2 - e2}.Layout),
		layout.Rigid(edRec.Layout),
		layout.Rigid(layout.Spacer{Width: w2 - e2}.Layout),
	)
}

func (ed *PackageEditor) layoutTypeName(gtx C, i int) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			if i != ed.focusedType || !ed.Focused() {
				return D{}
			}
			paint.FillShape(gtx.Ops, color.NRGBA{A: 64},
				clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(4)).Op(gtx.Ops))
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
				return material.Body1(theme, ed.pkg.Types[i].Name).Layout(gtx)
			})
		}),
	)
}
