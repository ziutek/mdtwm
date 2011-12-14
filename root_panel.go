package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

// Box for root window
type RootPanel struct {
	commonBox
}

func NewRootPanel() *RootPanel {
	var p RootPanel
	p.init(Window(screen.Root))
	p.SetClass(cfg.Instance, cfg.Class)
	p.SetName("mdtwm root")
	// Supported WM properities
	/*root.ChangeProp(xgb.PropModeReplace, AtomNetSupported,
	xgb.AtomAtom,	...)*/
	p.w.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, p.w)
	// Event mask for WM root
	p.w.SetEventMask(xgb.EventMaskSubstructureRedirect |
		xgb.EventMaskStructureNotify |
		//xgb.EventMaskPointerMotion |
		xgb.EventMaskPropertyChange |
		xgb.EventMaskEnterWindow)
	p.grabInput(p.w)
	return &p
}

func (p *RootPanel) SetPosSize(x, y, width, height int16) {
	panic("Can't change position of root window")
}

func (p *RootPanel) SetFocus(f bool) {
	return
}

// Inserts a box into panel 
func (p *RootPanel) Insert(b Box) {
	b.SetParent(p)
	p.children.PushBack(b)
	g := p.w.Geometry()
	b.SetPosSize(g.X, g.Y, g.W, g.H)
	b.SetName("mdtwm desktop")
	b.Window().Map()
}

func (p *RootPanel) Remove(b Box) {
	p.children.Remove(b)
}
