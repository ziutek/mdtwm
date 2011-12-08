package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

func mapRequest(e xgb.MapRequestEvent) {
	l.Print("MapRequestEvent")

	w := Window(e.Window)
	if allDesks.BoxByWindow(w) == nil {
		winAdd(w)
	}
	w.Map()
	winFocus(w)
}

func enterNotify(e xgb.EnterNotifyEvent) {
	if e.Mode != xgb.NotifyModeNormal {
		return
	}
	w := Window(e.Event)
	if currentDesk.Childs.BoxByWindow(w) != nil {
		winFocus(w)
	}
}

func destroyNotify(e xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent")
}

func configureNotify(e xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent")
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent")
}

func keyPress(ev xgb.KeyPressEvent) {
	l.Print("KeyPressEvent")
}

func buttonPress(ev xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent")
}

