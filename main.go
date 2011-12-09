package main

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"x-go-binding.googlecode.com/hg/xgb"
)

var (
	conn   *xgb.Conn
	screen *xgb.ScreenInfo
	root   Window

	// Desk in mdtwm means workspace. Desk contains panels. Panel contains
	// windows. All they are described by Box structure.
	allDesks      BoxList
	currentDesk   *Box
	currentPanel  *Box
	currentWindow *Box

	l = log.New(os.Stderr, "mdtwm: ", 0)
)

func main() {
	signals()
	connect()
	setupAtoms()
	loadConfig()
	setupWm()
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
				case syscall.SIGWINCH:
					continue
				}
			}
			l.Printf("Signal %s received and ignored", sig)
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
		l.Fatalf("Can't connect to %s display: %s", display, err)
	}
	screen = conn.DefaultScreen()
	root = Window(screen.Root)
}

func setupWm() {
	l.Print("setupDesks")
	// Spported atoms
	/* root.ChangeProp(xgb.PropModeReplace, AtomNetSupported,
	xgb.AtomAtom,	...) */
	root.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, &root)
	root.SetName("mdtwm")

	// Setup list of desk (for now there is only one desk)
	allDesks = NewBoxList()
	currentDesk = NewBox()
	currentDesk.Window = root
	allDesks.PushBack(currentDesk)
	// Setup panels from configured layout
	for _, g := range cfg.Layout {
		currentDesk.Children.PushBack(newPanel(g))
	}
	currentPanel = currentDesk.Children.Front()
	// Initial value of currentWindow
	currentWindow = currentDesk
}

func manageExistingWindows() {
	l.Print("manageExistingWindows")
	tr, err := conn.QueryTree(root.Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		w := Window(id)
		winAdd(w, root)
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
			//l.Print("WaitForEvent error: ", err)
			continue
		}
		switch e := event.(type) {
		case xgb.MapRequestEvent:
			mapRequest(e)
		case xgb.EnterNotifyEvent:
			enterNotify(e)
		case xgb.DestroyNotifyEvent:
			destroyNotify(e)
		case xgb.ConfigureNotifyEvent:
			configureNotify(e)
		case xgb.ConfigureRequestEvent:
			configureRequest(e)
		case xgb.KeyPressEvent:
			keyPress(e)
		case xgb.ButtonPressEvent:
			buttonPress(e)
		default:
			l.Print("Unhandled event: ", reflect.TypeOf(e))

		}
	}
}
