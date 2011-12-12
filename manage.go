package main

import (
	"math"
	"unicode/utf16"
	"x-go-binding.googlecode.com/hg/xgb"
)

const (
	/* EventMask = xgb.EventMaskPropertyChange |
		xgb.EventMaskButtonRelease |
		xgb.EventMaskPointerMotion |
		xgb.EventMaskExposure | // window needs to be redrawn
		xgb.EventMaskStructureNotify | // window gets destroyed
		xgb.EventMaskSubstructureRedirect | // app tries to resize itself
		xgb.EventMaskSubstructureNotify | // subwindows get notifies
		xgb.EventMaskFocusChange |
		xgb.EventMaskEnterWindow */
	EventMask = xgb.EventMaskEnterWindow
)

var x int16

func manageWindow(w Window, panel *Box) {
	l.Print("manageWindow: ", w)
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if allDesks.BoxByWindow(w, true) != nil {
		l.Printf("  %s - alredy managed", w)
		return
	}
	attr := w.Attrs()
	// Don't manage and don't map if unvisible or has OverrideRedirect flag
	if attr.MapState != xgb.MapStateViewable || attr.OverrideRedirect {
		return
	}

	b := NewBox(BoxTypeWindow, w)

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
			l.Printf("Window %s should be treated as float", w)
		}
	} else {
		l.Printf("Can't get AtomNetWmWindowType from %s: %s", w, err)
	}
	// Update informations
	b.Name = w.Name()
	b.NameX = utf16.Encode([]rune(b.Name))

	insertNewBox(panel, b)
}

func insertNewBox(panel, b *Box) {
	l.Print("insertNewBox: ", b.Window)
	w := b.Window
	// Set window attributes
	if b.Type == BoxTypeWindow {
		w.SetBorderWidth(cfg.BorderWidth)
		w.SetBorderColor(cfg.NormalBorderColor)
	}
	// Grab right mouse buttons for WM actions 
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 3, xgb.ButtonMaskAny)
	// Add window to found parentBox
	w.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmapNotify due to reparenting
	w.Reparent(panel.Window, 0, 0)
	w.SetEventMask(EventMask)
	// Update geometry of windows in panel
	panel.Children.PushBack(b)
	tile(panel)
	// Show the window
	w.Map()
}

func winFocus(w Window) {
	l.Print("Focusing window: ", w)
	if currentDesk.Window == w {
		currentDesk.Window.SetInputFocus()
		currentPanel = currentDesk
	}
	// Iterate over all boxes in current desk
	panel := currentPanel
	bi := currentDesk.Children.FrontIter(true)
	for b := bi.Next(); b != nil; b = bi.Next() {
		if b.Type == BoxTypeWindow {
			if b.Window == w {
				w.SetBorderColor(cfg.FocusedBorderColor)
				currentPanel = panel
				w.SetInputFocus()
			} else {
				b.Window.SetBorderColor(cfg.NormalBorderColor)
			}
		} else {
			panel = b
			if b.Window == w {
				currentPanel = panel
				b.Window.SetInputFocus()
			}
		}
	}
}

// Panel is a box with transparent window in which a real windows or other
// panels are placed.
func newPanel(typ BoxType, parent *Box) {
	l.Print("newPanel")
	if parent.Type == BoxTypeWindow {
		panic("Can't create panel in BoxTypeWindow box")
	}
	w := NewWindow(parent.Window, Geometry{0, 0, 1, 1},
		xgb.WindowClassInputOutput, 0)
	w.SetName("mdtwm panel")
	insertNewBox(parent, NewBox(typ, w))
}
