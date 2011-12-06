package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

func mapRequest(ev xgb.MapRequestEvent) {
	l.Print("MapRequestEvent: ", ev)

	w := Window(ev.Window)
	if !windows.Contains(w) {
		winAdd(w)
	}
	w.Map()
	winFocus(w)
}

func enterNotify(ev xgb.EnterNotifyEvent) {
	l.Print("EnterNotifyEvent: ", ev)
	switch ev.Mode {
	case xgb.NotifyModeNormal:
		l.Print("NotifyModeNormal")
	case xgb.NotifyModeGrab:
		l.Print("NotifyModeGrab")
	case xgb.NotifyModeUngrab:
		l.Print("NotifyModeUngrab")
	case xgb.NotifyModeWhileGrabbed:
		l.Print("NotifyModeWhileGrabbed")
	default:
		l.Print("unknown notify mode")
	}

	w := Window(ev.Event)
	if windows.Contains(w) {
		winFocus(w)
	}
}

func destroyNotify(ev xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent: ", ev)
}

func configureNotify(ev xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent")
}

func configureRequest(ev xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent")
}

func keyPress(ev xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", ev)
}

func buttonPress(ev xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent", ev)
}

