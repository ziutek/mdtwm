package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

// Box for WM window (panel)
type Panel struct {
	commonBox

	typ   Orientation // panel type (vertical or horizontal)
	ratio float64     // ratio of width/height betwen two neighboring subwindows
}

// New Panel has parent set to nil and its window
// parent is root window.
// ratio == 1 means all subwindows in panel are equal in size.
func NewPanel(typ Orientation, ratio float64) *Panel {
	var p Panel
	p.init(NewWindow(root.Window(), Geometry{0, 0, 1, 1, 0},
		xgb.WindowClassInputOutput, 0))
	p.typ = typ
	p.ratio = ratio
	p.SetClass(cfg.Instance, cfg.Class)
	p.SetName("mdtwm panel")
	p.w.SetBackPixmap(xgb.BackPixmapParentRelative)
	p.w.SetEventMask(boxEventMask)
	p.grabInput(root.Window())
	return &p
}

func (p *Panel) SetPosSize(x, y, width, height int16) {
	p.w.SetGeometry(Geometry{x, y, width, height, 0})
}

func (p *Panel) SetFocus(f bool) {
	if f {
		currentBox = p
		p.w.SetInputFocus()
	}
}

// Inserts a box into panel 
func (p *Panel) Insert(b Box) {
	b.SetParent(p)
	// TODO: Implement finding of best place to insert
	p.children.PushBack(b)
	// Rearange panel and show new box
	p.tile()
	b.Window().Map()
}

func (p *Panel) Remove(b Box) {
	p.children.Remove(b)
	p.tile()
}

// Update geometry for boxes in panel
func (p *Panel) tile() {
	n := int16(p.children.Len())
	if n == 0 {
		return // there is no boxes in panel
	}
	pg := p.w.Geometry()
	if p.typ == Vertical {
		l.Print("tile V in: ", p)
		hg := NewSizeGen(pg.H, n, p.ratio)
		i := p.children.FrontIter(false)
		y, w := int16(0), pg.W
		for ; n > 1; n-- {
			h := hg.Next()
			i.Next().SetPosSize(0, y, w, h)
			y += h
		}
		// Last window obtain all remaining space
		i.Next().SetPosSize(0, y, w, pg.H-y)
	} else {
		l.Print("tile H in:", p)
		wg := NewSizeGen(pg.W, n, p.ratio)
		x, h := int16(0), pg.H
		i := p.children.FrontIter(false)
		for ; n > 1; n-- {
			w := wg.Next()
			i.Next().SetPosSize(x, 0, w, h)
			x += w
		}
		// Last window obtain all remaining space
		i.Next().SetPosSize(x, 0, pg.W-x, h)
	}
}

type SizeGen struct {
	s, ratio float64
}

func NewSizeGen(allSpace, n int16, ratio float64) *SizeGen {
	d := 1.0
	for ; n > 1; n-- {
		d = d*ratio + 1
	}
	return &SizeGen{float64(allSpace) / d, ratio}
}

func (g *SizeGen) Next() (s int16) {
	s = int16(g.s)
	g.s *= g.ratio
	return
}
