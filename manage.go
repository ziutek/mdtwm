package main

import (
	"math"
	"x-go-binding.googlecode.com/hg/xgb"
)

var (
	windows = NewWindowList()
)

const (
	ChildEventMask = xgb.EventMaskPropertyChange |
		xgb.EventMaskStructureNotify |
		xgb.EventMaskFocusChange

	FrameEventMask = xgb.EventMaskButtonPress |
		xgb.EventMaskButtonRelease |
		//xgb.EventMaskPointerMotion |
		xgb.EventMaskExposure |  // window needs to be redrawn
		xgb.EventMaskStructureNotify | // frame gets destroyed
		xgb.EventMaskSubstructureRedirect | // app tries to resize itself
		xgb.EventMaskSubstructureNotify | // subwindows get notifies
		xgb.EventMaskEnterWindow
)

var x int16

func winAdd(w Window) {
	l.Print("manageWindow: ", w)
	if cfg.Ignore.Contains(w.Class()) {
		return
	}
	// Don't map if unvisible or has OverrideRedirect flag
	wa := w.Attrs()
	if wa.MapState != xgb.MapStateViewable || wa.OverrideRedirect {
		return
	}
	wm_type := propReplyToAtoms(w.Prop(AtomNetWmWindowType, math.MaxUint32))
	if wm_type.Contains(AtomNetWmWindowTypeDock) {
		l.Printf("Window %v is of type dock")
	}
	if wm_type.Contains(AtomNetWmWindowTypeDialog) ||
		wm_type.Contains(AtomNetWmWindowTypeUtility) ||
		wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
		wm_type.Contains(AtomNetWmWindowTypeSplash) {
		l.Printf("Window %v should be floating")
	}
	// TODO: FrameMask for frame window not for child
	w.SetEventMask(ChildEventMask | FrameEventMask)

	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)

	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	w.SetGeometry(x, 0, 500, 700)
	x += 504

	windows.PushFront(w)
}

func winFocus(w Window) {
	l.Print("Focusing window: ", w)
	for wi := windows.FrontIter(); !wi.Done(); {
		f := wi.Next()
		if f == w {
			f.SetBorderColor(cfg.FocusedBorderColor)
			f.SetInputFocus()
		} else {
			f.SetBorderColor(cfg.NormalBorderColor)
		}
	}
}
