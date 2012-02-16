package mdtwm

import (
	"github.com/ziutek/mdtwm/xgb_patched"
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

func (p *Panel) SetPosSize(x, y, width, height int16) {
	p.x, p.y, p.width, p.height = x, y, width, height
	p.w.SetGeometry(Geometry{x, y, width, height, 0})
	p.tile()
}

func (p *Panel) SetFocus(f bool, t xgb.Timestamp) {
	if f {
		p.w.SetInputFocus(t)
	}
}

// Inserts b next to mark. Can use x, y (coordinates in mark) for decision.
func (p *Panel) InsertNextTo(b, mark Box, x, y int16) {
	var ok bool
	mg := mark.Geometry()
	if p.typ == Vertical && y < mg.H/2 || p.typ == Horizontal && x < mg.W/2 {
		ok = p.children.InsertBefore(b, mark)
	} else {
		ok = p.children.InsertAfter(b, mark)
	}
	if !ok {
		p.children.PushBack(b)
	}
	p.insertCommon(b)
}

func (p *Panel) InsertBefore(b, mark Box) {
	if !p.children.InsertBefore(b, mark) {
		p.children.PushBack(b)
	}
	p.insertCommon(b)
}

func (p *Panel) Append(b Box) {
	p.children.PushBack(b)
	p.insertCommon(b)
}

// Inserts a box into panel 
func (p *Panel) insertCommon(b Box) {
	b.SetParent(p)
	if !b.Float() {
		// Rearange panel and show new box
		p.tile()
	}
	b.Window().Map()
	if w, ok := b.(*BoxedWindow); ok {
		w.SetWmState(WmStateNormal)
		w.UpdateNetWmDesktop()
	}
}

func (p *Panel) Remove(b Box) {
	b.SetParent(root)
	p.children.Remove(b)
	p.tile()
}

// Update geometry for boxes in panel
func (p *Panel) tile() {
	var (
		n int16
		b Box
	)
	for b = p.children.Front(); b != nil; b = b.Next() {
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
		y, w := int16(0), p.width
		for b = p.children.Front(); n > 1; b, n = b.Next(), n - 1 {
			if b.Float() {
				continue
			}
			h := hg.Next()
			b.SetPosSize(0, y, w, h)
			y += h
		}
		// Last window obtain all remaining space
		b.SetPosSize(0, y, w, p.height-y)
	} else {
		d.Print("Tile H in:", p)
		wg := NewSizeGen(p.width, n, p.ratio)
		x, h := int16(0), p.height
		for b = p.children.Front(); n > 1; b, n = b.Next(), n - 1 {
			if b.Float() {
				continue
			}
			w := wg.Next()
			b.SetPosSize(x, 0, w, h)
			x += w
		}
		// Last window obtain all remaining space
		b.SetPosSize(x, 0, p.width-x, h)
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
