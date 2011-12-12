package main

// Box for APP window
type WindowBox struct {
	commonBox
}

func NewWindowBox(w Window) *WindowBox {
	var b WindowBox
	b.init(w)
	b.SetBorderWidth(cfg.BorderWidth)
	b.SetBorderColor(cfg.NormalBorderColor)
	return &b
}

func (w *WindowBox) SetPosSize(x, y, width, height int16) {
	bb := 2 * cfg.BorderWidth
	w.SetGeometry(Geometry{x, y, width - bb, height - bb, cfg.BorderWidth})
}

func (w *WindowBox) SetFocus(f bool) {
	if f {
		currentPanel = w.parent
		w.Window.SetInputFocus()
		w.SetBorderColor(cfg.FocusedBorderColor)
	} else {
		w.Window.SetBorderColor(cfg.NormalBorderColor)
	}
}
