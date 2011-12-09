package main

import (
	"fmt"
	"reflect"
	"unsafe"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Geometry struct {
	X, Y int16
	W, H uint16
}

func (g Geometry) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d)", g.X, g.Y, g.W, g.H)
}

func (g Geometry) Enlarge(i int) Geometry {
	d := i + i
	return Geometry{
		Int16(int(g.X) - i), Int16(int(g.Y) - i),
		Uint16(int(g.W) + d), Uint16(int(g.H) + d),
	}
}

func (g Geometry) EnlargeH(i int) Geometry {
	return Geometry{Int16(int(g.X) - i), g.Y, Uint16(int(g.W) + i + i), g.H}
}

func (g Geometry) EnlargeV(i int) Geometry {
	return Geometry{g.X, Int16(int(g.Y) - i), g.W, Uint16(int(g.H) + i + i)}
}

type Window xgb.Id

// Creates unmaped window with border == 0
func NewWindow(parent Window, g Geometry, class uint16,
	mask uint32, vals ...uint32) Window {
	id := conn.NewId()
	conn.CreateWindow(
		xgb.WindowClassCopyFromParent,
		id, parent.Id(),
		g.X, g.Y, g.W, g.H, 0,
		class, xgb.WindowClassCopyFromParent,
		mask, vals,
	)
	return Window(id)
}

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

func (w Window) SetBorderWidth(width uint16) {
	w.Configure(xgb.ConfigWindowBorderWidth, uint32(width))
}

func (w Window) SetInputFocus() {
	conn.SetInputFocus(xgb.InputFocusPointerRoot, w.Id(), xgb.TimeCurrentTime)
}

func (w Window) Geometry() Geometry {
	g, err := conn.GetGeometry(w.Id())
	if err != nil {
		l.Fatalf("Can't get geometry of window %s: %s", w, err)

	}
	return Geometry{g.X, g.Y, g.Width, g.Height}
}

func (w Window) SetGeometry(g Geometry) {
	w.Configure(
		xgb.ConfigWindowX|xgb.ConfigWindowY|
			xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
		uint32(g.X), uint32(g.Y), uint32(g.W), uint32(g.H),
	)
}

func (w Window) Attrs() *xgb.GetWindowAttributesReply {
	a, err := conn.GetWindowAttributes(w.Id())
	if err != nil {
		l.Fatalf("Can't get attributes of window %s: %s", w, err)
	}
	return a
}

func (w Window) Prop(prop xgb.Id, max uint32) (*xgb.GetPropertyReply, error) {
	return conn.GetProperty(false, w.Id(), prop, xgb.GetPropertyTypeAny, 0, max)
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
	// We prefer utf8 version
	if p, err := w.Prop(AtomNetWmName, 128); err != nil {
		return string(p.Value)
	}
	if p, err := w.Prop(xgb.AtomWmName, 128); err != nil {
		return string(p.Value)
	}
	return ""
}

func (w Window) SetName(name string) {
	w.ChangeProp(xgb.PropModeReplace, xgb.AtomWmName, xgb.AtomString, name)
	w.ChangeProp(xgb.PropModeReplace, AtomNetWmName, AtomUtf8String, name)
}

func (w Window) Class() string {
	p, err := w.Prop(xgb.AtomWmClass, 128)
	if err != nil {
		return ""
	}
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

func (w Window) GrabButton(ownerEvents bool, eventMask uint16,
	pointerMode, keyboardMode byte, confineTo Window, cursor xgb.Id,
	button byte, modifiers uint16) {
	conn.GrabButton(ownerEvents, w.Id(), eventMask, pointerMode, keyboardMode,
		confineTo.Id(), cursor, button, modifiers)
}

func (w Window) UngrabButton(button byte, modifiers uint16) {
	conn.UngrabButton(button, w.Id(), modifiers)
}

func (w Window) GrabKey(ownerEvents bool, modifiers uint16,
	key, pointerMode, keyboardMode byte) {
	conn.GrabKey(ownerEvents, w.Id(), modifiers, key, pointerMode, keyboardMode)
}

func (w Window) UngrabKey(key byte, modifiers uint16) {
	conn.UngrabKey(key, w.Id(), modifiers)
}
