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

	// Desk in mdtwm means workspace. Desk contains panels. Panel contains
	// panels or windows.
	root         ParentBox
	currentDesk  ParentBox
	currentPanel ParentBox

	l = log.New(os.Stderr, "mdtwm: ", 0)
)

func main() {
	signals()
	connect()
	setupAtoms()
	loadConfig()
	setupWm()
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
}

func setupWm() {
	l.Print("setupWm")
	// Setup root window (RootPanel)
	root = NewRootPanel()
	// Setup list of desk (for now there is only one desk)
	currentDesk = NewPanel(Horizontal)
	root.Insert(currentDesk)
	// Setup two main panels
	// TODO: Use configuration for this
	currentDesk.Insert(NewPanel(Vertical))
	currentDesk.Insert(NewPanel(Vertical))
	// Initial value of currentPanel and currentWindow
	currentPanel = currentDesk.Children().Front().(*Panel)
}

func manageExistingWindows() {
	l.Print("manageExistingWindows")
	tr, err := conn.QueryTree(root.Window().Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		manage(Window(id), currentPanel, true)
	}
}

func eventLoop() {
	l.Print("eventLoop")
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
		case xgb.UnmapNotifyEvent:
			unmapNotify(e)
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
