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
	w := Window(e.Event)
	currentDesk.SetFocus(currentDesk.Window() == w)
	// Iterate over all boxes in current desk
	bi := currentDesk.Children().FrontIter(true)
	for b := bi.Next(); b != nil; b = bi.Next() {
		b.SetFocus(b.Window() == w)
	}

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

var moveStart struct {
	b     Box
	x, y int16
}

func buttonPress(e xgb.ButtonPressEvent) {
	l.Println("ButtonPressEvent:", Window(e.Event), Window(e.Child))
	l.Print("   ", currentBox)
	moveStart.b = currentBox
	moveStart.x, moveStart.y = e.RootX, e.RootY
}

func buttonRelease(e xgb.ButtonReleaseEvent) {
	l.Println("ButtonReleaseEvent:", Window(e.Event), Window(e.Child))
	if moveStart.b == nil {
		return
	}
	defer func(){moveStart.b = nil}()
	l.Println("   ", currentBox)
	return
	moveStart.b.Parent().Remove(moveStart.b)
	currentPanel().Insert(moveStart.b)
}

func motionNotify(e xgb.MotionNotifyEvent) {
	l.Print("xgb.MotionNotifyEvent: ", Window(e.Event), Window(e.Child))
}
