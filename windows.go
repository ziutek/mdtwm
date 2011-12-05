package main

import (
	"container/list"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Window xgb.Id

func (w Window) String() string {
	return w.Name()
}

func (w Window) Id() xgb.Id {
	return xgb.Id(w)
}

func (w Window) SetAttrs(mask uint32, vals ...uint32) {
	conn.ChangeWindowAttributes(w.Id(), mask, vals)
}

func (w Window) SetBorderColor(pixel uint32) {
	 w.SetAttrs(xgb.CWBorderPixel, pixel)
}

func (w Window) SetBorderWidth(width uint32) {
	conn.ConfigureWindow(w.Id(), xgb.ConfigWindowBorderWidth, []uint32{width})
}

func (w Window) SetInputFocus() {
	conn.SetInputFocus(xgb.InputFocusPointerRoot, w.Id(), xgb.TimeCurrentTime)
}

func (w Window) Geometry() (x, y, width, height int) {
	g, err := conn.GetGeometry(w.Id())
	if err != nil {
		l.Fatal("Can't get geometry of window %v: %v", w, err)
	}
	return int(g.X), int(g.Y), int(g.Width), int(g.Height)
}

func (w Window) Attrs() *xgb.GetWindowAttributesReply {
	a, err := conn.GetWindowAttributes(w.Id())
    if err != nil {
        l.Fatalf("Can't get attributes of window %v: %v", w, err)
    }
	return a
}

func (w Window) Property(prop xgb.Id, max uint32) *xgb.GetPropertyReply {
	p, err := conn.GetProperty(false, w.Id(), prop, xgb.GetPropertyTypeAny, 0,
		max)
	if err != nil {
		l.Fatalf("Can't get property of window %v: %v", w, err)
	}
	return p
}

func (w Window) Name() string {
	p := w.Property(xgb.AtomWmName, 128)
	return string(p.Value)
}

func (w Window) Class() string {
    p := w.Property(xgb.AtomWmClass, 128)
    return string(p.Value)
}

func (w Window) Map() {
	conn.MapWindow(w.Id())
}

var windows = list.New()
