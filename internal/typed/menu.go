package typed

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Menu struct {
	*Core
	Items, filteredItems []*MenuItem
	FocusedItem          int
	HideEditor           bool

	ed      widget.Editor
	oldText string
	list    layout.List
}

type MenuItem struct {
	Text   string
	Widget layout.Widget
	Value  any
}

type MenuEvent interface{ isMenuEvent() }

type MenuSubmitEvent struct{ Item *MenuItem }
type MenuCancelEvent struct{}

func (MenuSubmitEvent) isMenuEvent() {}
func (MenuCancelEvent) isMenuEvent() {}

func NewMenu(core *Core) *Menu {
	return &Menu{
		Core: core,
		ed: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		list: layout.List{Axis: layout.Vertical},
	}
}

func (m *Menu) Add(text string, value any) {
	m.AddWidget(text, material.Body1(theme, text).Layout, value)
}

func (m *Menu) AddWidget(text string, w layout.Widget, value any) {
	m.Items = append(m.Items, &MenuItem{
		Text:   text,
		Widget: w,
		Value:  value,
	})
}

func (m *Menu) Focus(gtx C) {
	m.ed.SetText("")
	m.oldText = ""
	m.applyFilter()
	m.list.Position = layout.Position{}
	gtx.Execute(key.FocusCmd{Tag: &m.ed})
}
func (m *Menu) Focused(gtx C) bool { return gtx.Focused(&m.ed) }

func (m *Menu) applyFilter() {
	text := m.ed.Text()
	var prefix, contains []*MenuItem
	for _, i := range m.Items {
		if strings.HasPrefix(i.Text, text) {
			prefix = append(prefix, i)
		} else if containsAll(i.Text, text) {
			contains = append(contains, i)
		}
	}
	filtered := append(prefix, contains...)
	if len(filtered) > 0 {
		m.filteredItems = filtered
		m.FocusedItem = 0
		m.oldText = text
	} else {
		_, pos := m.ed.CaretPos()
		m.ed.SetText(m.oldText)
		m.ed.SetCaret(pos, pos)
	}
}

func containsAll(s, t string) bool {
	for _, r := range t {
		var ok bool
		_, s, ok = strings.Cut(s, string(r))
		if !ok {
			return false
		}
	}
	return true
}

func (m *Menu) Update(gtx C) MenuEvent {
	for {
		e, ok := gtx.Event(
			key.Filter{Focus: &m.ed, Name: "↑"},
			key.Filter{Focus: &m.ed, Name: "↓"},
			key.Filter{Focus: &m.ed, Name: "⎋"},
		)
		if !ok {
			break
		}
		if e, ok := e.(key.Event); ok && e.State == key.Press {
			switch e.Name {
			case "↑":
				m.FocusedItem--
				m.updateScroll()
			case "↓":
				m.FocusedItem++
				m.updateScroll()
			case "⎋":
				return MenuCancelEvent{}
			}
		}
	}
	for {
		e, ok := m.ed.Update(gtx)
		if !ok {
			break
		}
		switch e.(type) {
		case widget.ChangeEvent:
			m.applyFilter()
		case widget.SubmitEvent:
			return MenuSubmitEvent{Item: m.filteredItems[m.FocusedItem]}
		}
	}
	return nil
}

const maxMenuItems = 12

func (m *Menu) updateScroll() {
	m.FocusedItem = (m.FocusedItem + len(m.filteredItems)) % len(m.filteredItems)
	pos := &m.list.Position
	first := pos.First
	last := first + (maxMenuItems - 1)
	if pos.Offset > 0 {
		first++
	}
	if m.FocusedItem < first {
		pos.First = m.FocusedItem
		pos.Offset = 0
	}
	if m.FocusedItem > last {
		pos.First = m.FocusedItem - (maxMenuItems - 1)
		pos.Offset = 0
	}
}

func (m *Menu) Layout(gtx C) D {
	m.Update(gtx)

	if !m.Focused(gtx) {
		// m.ed must be laid out or it will never take the focus (Editor.Layout is where op.Event is called).
		gtx.Constraints.Max = image.Point{}
		return material.Editor(theme, &m.ed, "").Layout(gtx)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(m.layoutEditor),
		layout.Rigid(layout.Spacer{Height: 2}.Layout),
		layout.Rigid(m.layoutList),
	)
}

func (m *Menu) layoutEditor(gtx C) D {
	if m.HideEditor {
		gtx.Constraints.Max = image.Point{}
	}
	return material.Editor(theme, &m.ed, "").Layout(gtx)
}

func (m *Menu) layoutList(gtx C) D {
	macro := op.Record(gtx.Ops)
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			r := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, gtx.Dp(4))
			paint.FillShape(gtx.Ops, theme.Bg, r.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, theme.Fg,
				clip.Stroke{
					Path:  r.Path(gtx.Ops),
					Width: 1,
				}.Op(),
			)
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
				width := 0
				for i := range m.filteredItems {
					width = max(width, Record(gtx.Disabled(), func(gtx C) D { return m.layoutItem(gtx, i) }).Dims.Size.X)
				}
				itemHeight := Record(gtx.Disabled(), func(gtx C) D { return m.layoutItem(gtx, 0) }).Dims.Size.Y
				gtx.Constraints.Min.X = width
				gtx.Constraints.Max.Y = maxMenuItems * itemHeight
				return m.list.Layout(gtx, len(m.filteredItems), m.layoutItem)
			})
		}),
	)
	op.Defer(gtx.Ops, macro.Stop())
	return D{}
}

func (m *Menu) layoutItem(gtx C, i int) D {
	return layout.UniformInset(4).Layout(gtx, func(gtx C) D {
		w := m.filteredItems[i].Widget
		if i == m.FocusedItem {
			minWidth := gtx.Constraints.Min.X
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					margin := gtx.Dp(4)
					defer op.Offset(image.Pt(-margin, -margin)).Push(gtx.Ops).Pop()
					paint.FillShape(gtx.Ops, color.NRGBA{A: 64},
						clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min.Add(image.Pt(2*margin, 2*margin))}, 2*margin).Op(gtx.Ops))
					return D{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min.X = minWidth
					return w(gtx)
				}),
			)
		}
		return w(gtx)
	})
}
