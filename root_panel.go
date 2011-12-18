package main

import (
	"code.google.com/p/x-go-binding/xgb"
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
	p.w.SetEventMask(
		xgb.EventMaskSubstructureRedirect | // all config. req. redireted to WM
			xgb.EventMaskStructureNotify,
	)
	// Grab right mouse buttons for WM actions
	p.w.GrabButton(
		true, // Needed for EnterNotify events during grab
		xgb.EventMaskButtonPress|xgb.EventMaskButtonRelease,
			//|xgb.EventMaskPointerMotion*/,
		xgb.GrabModeAsync, xgb.GrabModeAsync,
		xgb.WindowNone, cfg.MoveCursor, 3,
		xgb.ButtonMaskAny,
	)
	// Grab keys for WM actions
	for k, _ := range cfg.Keys {
		p.w.GrabKey(true, cfg.ModMask, k, xgb.GrabModeAsync, xgb.GrabModeAsync)
	}
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
