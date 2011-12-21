package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

// Box for WM window (panel)
type Panel struct {
	commonBox

	typ   Orientation // panel type (vertical or horizontal)
	ratio float64     // ratio of size of two neighboringsubwindows
}

// New Panel has parent set to nil and its window
// parent is root window.
// ratio == 1 means all subwindows in panel are equal in size.
func NewPanel(typ Orientation, ratio float64) *Panel {
	var p Panel
	p.init(
		NewWindow(root.Window(), Geometry{0, 0, 1, 1, 0},
			xgb.WindowClassInputOutput, 0),
		xgb.EventMaskEnterWindow|xgb.EventMaskStructureNotify|
			xgb.EventMaskSubstructureRedirect,
	)
	p.width = 1
	p.height = 1
	p.typ = typ
	p.ratio = ratio
	p.SetClass(cfg.Instance, cfg.Class)
	p.SetName("mdtwm panel")
	p.w.SetBackPixmap(xgb.BackPixmapParentRelative)
	return &p
}

func (p *Panel) Geometry() Geometry {
	return Geometry{
		X: p.x, Y: p.y,
		W: p.width, H: p.height,
	}
}

func (p *Panel) ReqPosSize(x, y, width, height int16) {
	p.x, p.y, p.width, p.height = x, y, width, height
	p.w.SetGeometry(Geometry{x, y, width, height, 0})
}

func (p *Panel) SyncGeometry(g Geometry) {
	if g.B != 0 {
		l.Print("non-zero border width: ", g.B)
	}
	p.x, p.y, p.width, p.height = g.X, g.Y, g.W, g.H
	p.tile()
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
	if !b.Float() {
		// Rearange panel and show new box
		p.tile()
	}
	b.Window().Map()
	if w, ok := b.(*BoxedWindow); ok {
		w.SetWmState(WmStateNormal)
	}
}

func (p *Panel) Remove(b Box) {
	b.SetParent(root)
	p.children.Remove(b)
	p.tile()
}

// Update geometry for boxes in panel
func (p *Panel) tile() {
	i := p.children.FrontIter(false)
	var n int16
	for b := i.Next(); b != nil; b = i.Next() {
		if !b.Float() {
			n++
		}
	}
	if n == 0 {
		return // there is no boxes in panel
	}
	if p.typ == Vertical {
		d.Print("Tile V in: ", p)
		hg := NewSizeGen(p.height, n, p.ratio)
		i = p.children.FrontIter(false)
		y, w := int16(0), p.width
		for n > 1 {
			b := i.Next()
			if b.Float() {
				continue
			}
			h := hg.Next()
			b.ReqPosSize(0, y, w, h)
			y += h
			n--
		}
		// Last window obtain all remaining space
		i.Next().ReqPosSize(0, y, w, p.height-y)
	} else {
		d.Print("Tile H in:", p)
		wg := NewSizeGen(p.width, n, p.ratio)
		x, h := int16(0), p.height
		i = p.children.FrontIter(false)
		for n > 1 {
			b := i.Next()
			if b.Float() {
				continue
			}
			w := wg.Next()
			b.ReqPosSize(x, 0, w, h)
			x += w
			n--
		}
		// Last window obtain all remaining space
		i.Next().ReqPosSize(x, 0, p.width-x, h)
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
