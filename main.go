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
	// panels or windows.
	allDesks      BoxList
	currentDesk   *PanelBox
	currentPanel  *PanelBox

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
	root = RawWindow(screen.Root)
}

func setupWm() {
	l.Print("setupWm")
	// Spported atoms
	/* root.ChangeProp(xgb.PropModeReplace, AtomNetSupported,
	xgb.AtomAtom,	...) */
	root.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, root)
	root.SetClass("mdtwm", "Mdtwm")
	root.SetName("mdtwm root")
	// Setup list of desk (for now there is only one desk)
	allDesks = NewBoxList()
	currentDesk = DeskPanelBox(Horizontal)
	currentDesk.Map()
	allDesks.PushBack(currentDesk)
	// Setup two main panels
	// TODO: Use configuration for this
	currentDesk.Insert(NewPanelBox(Vertical))
	currentDesk.Insert(NewPanelBox(Vertical))
	// Initial value of currentPanel and currentWindow
	currentPanel = currentDesk.Children().Front().(*PanelBox)
}

func manageExistingWindows() {
	l.Print("manageExistingWindows")
	tr, err := conn.QueryTree(root.Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		manage(RawWindow(id), currentPanel)
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
		/*xgb.EventMaskSubstructureRedirect |
			xgb.EventMaskStructureNotify |
			//xgb.EventMaskPointerMotion |
			xgb.EventMaskPropertyChange |
			xgb.EventMaskEnterWindow,*/
		xgb.EventMaskSubstructureRedirect | xgb.EventMaskEnterWindow,
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
