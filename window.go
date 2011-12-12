package main

import (
	"bytes"
	"reflect"
	"unsafe"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Window interface {
	String() string
	Id() xgb.Id

	// Properties
	Prop(prop xgb.Id, max uint32) (*xgb.GetPropertyReply, error)
	ChangeProp(mode byte, prop, typ xgb.Id, data interface{})
	Name() string
	SetName(name string)
	Class() (instance, class string)
	SetClass(instance, class string)

	// Configuration
	Configure(mask uint16, vals ...uint32)
	SetGeometry(g Geometry)
	SetPosition(x, y int16)
	SetSize(width, height int16)
	SetBorderWidth(width int16)

	// Attributes
	Attrs() *xgb.GetWindowAttributesReply
	ChangeAttrs(mask uint32, vals ...uint32)
	SetBorderColor(pixel uint32)
	SetEventMask(mask uint32)

	// Other
	Geometry() Geometry
	Map()
	SetInputFocus()
	Reparent(parent Window, x, y int16)
	GrabButton(ownerEvents bool, eventMask uint16,
		pointerMode, keyboardMode byte,
		confineTo Window, cursor xgb.Id,
		button byte, modifiers uint16)
	UngrabButton(button byte, modifiers uint16)
	GrabKey(ownerEvents bool, modifiers uint16, key,
		pointerMode, keyboardMode byte)
	UngrabKey(key byte, modifiers uint16)
	ChangeSaveSet(mode byte)
}

type RawWindow xgb.Id

// Creates unmaped window
func NewRawWindow(parent Window, g Geometry, class uint16,
	mask uint32, vals ...uint32) RawWindow {
	id := conn.NewId()
	conn.CreateWindow(
		xgb.WindowClassCopyFromParent,
		id, parent.Id(),
		g.X, g.Y, Uint16(g.W), Uint16(g.H), Uint16(g.B),
		class, xgb.WindowClassCopyFromParent,
		mask, vals,
	)
	return RawWindow(id)
}

func (w RawWindow) String() string {
	return w.Name()
}

func (w RawWindow) Id() xgb.Id {
	return xgb.Id(w)
}

// Base methods from XGB

func (w RawWindow) Prop(prop xgb.Id, max uint32) (*xgb.GetPropertyReply, error) {
	return conn.GetProperty(false, w.Id(), prop, xgb.GetPropertyTypeAny, 0, max)
}

func (w RawWindow) ChangeProp(mode byte, prop, typ xgb.Id, data interface{}) {
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
	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int, reflect.Uint:
		p := reflect.New(d.Type())
		p.Elem().Set(d)
		d = p // now d is a pointer to an integer
		fallthrough
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
		panic("Property data should be an integer, a string, a pointer or a slice")
	}
	if format > 255 {
		panic("format > 255")
	}
	conn.ChangeProperty(mode, w.Id(), prop, typ, byte(format*8), content)
}

func (w RawWindow) Attrs() *xgb.GetWindowAttributesReply {
	a, err := conn.GetWindowAttributes(w.Id())
	if err != nil {
		l.Fatalf("Can't get attributes of window %s: %s", w, err)
	}
	return a
}

func (w RawWindow) ChangeAttrs(mask uint32, vals ...uint32) {
	conn.ChangeWindowAttributes(w.Id(), mask, vals)
}

func (w RawWindow) Configure(mask uint16, vals ...uint32) {
	conn.ConfigureWindow(w.Id(), mask, vals)
}

func (w RawWindow) Geometry() Geometry {
	g, err := conn.GetGeometry(w.Id())
	if err != nil {
		l.Fatalf("Can't get geometry of window %s: %s", w, err)

	}
	return Geometry{
		g.X, g.Y,
		Int16(g.Width), Int16(g.Height),
		Int16(g.BorderWidth),
	}
}

func (w RawWindow) Map() {
	conn.MapWindow(w.Id())
}

