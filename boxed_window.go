package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"math"
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
	// Set initial geometry for new window
	if g, ok := w.Geometry(); ok {
		// Place new window on center of its parent
		bb := cfg.BorderWidth * 2
		b.width = g.W + bb
		b.height = g.H + bb
		_, _, w, h := b.parent.PosSize()
		b.x = (w - b.width) / 2
		b.y = (h - b.height) / 2
		g.X = b.x
		g.Y = b.y
		g.B = cfg.BorderWidth
		b.w.SetGeometry(g)
	}
	// Set window normal hints
	b.SetHints(propToHints(w.Prop(xgb.AtomWmNormalHints, math.MaxUint32)))
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
	b.w.SetGeometry(Geometry{
		x, y,
		width - bb, height - bb,
		cfg.BorderWidth,
	})
}

func (b *BoxedWindow) SetFocus(f bool, t xgb.Timestamp) {
	if f {
		b.w.SetInputFocus(t)
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

func (b *BoxedWindow) UpdateNetWmDesktop() {
	b.w.ChangeProp(xgb.PropModeReplace, AtomNetWmDesktop, xgb.AtomWindow,
		currentDesk.Window())
}

func (b *BoxedWindow) Protocols() IdList {
	p := b.Window().Prop(AtomWmProtocols, math.MaxUint32)
	return atomList(p)
}

func (b *BoxedWindow) SendMessage(typ xgb.Id, w Window) {
	m := xgb.ClientMessageEvent{
		Format: 32,
		Window: w.Id(),
		Type:   AtomWmProtocols,
	}
	m.Data.Data32[0] = uint32(typ)
	m.Data.Data32[1] = uint32(xgb.TimeCurrentTime)

	b.w.Send(false, xgb.EventMaskNoEvent, m)
}

func propToHints(prop *xgb.GetPropertyReply) (h Hints) {
	p := prop32(prop)
	if len(p) > 5 {
		h.MinW = int16(p[5])
	}
	if len(p) > 6 {
		h.MinH = int16(p[6])
	}
	if len(p) > 7 {
		h.MaxW = int16(p[7])
	}
	if len(p) > 8 {
		h.MaxH = int16(p[8])
	}
	if len(p) > 9 {
		h.IncW = int16(p[9])
	}
	if len(p) > 10 {
		h.IncH = int16(p[10])
	}
	if len(p) > 12 {
		h.MinAspect[0] = int16(p[11])
		h.MinAspect[1] = int16(p[12])
	}
	if len(p) > 14 {
		h.MaxAspect[0] = int16(p[13])
		h.MaxAspect[1] = int16(p[14])
	}
	if len(p) > 15 {
		h.BaseW = int16(p[15])
	}
	if len(p) > 16 {
		h.BaseH = int16(p[16])
	}
	if len(p) > 17 {
		h.Gravity = byte(p[17])
	}
	return
}
