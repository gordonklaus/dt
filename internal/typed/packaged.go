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
	ed.ed = NewTypeNameEditor(ed, pkg.Types[0], loader)
	return ed
}

func (ed *PackageEditor) Layout(gtx C) D {
	if gtx.Focused(nil) {
		ed.Focus(gtx)
	}

	for {
		e, ok := gtx.Event(key.Filter{Name: "S", Required: key.ModShortcut})
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.Event:
			if e.Name == "S" {
				ed.loader.Store(&types.PackageID_Current{})
			}
		}
	}

events:
	for {
		var e key.Event
		switch {
		default:
			break events
		case ed.FocusEvent(gtx):
		case ed.Event(gtx, &e, 0, 0, "→"):
			ed.ed.Focus(gtx)
		case ed.Event(gtx, &e, 0, key.ModShift, "↑"):
			if ed.focusedType > 0 && ed.Focused(gtx) {
				ed.focusedType--
				if e.Modifiers == key.ModShift {
					ed.pkg.Types[ed.focusedType], ed.pkg.Types[ed.focusedType+1] = ed.pkg.Types[ed.focusedType+1], ed.pkg.Types[ed.focusedType]
				} else {
					ed.ed = NewTypeNameEditor(ed, ed.pkg.Types[ed.focusedType], ed.loader)
				}
			}
		case ed.Event(gtx, &e, 0, key.ModShift, "↓"):
			if ed.focusedType < len(ed.pkg.Types)-1 && ed.Focused(gtx) {
				ed.focusedType++
				if e.Modifiers == key.ModShift {
					ed.pkg.Types[ed.focusedType], ed.pkg.Types[ed.focusedType-1] = ed.pkg.Types[ed.focusedType-1], ed.pkg.Types[ed.focusedType]
				} else {
					ed.ed = NewTypeNameEditor(ed, ed.pkg.Types[ed.focusedType], ed.loader)
				}
			}
		case ed.Event(gtx, &e, 0, key.ModShift, "⏎", "⌤"):
			var id uint64
		restart:
			for _, nt := range ed.pkg.Types {
				if nt.ID == id {
					id++
					goto restart
				}
			}
			n := &types.TypeName{ID: id}
			if e.Modifiers != key.ModShift && len(ed.pkg.Types) > 0 {
				ed.focusedType++
			}
			ed.pkg.Types = slices.Insert(ed.pkg.Types, ed.focusedType, n)
			ed.ed = NewTypeNameEditor(ed, n, ed.loader)
			ed.ed.Focus(gtx)
		case ed.Event(gtx, &e, 0, key.ModShift, "⌫", "⌦"):
			// TODO: Check if this type is referenced elsewhere and, if so, ask the user if they want to delete those references.
			if len(ed.pkg.Types) == 1 {
				ed.pkg.Types = []*types.TypeName{{}}
				ed.ed = NewTypeNameEditor(ed, ed.pkg.Types[0], ed.loader)
				ed.ed.Focus(gtx)
				break
			}
			ed.pkg.Types = slices.Delete(ed.pkg.Types, ed.focusedType, ed.focusedType+1)
			if e.Name == "⌫" && e.Modifiers == 0 && ed.focusedType > 0 || ed.focusedType == len(ed.pkg.Types) {
				ed.focusedType--
			}
			ed.ed = NewTypeNameEditor(ed, ed.pkg.Types[ed.focusedType], ed.loader)
		}
	}

	listRec := Record(gtx, func(gtx C) D {
		return ed.list.Layout(gtx, len(ed.pkg.Types), ed.layoutTypeName)
	})
	edRec := Record(gtx, ed.ed.Layout)

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
			if i != ed.focusedType || !ed.Focused(gtx) {
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
