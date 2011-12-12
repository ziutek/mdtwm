package main

import (
	"math"
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

func manage(w Window, panel *PanelBox) {
	// We need to use b because b implements caching (TODO)
	b := NewWindowBox(w)
	l.Print("manage: ", b)
	_, class := b.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if allDesks.BoxByWindow(b, true) != nil {
		l.Printf("  %s - alredy managed", b)
		return
	}
	attr := b.Attrs()
	// Don't manage and don't map if unvisible or has OverrideRedirect flag
	if attr.MapState != xgb.MapStateViewable || attr.OverrideRedirect {
		return
	}
	// Check window type
	p, err := b.Prop(AtomNetWmWindowType, math.MaxUint32)
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
	// Insert new box in a panel
	panel.Insert(b)
}
