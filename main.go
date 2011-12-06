package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"x-go-binding.googlecode.com/hg/xgb"
)

var (
	display string
	conn    *xgb.Conn
	screen  *xgb.ScreenInfo
	root    Window

	l = log.New(os.Stderr, "mdtwm: ", 0)
)

func main() {
	signals()
	connect()
	loadConfig()
	grabKeys()
	manageExistingWindows()
	eventLoop()
}

func signals() {
	go func() {
		for sig := range signal.Incoming {
			if s, ok := sig.(os.UnixSignal); ok {
				switch s {
				case syscall.SIGTERM, syscall.SIGINT:
					os.Exit(0)
				}
			}
			l.Printf("Signal %v received and ignored", sig)
		}
	}()
}

func connect() {
	l.Print("init")
	display = os.Getenv("DISPLAY")
	if display == "" {
		l.Fatal("There is not DISPLAY environment variable")
	}
	var err error
	conn, err = xgb.Dial(display)
	if err != nil {
		l.Fatalf("Can't connect to %s display: %v", display, err)
	}
	screen = conn.DefaultScreen()
	root = Window(screen.Root)
}

func manageExistingWindows() {
	tr, err := conn.QueryTree(root.Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		manageWindow(id)
	}
}

func manageWindow(id xgb.Id) {
	w := Window(id)
	l.Print("manageWindow: ", w)
	if cfg.Ignore.Contains(w.Class()) {
		return
	}
	wa := w.Attrs()
	if wa.OverrideRedirect || wa.MapState != xgb.MapStateViewable {
		return
	}
	w.EventMask(xgb.EventMaskPropertyChange | xgb.EventMaskStructureNotify |
		xgb.EventMaskFocusChange)
	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)

	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	windows.PushFront(w)
	w.Map()
}

func grabKeys() {
	l.Print("grabKeys")
	conn.GrabKey(true, root.Id(), xgb.ModMask4, 36, xgb.GrabModeAsync,
		xgb.GrabModeAsync) // Win + Return
}

func eventLoop() {
	l.Print("eventLoop")

	// Init event
	root.EventMask(
		xgb.EventMaskSubstructureRedirect |
		xgb.EventMaskStructureNotify |
		xgb.EventMaskSubstructureNotify |
		xgb.EventMaskPointerMotion |
		xgb.EventMaskPropertyChange |
		xgb.EventMaskEnterWindow,
	)
	// Event loop
	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			l.Print("WaitForEvent error: ", err)
			continue
		}
		switch ev := event.(type) {
		case *xgb.KeyPressEvent:
			l.Print("KeyPressEvent: ", ev)
		case *xgb.MapRequestEvent:
			l.Print("MapRequestEvent: ", ev)
		case *xgb.EnterNotifyEvent:
			enterNotifyHandler(ev)
		case *xgb.ButtonPressEvent:
			l.Print("ButtonPressEvent", ev)
		case *xgb.DestroyNotifyEvent:
			l.Print("DestroyNotifyEvent: ", ev)
		case *xgb.ConfigureNotifyEvent:
			l.Print("ConfigureNotifyEvent")
		case *xgb.ConfigureRequestEvent:
			l.Print("ConfigureRequestEvent")
		}
	}
}

func enterNotifyHandler(ev *xgb.EnterNotifyEvent) {
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

	for el := windows.Front(); el != nil; el = el.Next() {
		w := el.Value.(Window)
		l.Print(w)
		if w.Id() == ev.Event {
			w.SetBorderColor(cfg.FocusedBorderColor)
		}
	}
}
