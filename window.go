package main

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Window xgb.Id

// Creates unmaped window
func NewWindow(parent Window, g Geometry, class uint16,
	mask uint32, vals ...uint32) Window {
	id := conn.NewId()
	conn.CreateWindow(
		xgb.WindowClassCopyFromParent,
		id, parent.Id(),
		g.X, g.Y, Uint16(g.W), Uint16(g.H), Uint16(g.B),
		class, xgb.WindowClassCopyFromParent,
		mask, vals,
	)
	return Window(id)
}

func (w Window) String() string {
	return fmt.Sprint(w.Id())
}

func (w Window) Id() xgb.Id {
	return xgb.Id(w)
}

func (w Window) Map() {
	conn.MapWindow(w.Id())
}

func (w Window) Reparent(parent Window, x, y int16) {
	conn.ReparentWindow(w.Id(), parent.Id(), x, y)
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

func (w Window) SetInputFocus() {
	conn.SetInputFocus(xgb.InputFocusPointerRoot, w.Id(), xgb.TimeCurrentTime)
}

// Configuration

func (w Window) Configure(mask uint16, vals ...interface{}) {
	data := make([]uint32, len(vals))
	for i, v := range vals {
		r := reflect.ValueOf(v)
		switch r.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			data[i] = uint32(r.Uint())
		case reflect.Int8, reflect.Int16, reflect.Int32:
			data[i] = uint32(r.Int())
		default:
			panic(fmt.Sprintf(
				"vals[%d] type is %s; accepted: int8-32, uint8-32 ",
				i, r.Type(),
			))
		}

	}
	conn.ConfigureWindow(w.Id(), mask, data)
}

func (w Window) Geometry() Geometry {
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

func (w Window) SetGeometry(g Geometry) {
	w.Configure(
		xgb.ConfigWindowX|xgb.ConfigWindowY|
			xgb.ConfigWindowWidth|xgb.ConfigWindowHeight|
			xgb.ConfigWindowBorderWidth,
		g.X, g.Y, Pint16(g.W), Pint16(g.H), g.B,
	)
}

func (w Window) SetPosition(x, y int16) {
	w.Configure(xgb.ConfigWindowX|xgb.ConfigWindowY, x, y)
}

func (w Window) SetSize(width, height int16) {
	w.Configure(xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
		Pint16(width), Pint16(height))
}

func (w Window) SetBorderWidth(width int16) {
	w.Configure(xgb.ConfigWindowBorderWidth, width)
}

// Attributes

func (w Window) Attrs() *xgb.GetWindowAttributesReply {
	a, err := conn.GetWindowAttributes(w.Id())
	if err != nil {
		l.Fatalf("Can't get attributes of window %s: %s", w, err)
	}
	return a
}

func (w Window) ChangeAttrs(mask uint32, vals ...uint32) {
	conn.ChangeWindowAttributes(w.Id(), mask, vals)
}

func (w Window) SetBorderColor(pixel uint32) {
	w.ChangeAttrs(xgb.CWBorderPixel, pixel)
}

func (w Window) SetBgColor(pixel uint32) {
	w.ChangeAttrs(xgb.CWBackPixel, pixel)
}

func (w Window) SetEventMask(mask uint32) {
	w.ChangeAttrs(xgb.CWEventMask, mask)
}

// Properities

func (w Window) Prop(prop xgb.Id, max uint32) (*xgb.GetPropertyReply, error) {
	return conn.GetProperty(false, w.Id(), prop, xgb.GetPropertyTypeAny, 0, max)
}

func (w Window) ChangeProp(mode byte, prop, typ xgb.Id, val interface{}) {
	if val == nil {
		panic("nil property")
	}
	var (
		format  int
		content []byte
	)
	d := reflect.ValueOf(val)
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
		panic("Property value should be an integer, a string, a pointer or a slice")
	}
	if format > 255 {
		panic("format > 255")
	}
	conn.ChangeProperty(mode, w.Id(), prop, typ, byte(format*8), content)
}

// Class properity is implemented in Window because it is needed to check if
// WM need to ignore some window
func (w Window) Class() (instance, class string) {
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
