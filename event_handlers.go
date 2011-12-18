package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

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
	unmapNotify(xgb.UnmapNotifyEvent{Event: e.Event, Window: e.Window})
}

func unmapNotify(e xgb.UnmapNotifyEvent) {
	l.Print("xgb.UnmapNotifyEvent: ", Window(e.Event), e)
	w := Window(e.Event)
	if b := root.Children().BoxByWindow(w, true); b != nil {
		b.Parent().Remove(b)
	}
}

func configureNotify(e xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent: ", Window(e.Window))
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	if root.Children().BoxByWindow(w, true) == nil {
		// Unmanaged window - execute its request
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
	} else {
		l.Print("ConfigureRequestEvent from managed window")
	}
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
	if _, ok := currentBox.(ParentBox); ok {
		// For now, we don't move panels
		return
	}
	move.b = currentBox
	move.x, move.y = e.RootX, e.RootY
	move.b.Parent().Remove(move.b)
}

func buttonRelease(e xgb.ButtonReleaseEvent) {
	l.Println("ButtonReleaseEvent:", Window(e.Event), Window(e.Child))
	if move.b == nil {
		return
	}
	// Border coolor will be set properly by EnterNotify event.
	move.b.Window().SetBorderColor(cfg.NormalBorderColor)
	currentPanel().Insert(move.b)
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
