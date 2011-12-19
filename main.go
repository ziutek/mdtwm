package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"io"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	signals()
	connect()
	setupAtoms()
	configure()
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
				case syscall.SIGWINCH:
					continue
				case syscall.SIGCHLD:
					var status syscall.WaitStatus
					_, err := syscall.Wait4(-1, &status, 0, nil)
					if err != nil {
						l.Print("syscal.Wait4: ", err)
					}
					continue
				}
			}
			l.Printf("Signal '%s' received and ignored", sig)
		}
	}()
}

func connect() {
	var err error
	conn, err = xgb.Dial("")
	if err != nil {
		l.Fatal("Can't connect to display: ", err)
	}
	screen = conn.DefaultScreen()
}

func manageExistingWindows() {
	tr, err := conn.QueryTree(root.Window().Id())
	if err != nil {
		l.Fatal("Can't get a list of existing windows: ", err)
	}
	for _, id := range tr.Children {
		manage(Window(id), currentPanel(), true)
	}
}

func eventLoop() {
	for {
		event, err := conn.WaitForEvent()
		if err != nil {
			if err == io.EOF {
				conn.Close()
				os.Exit(0)
			}
			l.Print("WaitForEvent error: ", err)
			continue
		}
		switch e := event.(type) {
		// *Request events
		case xgb.MapRequestEvent:
			mapRequest(e)
		case xgb.ConfigureRequestEvent:
			configureRequest(e)
		// *Notify events
		case xgb.EnterNotifyEvent:
			enterNotify(e)
		case xgb.UnmapNotifyEvent:
			unmapNotify(e)
		case xgb.DestroyNotifyEvent:
			destroyNotify(e)
		// Keyboard and mouse events
		case xgb.KeyPressEvent:
			keyPress(e)
		case xgb.ButtonPressEvent:
			buttonPress(e)
		case xgb.ButtonReleaseEvent:
			buttonRelease(e)
		case xgb.MotionNotifyEvent:
			motionNotify(e)
		default:
			l.Print("*** Unhandled event: ", reflect.TypeOf(e))

		}
	}
}
