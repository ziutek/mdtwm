package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

const rightButtonEventMask = xgb.EventMaskButtonPress |
	xgb.EventMaskButtonRelease | xgb.EventMaskPointerMotion

// Box for root window
type RootPanel struct {
	commonBox
}

func NewRootPanel() *RootPanel {
	var p RootPanel
	p.init(
		Window(screen.Root),
		xgb.EventMaskSubstructureRedirect|xgb.EventMaskStructureNotify,
	)
	p.width = int16(screen.WidthInPixels)
	p.height = int16(screen.HeightInPixels)
	p.w.ChangeAttrs(xgb.CWCursor, uint32(cfg.DefaultCursor))
	p.SetClass(cfg.Instance, cfg.Class)
	p.SetName("mdtwm root")
	p.w.ChangeProp(xgb.PropModeReplace, xgb.AtomCursor, xgb.AtomWindow, p.w)
	// Supported WM properities
	/*root.ChangeProp(xgb.PropModeReplace, AtomNetSupported,
	xgb.AtomAtom,	...)*/
	p.w.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, p.w)
	// Grab right mouse buttons for WM actions
	p.w.GrabButton(
		true, // Needed for EnterNotify events during grab
		rightButtonEventMask,
		xgb.GrabModeAsync, xgb.GrabModeAsync,
		xgb.WindowNone, cfg.DefaultCursor, 3,
		xgb.ButtonMaskAny,
	)
	// Grab keys for WM actions
	for k, _ := range cfg.Keys {
		p.w.GrabKey(true, cfg.ModMask, k, xgb.GrabModeAsync, xgb.GrabModeAsync)
	}
	return &p
}

func (p *RootPanel) Geometry() Geometry {
	return Geometry{
		X: p.x, Y: p.y,
		W: p.width, H: p.height,
	}
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
	b.SetPosSize(p.x, p.y, p.width, p.height)
	b.SetName("mdtwm desktop")
	b.Window().Map()
}

func (p *RootPanel) Remove(b Box, unmap bool) {
	if unmap {
		b.Window().Unmap()
	}
	p.children.Remove(b)
}
