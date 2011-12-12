package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

func mapRequest(e xgb.MapRequestEvent) {
	l.Print("MapRequestEvent: ", RawWindow(e.Window))
	w := RawWindow(e.Window)
	manage(w, currentPanel, false)
}

func enterNotify(e xgb.EnterNotifyEvent) {
	l.Print("EnterNotifyEvent: ", RawWindow(e.Event))
	if e.Mode != xgb.NotifyModeNormal {
		return
	}
	currentDesk.SetFocus(currentDesk.Id() == e.Event)
	// Iterate over all boxes in current desk
	bi := currentDesk.Children().FrontIter(true)
	for b := bi.Next(); b != nil; b = bi.Next() {
		b.SetFocus(b.Id() == e.Event)
	}

}

func destroyNotify(e xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent: ", RawWindow(e.Event))
	unmapNotify(xgb.UnmapNotifyEvent{Event: e.Event, Window: e.Window})
}

func unmapNotify(e xgb.UnmapNotifyEvent) {
	l.Print("xgb.UnmapNotifyEvent: ", RawWindow(e.Event), e)
	w := RawWindow(e.Event)
	if b := allDesks.BoxByWindow(w, true); b != nil {
		b.Parent().Remove(b)
	}
}

func configureNotify(e xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent: ", RawWindow(e.Window))
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", RawWindow(e.Window))
	w := RawWindow(e.Window)
	if allDesks.BoxByWindow(w, true) == nil {
		// Unmanaged window - execute its request
		mask := (xgb.ConfigWindowX|xgb.ConfigWindowY|
			xgb.ConfigWindowWidth|xgb.ConfigWindowHeight|
			xgb.ConfigWindowBorderWidth|xgb.ConfigWindowSibling|
			xgb.ConfigWindowStackMode) & e.ValueMask
		v := make([]interface{}, 0, 7)
		if mask & xgb.ConfigWindowX != 0 {
			v = append(v, e.X)
		}
		if mask & xgb.ConfigWindowY != 0 {
			v = append(v, e.Y)
		}
		if mask & xgb.ConfigWindowWidth != 0 {
			v = append(v, e.Width)
		}
		if mask & xgb.ConfigWindowHeight!= 0 {
			v = append(v, e.Height)
		}
		if mask & xgb.ConfigWindowBorderWidth!= 0 {
			v = append(v, e.BorderWidth)
		}
		if mask & xgb.ConfigWindowSibling != 0 {
			v = append(v, e.Sibling)
		}
		if mask & xgb.ConfigWindowStackMode != 0 {
			v = append(v, e.StackMode)
		}
		w.Configure(mask, v...)
	} else {
		l.Print("  alredy managed")
	}
}

func keyPress(e xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", RawWindow(e.Event))
}

func buttonPress(e xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent: ", RawWindow(e.Event))
}

