package typed

import (
	"image"
	"image/color"
	"log"
	"slices"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/gordonklaus/dt/types"
	"golang.org/x/exp/maps"
)

type PackageEditor struct {
	*Core

	KeyFocus
	list        layout.List
	focusedType int
	ed          *TypeNameEditor
}

func NewPackageEditor(core *Core) *PackageEditor {
	ed := &PackageEditor{
		Core: core,
		list: layout.List{
			Axis: layout.Vertical,
		},
	}
	if len(ed.Pkg.Types) == 0 {
		n := &types.TypeName{Parent: ed.Pkg}
		ed.Pkg.Types = []*types.TypeName{n}
		ed.Pkg.TypesByID[n.ID] = n
	}
	ed.ed = NewTypeNameEditor(ed, ed.Pkg.Types[0], core)
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
				if err := ed.Loader.Store(types.PackageID_Current{}); err != nil {
					log.Println(err)
				}
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
			if ed.focusedType > 0 {
				ed.focusedType--
				if e.Modifiers == key.ModShift {
					ed.Pkg.Types[ed.focusedType], ed.Pkg.Types[ed.focusedType+1] = ed.Pkg.Types[ed.focusedType+1], ed.Pkg.Types[ed.focusedType]
				} else {
					ed.ed = NewTypeNameEditor(ed, ed.Pkg.Types[ed.focusedType], ed.Core)
				}
			}
		case ed.Event(gtx, &e, 0, key.ModShift, "↓"):
			if ed.focusedType < len(ed.Pkg.Types)-1 {
				ed.focusedType++
				if e.Modifiers == key.ModShift {
					ed.Pkg.Types[ed.focusedType], ed.Pkg.Types[ed.focusedType-1] = ed.Pkg.Types[ed.focusedType-1], ed.Pkg.Types[ed.focusedType]
				} else {
					ed.ed = NewTypeNameEditor(ed, ed.Pkg.Types[ed.focusedType], ed.Core)
				}
			}
		case ed.Event(gtx, &e, 0, key.ModShift, "⏎", "⌤"):
			n := &types.TypeName{ID: nextID(ed.Pkg), Parent: ed.Pkg}
			if e.Modifiers != key.ModShift && len(ed.Pkg.Types) > 0 {
				ed.focusedType++
			}
			ed.Pkg.Types = slices.Insert(ed.Pkg.Types, ed.focusedType, n)
			ed.Pkg.TypesByID[n.ID] = n
			ed.ed = NewTypeNameEditor(ed, n, ed.Core)
			ed.ed.named.Edit(gtx)
		case ed.Event(gtx, &e, 0, key.ModShift, "⌫", "⌦"):
			// TODO: Check if this type is referenced elsewhere and, if so, ask the user if they want to delete those references.
			switch ed := ed.ed.typed.ed.(type) {
			case *StructTypeEditor:
				for len(ed.fields) > 0 {
					ed.deleteField(gtx, ed.fields[0], true)
				}
			case *EnumTypeEditor:
				for len(ed.elems) > 0 {
					ed.deleteElem(gtx, ed.elems[0], true)
				}
			}
			delete(ed.Pkg.TypesByID, ed.Pkg.Types[ed.focusedType].ID)
			if len(ed.Pkg.Types) == 1 {
				n := &types.TypeName{Parent: ed.Pkg}
				ed.Pkg.Types = []*types.TypeName{n}
				ed.Pkg.TypesByID[n.ID] = n
				ed.ed = NewTypeNameEditor(ed, n, ed.Core)
				ed.ed.Focus(gtx)
				break
			}
			ed.Pkg.Types = slices.Delete(ed.Pkg.Types, ed.focusedType, ed.focusedType+1)
			if e.Name == "⌫" && e.Modifiers == 0 && ed.focusedType > 0 || ed.focusedType == len(ed.Pkg.Types) {
				ed.focusedType--
			}
			ed.ed = NewTypeNameEditor(ed, ed.Pkg.Types[ed.focusedType], ed.Core)
		}
	}

	listRec := Record(gtx, func(gtx C) D {
		return ed.list.Layout(gtx, len(ed.Pkg.Types), ed.layoutTypeName)
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
				return material.Body1(theme, ed.Pkg.Types[i].Name).Layout(gtx)
			})
		}),
	)
}

func nextID(p *types.Package) uint64 {
	ids := maps.Keys(p.TypesByID)
	slices.Sort(ids)
	var id uint64
	for i := range ids {
		if id != ids[i] {
			break
		}
		id++
	}
	return id
}
