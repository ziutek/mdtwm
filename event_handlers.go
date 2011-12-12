package main

import (
	"x-go-binding.googlecode.com/hg/xgb"
)

func mapRequest(e xgb.MapRequestEvent) {
	w := RawWindow(e.Window)
	l.Print("MapRequestEvent: ", w)
	manage(w, currentPanel)
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
}

func configureNotify(e xgb.ConfigureNotifyEvent) {
	l.Print("ConfigureNotifyEvent: ", RawWindow(e.Window))
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", RawWindow(e.Window))
}

func keyPress(e xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", RawWindow(e.Event))
}

func buttonPress(e xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent: ", RawWindow(e.Event))
}

