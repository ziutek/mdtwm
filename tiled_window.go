package main

// Box for APP window
type TiledWindow struct {
	commonBox
}

// Warning! This function modifies some properities of window w.
func NewTiledWindow(w Window) *TiledWindow {
	var t TiledWindow
	t.init(w)
	t.w.SetBorderWidth(cfg.BorderWidth)
	t.w.SetBorderColor(cfg.NormalBorderColor)
	t.w.SetEventMask(boxEventMask)
	t.grabInput(root.Window())
	return &t
}

func (t *TiledWindow) SetPosSize(x, y, width, height int16) {
	bb := 2 * cfg.BorderWidth
	t.w.SetGeometry(Geometry{x, y, width - bb, height - bb, cfg.BorderWidth})
}

func (t *TiledWindow) SetFocus(f bool) {
	if f {
		currentPanel = t.parent
		t.w.SetInputFocus()
		t.w.SetBorderColor(cfg.FocusedBorderColor)
	} else {
		t.w.SetBorderColor(cfg.NormalBorderColor)
	}
}
