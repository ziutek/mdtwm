package main

import (
	"reflect"
	"unsafe"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Window xgb.Id

func (w Window) String() string {
	return w.Name()
}

func (w Window) Id() xgb.Id {
	return xgb.Id(w)
}

func (w Window) ChangeAttrs(mask uint32, vals ...uint32) {
	conn.ChangeWindowAttributes(w.Id(), mask, vals)
}

func (w Window) Configure(mask uint16, vals ...uint32) {
	conn.ConfigureWindow(w.Id(), mask, vals)
}

func (w Window) SetBorderColor(pixel uint32) {
	w.ChangeAttrs(xgb.CWBorderPixel, pixel)
}

func (w Window) SetBorderWidth(width uint32) {
	w.Configure(xgb.ConfigWindowBorderWidth, width)
}

func (w Window) SetInputFocus() {
	conn.SetInputFocus(xgb.InputFocusPointerRoot, w.Id(), xgb.TimeCurrentTime)
}

func (w Window) Geometry() (x, y int16, width, height uint16) {
	g, err := conn.GetGeometry(w.Id())
	if err != nil {
		l.Fatal("Can't get geometry of window %v: %v", w, err)
	}
	return g.X, g.Y, g.Width, g.Height
}

func (w Window) SetGeometry(x, y int16, width, height uint16) {
	w.Configure(
		xgb.ConfigWindowX|xgb.ConfigWindowY|
			xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
		uint32(x), uint32(y), uint32(width), uint32(height),
	)
}

func (w Window) Attrs() *xgb.GetWindowAttributesReply {
	a, err := conn.GetWindowAttributes(w.Id())
	if err != nil {
		l.Fatalf("Can't get attributes of window %v: %v", w, err)
	}
	return a
}

func (w Window) Prop(prop xgb.Id, max uint32) *xgb.GetPropertyReply {
	p, err := conn.GetProperty(false, w.Id(), prop, xgb.GetPropertyTypeAny, 0,
		max)
	if err != nil {
		l.Fatalf("Can't get property of window %v: %v", w, err)
	}
	return p
}

func (w Window) ChangeProp(mode byte, prop, typ xgb.Id, data interface{}) {
	if data == nil {
		panic("nil property")
	}
	var (
		format  int
		content []byte
	)
	d := reflect.ValueOf(data)
	switch d.Kind() {
	case reflect.String:
		format = 1
		content = []byte(d.String())
	case reflect.Ptr:
		format = int(d.Type().Elem().Size())
		length := format
		addr := unsafe.Pointer(d.Elem().UnsafeAddr())
		content = (*[1<<31 - 1]byte)(addr)[:length]
	case reflect.Slice:
		format = int(d.Type().Elem().Size())
		length := format * d.Len()
		addr := unsafe.Pointer(d.Index(0).UnsafeAddr())
		content = (*[1<<31 - 1]byte)(addr)[:length]
	default:
		panic("Property data should be a string, a pointer or a slice")
	}
	if format > 255 {
		panic("format > 255")
	}
	conn.ChangeProperty(mode, w.Id(), prop, typ, byte(format*8), content)
}

func (w Window) Name() string {
	p := w.Prop(xgb.AtomWmName, 128)
	return string(p.Value)
}

func (w Window) Class() string {
	p := w.Prop(xgb.AtomWmClass, 128)
	return string(p.Value)
}

func (w Window) Map() {
	conn.MapWindow(w.Id())
}

func (w Window) Reparent(parent Window, x, y int16) {
	conn.ReparentWindow(w.Id(), parent.Id(), x, y)
}

func (w Window) SetEventMask(mask uint32) {
	w.ChangeAttrs(xgb.CWEventMask, mask)
}

func (w Window) ChangeSaveSet(mode byte) {
	conn.ChangeSaveSet(mode, w.Id())
}
