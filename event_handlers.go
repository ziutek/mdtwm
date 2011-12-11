package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

func mapRequest(e xgb.MapRequestEvent) {
	l.Print("MapRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	winAdd(w, Window(e.Parent))
	w.Map()
	winFocus(w)
}

func enterNotify(e xgb.EnterNotifyEvent) {
	l.Print("EnterNotifyEvent: ", Window(e.Event))
	if e.Mode != xgb.NotifyModeNormal {
		return
	}
	w := Window(e.Event)
	if currentDesk.Children.BoxByWindow(w) != nil {
		winFocus(w)
	}
}

func destroyNotify(e xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent: ", Window(e.Event))
}

func configureNotify(e xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent: ", Window(e.Window))
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", Window(e.Window))
}

func keyPress(e xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", Window(e.Event))
}

func buttonPress(e xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent: ", Window(e.Event))
}

