package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"reflect"
	"time"
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
		l.Print("*** Unhandled event: ", reflect.TypeOf(e))
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
	removeWindow(Window(e.Event))
}

func unmapNotify(e xgb.UnmapNotifyEvent) {
	l.Print("xgb.UnmapNotifyEvent: ", Window(e.Event), e)
	removeWindow(Window(e.Event))
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
	w.SendEvent(false, xgb.EventMaskStructureNotify, cne)
}

func keyPress(e xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", Window(e.Event))
	// Simply run terminal
	if e.State == cfg.ModMask {
		cmd, ok := cfg.Keys[e.Detail]
		if !ok {
			l.Print("Unhandled key: ", e.Detail)
		}
		if err := cmd.Run(); err != nil {
			l.Printf("cmd(%s): %s", cmd.Param, err)
		}
	}
}

var move struct {
	b    Box
	x, y int16
}

func buttonPress(e xgb.ButtonPressEvent) {
	l.Println("ButtonPressEvent:", Window(e.Event), Window(e.Child))
	// Wait for double-click
	time.Sleep(200 * time.Millisecond)

	if _, ok := currentBox.(ParentBox); ok {
		// For now, we don't move panels
		return
	}
	conn.ChangeActivePointerGrab(cfg.MoveCursor, xgb.TimeCurrentTime,
		xgb.EventMaskButtonPress|xgb.EventMaskButtonRelease)
	move.b = currentBox
	move.x, move.y = e.RootX, e.RootY
}

func buttonRelease(e xgb.ButtonReleaseEvent) {
	l.Println("ButtonReleaseEvent:", Window(e.Event), Window(e.Child))
	if move.b == nil {
		return
	}
	// Border coolor will be set properly by EnterNotify event.
	move.b.Window().SetBorderColor(cfg.NormalBorderColor)
	// We need to save a current panel before remove a box from its panel,
	// beacause this box may be a current box
	cp := currentPanel()
	move.b.Parent().Remove(move.b)
	cp.Insert(move.b)
	l.Print("  ", currentPanel())
	//changeFocusTo(move.b.Window()) // Always set moved window focused
	move.b = nil
}

func motionNotify(e xgb.MotionNotifyEvent) {
	l.Print("xgb.MotionNotifyEvent: ", Window(e.Event), Window(e.Child))
	if move.b == nil {
		return
	}
	dx, dy := e.RootX-move.x, e.RootY-move.y
	x, y, w, h := move.b.PosSize()
	move.b.SetPosSize(x+dx, y+dy, w, h)
	move.x += dx
	move.y += dy
}
