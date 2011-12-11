package main

import (
	"math"
	"unicode/utf16"
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
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
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
	if err == nil {
		wm_type := propReplyAtoms(p)
		if wm_type.Contains(AtomNetWmWindowTypeDock) {
			l.Printf("Window %s is of type dock", w)
		}
		if cfg.Float.Contains(class) ||
			wm_type.Contains(AtomNetWmWindowTypeDialog) ||
			wm_type.Contains(AtomNetWmWindowTypeUtility) ||
			wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
			wm_type.Contains(AtomNetWmWindowTypeSplash) {
			b.Float = true
		}
	} else {
		l.Printf("Can't get AtomNetWmWindowType from %s: %s", w, err)
	}
	// Update informations
	b.Name = w.Name()
	b.NameX = utf16.Encode([]rune(b.Name))

	// Grab left and right mouse buttons for click to focus/rasie
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 1, xgb.ButtonMaskAny)
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 3, xgb.ButtonMaskAny)

	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)
	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	// Find box in which we have to put this window
	parentBox := currentPanel
	if parent != root {
		parentBox = parentBox.Children.BoxByWindow(parent)
	}
	// Add window to found parentBox
	w.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmapNotify due to reparenting
	w.Reparent(parentBox.Window, 0, 0)
	w.SetEventMask(WindowEventMask) // set desired event mask
	parentBox.Children.PushBack(b)
	// Update geometry of windows in parentBox
	tile(parentBox)
}

func winFocus(w Window) {
	l.Print("Focusing window: ", w)
	// Iterate over panels
	pi := currentDesk.Children.FrontIter(false)
	for p := pi.Next(); p != nil; p = pi.Next() {
		// Iterate over full tree of windows in panel
		bi := p.Children.FrontIter(true);
		for b := bi.Next(); b != nil; b = bi.Next() {
			if b.Window == w {
				b.Window.SetBorderColor(cfg.FocusedBorderColor)
				w.SetInputFocus()
				currentPanel = p // Change current panel
			} else {
				b.Window.SetBorderColor(cfg.NormalBorderColor)
			}
		}
	}
}

// Panel is a box with InputOnly (transparent) window in which a real windows
// are placed. A desk is organized as collection of panels in some layout
func newPanel(g Geometry) *Box {
	p := NewBox()
	p.Window = NewWindow(root, g, xgb.WindowClassInputOutput,
		xgb.CWOverrideRedirect|xgb.CWEventMask, 1, WindowEventMask)
	p.Window.SetName("mdtwm panel")
	return p
}
