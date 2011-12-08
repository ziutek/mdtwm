package main

import (
	"math"
	"x-go-binding.googlecode.com/hg/xgb"
)

const (
	WindowEventMask = xgb.EventMaskPropertyChange |
		xgb.EventMaskButtonRelease |
		//xgb.EventMaskPointerMotion |
		xgb.EventMaskExposure | // window needs to be redrawn
		xgb.EventMaskStructureNotify | // window gets destroyed
		xgb.EventMaskSubstructureRedirect | // app tries to resize itself
		xgb.EventMaskSubstructureNotify | // subwindows get notifies
		xgb.EventMaskEnterWindow |
		xgb.EventMaskFocusChange
)

var x int16

func winAdd(w, parent Window) {
	l.Print("manageWindow: ", w)
	if cfg.Ignore.Contains(w.Class()) {
		return
	}
	b := NewBox()
	b.Window = w

	// Don't map if unvisible or has OverrideRedirect flag
	attr := w.Attrs()
	if attr.MapState != xgb.MapStateViewable || attr.OverrideRedirect {
		return
	}
	// Check window type
	p, err := w.Prop(AtomNetWmWindowType, math.MaxUint32)
	if err != nil {
		wm_type := propReplyAtoms(p)
		if wm_type.Contains(AtomNetWmWindowTypeDock) {
			l.Printf("Window %s is of type dock", w)
		}
		if cfg.Float.Contains(w.Class()) ||
			wm_type.Contains(AtomNetWmWindowTypeDialog) ||
			wm_type.Contains(AtomNetWmWindowTypeUtility) ||
			wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
			wm_type.Contains(AtomNetWmWindowTypeSplash) {
			b.Float = true
		}
	} else {
		l.Printf("Can't get AtomNetWmWindowType from %s: %s", w, err)
	}
	// Grab left and right mouse buttons for click to focus/rasie
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 1, xgb.ButtonMaskAny)
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 3, xgb.ButtonMaskAny)

	w.SetEventMask(WindowEventMask)
	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)
	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	// Find box in which we have to put this window
	pbox := currentPanel
	if parent != root {
		pbox = pbox.Children.BoxByWindow(parent)
	}
	tile(b, pbox)
}

func tile(b, parent *Box) {
	parent.Children.PushFront(b)
	b.Window.Reparent(parent.Window, 100, 100)
	//g := b.Geometry()
}

func winFocus(w Window) {
	l.Print("Focusing window: ", w)
	for bi := currentDesk.Children.FrontIter(true); !bi.Done(); {
		b := bi.Next()
		if b.Window == w {
			b.Window.SetBorderColor(cfg.FocusedBorderColor)
			w.SetInputFocus()
		} else {
			b.Window.SetBorderColor(cfg.NormalBorderColor)
		}
	}
}

// Panel is a box with InputOnly (transparent) window in which a real windows
// are placed. A desk is organized as collection of panels in some layout
func newPanel(g Geometry) *Box {
	p := NewBox()
	p.Window = NewWindow(root, g, xgb.WindowClassInputOnly,
		xgb.CWOverrideRedirect|xgb.CWEventMask, 1, WindowEventMask)
	return p
}