func (w RawWindow) Reparent(parent Window, x, y int16) {
	conn.ReparentWindow(w.Id(), parent.Id(), x, y)
}
func (w RawWindow) ChangeSaveSet(mode byte) {
	conn.ChangeSaveSet(mode, w.Id())
}

func (w RawWindow) GrabButton(ownerEvents bool, eventMask uint16,
	pointerMode, keyboardMode byte, confineTo Window, cursor xgb.Id,
	button byte, modifiers uint16) {
	conn.GrabButton(ownerEvents, w.Id(), eventMask, pointerMode, keyboardMode,
		confineTo.Id(), cursor, button, modifiers)
}

func (w RawWindow) UngrabButton(button byte, modifiers uint16) {
	conn.UngrabButton(button, w.Id(), modifiers)
}

func (w RawWindow) GrabKey(ownerEvents bool, modifiers uint16,
	key, pointerMode, keyboardMode byte) {
	conn.GrabKey(ownerEvents, w.Id(), modifiers, key, pointerMode, keyboardMode)
}

func (w RawWindow) UngrabKey(key byte, modifiers uint16) {
	conn.UngrabKey(key, w.Id(), modifiers)
}

// Utility methods

func (w RawWindow) SetGeometry(g Geometry) {
	w.Configure(
		xgb.ConfigWindowX|xgb.ConfigWindowY|
			xgb.ConfigWindowWidth|xgb.ConfigWindowHeight|
			xgb.ConfigWindowBorderWidth,
		uint32(g.X), uint32(g.Y),
		uint32(Pint16(g.W)), uint32(Pint16(g.H)),
		uint32(g.B),
	)
}

func (w RawWindow) SetPosition(x, y int16) {
	w.Configure(xgb.ConfigWindowX|xgb.ConfigWindowY, uint32(x), uint32(y))
}

func (w RawWindow) SetSize(width, height int16) {
	w.Configure(xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
		uint32(Pint16(width)), uint32(Pint16(height)))
}

func (w RawWindow) SetBorderWidth(width int16) {
	w.Configure(xgb.ConfigWindowBorderWidth, uint32(width))
}

func (w RawWindow) SetBorderColor(pixel uint32) {
	w.ChangeAttrs(xgb.CWBorderPixel, pixel)
}

func (w RawWindow) SetInputFocus() {
	conn.SetInputFocus(xgb.InputFocusPointerRoot, w.Id(), xgb.TimeCurrentTime)
}

func (w RawWindow) Name() string {
	// We prefer utf8 version
	if p, err := w.Prop(AtomNetWmName, 128); err == nil && len(p.Value) > 0 {
		return string(p.Value)
	}
	if p, err := w.Prop(xgb.AtomWmName, 128); err == nil && len(p.Value) > 0 {
		return string(p.Value)
	}
	return ""
}

func (w RawWindow) SetName(name string) {
	w.ChangeProp(xgb.PropModeReplace, xgb.AtomWmName, xgb.AtomString, name)
	w.ChangeProp(xgb.PropModeReplace, AtomNetWmName, AtomUtf8String, name)
}

func (w RawWindow) Class() (instance, class string) {
	p, err := w.Prop(xgb.AtomWmClass, 128)
	if err != nil {
		return
	}
	v := p.Value
	i := bytes.IndexByte(v, 0)
	if i == -1 {
		return
	}
	return string(v[:i]), string(v[i+1 : len(v)-1])
}

func (w RawWindow) SetClass(instance, class string) {
	v := make([]byte, 0, len(instance) + len(class) + 2)
	v = append(v, instance...)
	v = append(v, 0)
	v = append(v, class...)
	w.ChangeProp(xgb.PropModeReplace,xgb.AtomWmClass, xgb.AtomString, v)
}

func (w RawWindow) SetEventMask(mask uint32) {
	w.ChangeAttrs(xgb.CWEventMask, mask)
}
