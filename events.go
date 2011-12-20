package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

func handleEvent(event xgb.Event) {
	switch e := event.(type) {
	// *Request events
	case xgb.MapRequestEvent:
		mapRequest(e)
	case xgb.ConfigureRequestEvent:
		configureRequest(e)
	// *Notify events
	case xgb.EnterNotifyEvent:
		enterNotify(e)
	case xgb.UnmapNotifyEvent:
		unmapNotify(e)
	case xgb.DestroyNotifyEvent:
		destroyNotify(e)
	// Keyboard and mouse events
	case xgb.KeyPressEvent:
		keyPress(e)
	case xgb.ButtonPressEvent:
		buttonPress(e)
	case xgb.ButtonReleaseEvent:
		buttonRelease(e)
	case xgb.MotionNotifyEvent:
		motionNotify(e)
	default:
		l.Printf("*** Unhandled event: %T", e)
	}
}

func mapRequest(e xgb.MapRequestEvent) {
	l.Print("MapRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	manage(w, currentPanel(), false)
}

func enterNotify(e xgb.EnterNotifyEvent) {
	l.Print("EnterNotifyEvent: ", Window(e.Event))
	if e.Mode != xgb.NotifyModeNormal {
		return
	}
	changeFocusTo(Window(e.Event))
}

func destroyNotify(e xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent: ", Window(e.Event))
	removeWindow(Window(e.Event), false)
}

func unmapNotify(e xgb.UnmapNotifyEvent) {
	l.Print("xgb.UnmapNotifyEvent: ", Window(e.Event), e)
	removeWindow(Window(e.Event), false)
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	b := root.Children().BoxByWindow(w, true)
	if b == nil || b.Float() {
		// We accept request from floating windows.
		// Unmanaged window will be configured by manage() function so
		// now we can simply execute its request.
		mask := (xgb.ConfigWindowX | xgb.ConfigWindowY |
			xgb.ConfigWindowWidth | xgb.ConfigWindowHeight |
			xgb.ConfigWindowBorderWidth | xgb.ConfigWindowSibling |
			xgb.ConfigWindowStackMode) & e.ValueMask
		v := make([]interface{}, 0, 7)
		if mask&xgb.ConfigWindowX != 0 {
			v = append(v, e.X)
		}
		if mask&xgb.ConfigWindowY != 0 {
			v = append(v, e.Y)
		}
		if mask&xgb.ConfigWindowWidth != 0 {
			v = append(v, e.Width)
		}
		if mask&xgb.ConfigWindowHeight != 0 {
			v = append(v, e.Height)
		}
		if mask&xgb.ConfigWindowBorderWidth != 0 {
			v = append(v, e.BorderWidth)
		}
		if mask&xgb.ConfigWindowSibling != 0 {
			v = append(v, e.Sibling)
		}
		if mask&xgb.ConfigWindowStackMode != 0 {
			v = append(v, e.StackMode)
		}
		w.Configure(mask, v...)
		return
	}
	// Force box configuration
	g := b.Geometry()
	cne := xgb.ConfigureNotifyEvent{
		Event:        e.Window,
		Window:       e.Window,
		AboveSibling: xgb.WindowNone,
		X:            g.X,
		Y:            g.Y,
		Width:        Pint16(g.W),
		Height:       Pint16(g.H),
		BorderWidth:  Pint16(g.B),
	}
	w.Send(false, xgb.EventMaskStructureNotify, cne)
}
