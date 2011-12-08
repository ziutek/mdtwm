package main

import (
	"math"
	"x-go-binding.googlecode.com/hg/xgb"
)

const (
	WindowEventMask = xgb.EventMaskPropertyChange |
		xgb.EventMaskStructureNotify |
		xgb.EventMaskFocusChange

	FrameEventMask = xgb.EventMaskButtonPress |
		xgb.EventMaskButtonRelease |
		//xgb.EventMaskPointerMotion |
		xgb.EventMaskExposure | // window needs to be redrawn
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
	attr := w.Attrs()
	if attr.MapState != xgb.MapStateViewable || attr.OverrideRedirect {
		return
	}
	wm_type := propReplyAtoms(w.Prop(AtomNetWmWindowType, math.MaxUint32))
	if wm_type.Contains(AtomNetWmWindowTypeDock) {
		l.Printf("Window %v is of type dock")
	}
	// Grab left and right mouse buttons for click to focus/rasie
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 1, xgb.ButtonMaskAny)
	w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 3, xgb.ButtonMaskAny)

	b := NewBox()
	x, y, width, height := w.Geometry()
	b.Window = w
	b.Frame = CreateWindow(
		root, // TODO: Obtain parent box and set parent of frame from it
		x-int16(cfg.BorderWidth), y-int16(cfg.BorderWidth),
		width+cfg.BorderWidth*2, height+cfg.BorderWidth*2,
		xgb.WindowClassInputOutput,
		xgb.CWOverrideRedirect|xgb.CWEventMask,
		1, FrameEventMask,
	)
	// Check it the window should be floating
	if cfg.Float.Contains(w.Class()) ||
		wm_type.Contains(AtomNetWmWindowTypeDialog) ||
		wm_type.Contains(AtomNetWmWindowTypeUtility) ||
		wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
		wm_type.Contains(AtomNetWmWindowTypeSplash) {
		b.Float = true
	}
	// TODO: FrameMask for frame window not for child
	w.SetEventMask(WindowEventMask)

	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)

	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	w.SetGeometry(x, 0, 500, 700)
	x += 504

	currentDesk.Childs.PushFront(b)
}

func winFocus(w Window) {
	l.Print("Focusing window: ", w)
	for bi := currentDesk.Childs.FrontIter(true); !bi.Done(); {
		b := bi.Next()
		if b.Window == w {
			b.Frame.SetBorderColor(cfg.FocusedBorderColor)
			w.SetInputFocus()
		} else {
			b.Frame.SetBorderColor(cfg.NormalBorderColor)
		}
	}
}
