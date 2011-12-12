package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

// Box for WM window (panel)
type PanelBox struct {
	commonBox

	typ Orientation // panel type (vertical or horizontal)
}

// New PanelBox has parent set to nil and its window
// parent is root window. 
func NewPanelBox(typ Orientation) *PanelBox {
	var p PanelBox
	p.init(NewRawWindow(root, Geometry{0, 0, 1, 1, 0},
		xgb.WindowClassInputOutput, 0))
	p.typ = typ
	p.SetClass("mdtwm", "Mdtwm")
	p.SetName("mdtwm panel")
	return &p
}

func DeskPanelBox(typ Orientation) *PanelBox {
	var p PanelBox
	p.init(NewRawWindow(root, root.Geometry(), xgb.WindowClassInputOutput, 0))
	p.typ = typ
	p.SetClass("mdtwm", "Mdtwm")
	p.SetName("mdtwm desktop")
	return &p
}

func (p *PanelBox) SetPosSize(x, y, width, height int16) {
	p.SetGeometry(Geometry{x, y, width, height, 0})
}

func (p *PanelBox) SetFocus(f bool) {
	if f {
		currentPanel = p
		p.SetInputFocus()
	}
}

// Inserts a box into panel 
func (p *PanelBox) Insert(b Box) {
	b.SetParent(p)
	// TODO: Implement finding of best place to insert
	p.children.PushBack(b)
	// Rearange panel and show new box
	p.tile()
	b.Map()
}

func (p *PanelBox) Remove(b Box) {
	p.children.Remove(b)
	p.tile()
}

// Update geometry for boxes in panel
func (p *PanelBox) tile() {
	n := int16(p.children.Len())
	if n == 0 {
		return // there is no boxes in panel
	}
	pg := p.Geometry()
	if p.typ == Vertical {
		l.Print("tile V in: ", p)
		h := pg.H / n
		// Set new geometry for all boxes in panel
		y, w, h := int16(0), pg.W, h
		i := p.children.FrontIter(false)
		for ; n > 1; n-- {
			i.Next().SetPosSize(0, y, w, h)
			y += h
		}
		// Last window obtain all remaining space
		i.Next().SetPosSize(0, y, w, pg.H-y)
	} else {
		l.Print("tile H in:", p)
		w := pg.W / n
		// Set new geometry for all boxes in panel
		x, w, h := int16(0), w, pg.H
		i := p.children.FrontIter(false)
		for ; n > 1; n-- {
			i.Next().SetPosSize(x, 0, w, h)
			x += w
		}
		// Last window obtain all remaining space
		i.Next().SetPosSize(x, 0, pg.W-x, h)
	}
}
