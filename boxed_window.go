package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

// Box for APP window
type BoxedWindow struct {
	commonBox
}

// Warning! This function modifies some properities of window w.
func NewBoxedWindow(w Window) *BoxedWindow {
	var b BoxedWindow
	w.ChangeSaveSet(xgb.SetModeInsert)
	b.init(w, xgb.EventMaskEnterWindow|xgb.EventMaskStructureNotify)
	b.w.SetBorderWidth(cfg.BorderWidth)
	b.w.SetBorderColor(cfg.NormalBorderColor)
	return &b
}

func (b *BoxedWindow) Geometry() Geometry {
	bb := cfg.BorderWidth * 2
	return Geometry{
		X: b.x, Y: b.y,
		W: b.width - bb, H: b.height - bb,
		B: cfg.BorderWidth,
	}
}

func (b *BoxedWindow) SetPosSize(x, y, width, height int16) {
	b.x, b.y, b.width, b.height = x, y, width, height
	bb := 2 * cfg.BorderWidth
	b.w.SetGeometry(Geometry{x, y, width - bb, height - bb, cfg.BorderWidth})
}

func (b *BoxedWindow) SetFocus(f bool) {
	if f {
		currentBox = b
		b.w.SetInputFocus()
		b.w.SetBorderColor(cfg.FocusedBorderColor)
	} else {
		b.w.SetBorderColor(cfg.NormalBorderColor)
	}
}

type WmState uint32
const (
	WmStateWithdrawn = WmState(iota)
	WmStateNormal
	WmStateIconic
)

func (b *BoxedWindow) SetWmState(state WmState) {
	data := []uint32{uint32(state), uint32(xgb.WindowNone)}
	b.w.ChangeProp(xgb.PropModeReplace, AtomWmState, AtomWmState, data)
}
