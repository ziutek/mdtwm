package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"x-go-binding.googlecode.com/hg/xgb"
)

var (
	conn        *xgb.Conn
	screen      *xgb.ScreenInfo
	root        Window
	currentDesk *Box // desk in mdtwm means workspace
	allDesks    BoxList

	l = log.New(os.Stderr, "mdtwm: ", 0)
)

func main() {
	signals()
	connect()
	loadConfig()
	setupAtoms()
	wmProperties()
	grabKeys()
	manageExistingWindows()
	eventLoop()
}

func signals() {
	l.Print("signals")
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
	display := os.Getenv("DISPLAY")
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

	// Setup current desk (for now there is only one desk)
	currentDesk = NewBox()
	currentDesk.Frame = root
	currentDesk.Window = root

	allDesks = NewBoxList()
	allDesks.PushFront(currentDesk)
}

func wmProperties() {
	l.Print("wmProperties")
	// Spported atoms
	/* root.ChangeProp(xgb.PropModeReplace, AtomNetSupported,
	xgb.AtomAtom,	...) */
	root.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, &root)
	root.ChangeProp(xgb.PropModeReplace, AtomNetWmName, AtomUtf8String, "mdtwm")
}

func manageExistingWindows() {
	l.Print("manageExistingWindows")
	tr, err := conn.QueryTree(root.Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		w := Window(id)
		winAdd(w)
		w.Map()
	}
}

func grabKeys() {
	l.Print("grabKeys")
	// Win + Return
	root.GrabKey(true, xgb.ModMask4, 36, xgb.GrabModeAsync, xgb.GrabModeAsync)
}

func eventLoop() {
	l.Print("eventLoop")
	root.SetEventMask(
		xgb.EventMaskSubstructureRedirect |
			xgb.EventMaskStructureNotify |
			//xgb.EventMaskPointerMotion |
			xgb.EventMaskPropertyChange |
			xgb.EventMaskEnterWindow,
	)
	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			l.Print("WaitForEvent error: ", err)
			continue
		}
		switch ev := event.(type) {
		case xgb.MapRequestEvent:
			mapRequest(ev)
		case xgb.EnterNotifyEvent:
			enterNotify(ev)
		case xgb.DestroyNotifyEvent:
			destroyNotify(ev)
		case xgb.ConfigureNotifyEvent:
			configureNotify(ev)
		case xgb.ConfigureRequestEvent:
			configureRequest(ev)
		case xgb.KeyPressEvent:
			keyPress(ev)
		case xgb.ButtonPressEvent:
			buttonPress(ev)
		default:
			l.Print("Unknown event: ", ev)
		}
	}
}
