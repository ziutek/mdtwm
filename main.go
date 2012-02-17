package main

import (
	"github.com/ziutek/mdtwm/xgb_patched"
	"os"
	"os/signal"
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
	sc := make(chan os.Signal, 3)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGCHLD)
	go func() {
		for sig := range sc {
			if s, ok := sig.(syscall.Signal); ok {
				switch s {
				case syscall.SIGTERM, syscall.SIGINT:
					os.Exit(0)
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
	r, err := conn.GetKeyboardMapping(
		conn.Setup.MinKeycode,
		conn.Setup.MaxKeycode-conn.Setup.MinKeycode+1,
	)
	if err != nil {
		l.Fatal("Can't get keybboard mapping: ", err)
	}
	// Setup keycode <-> keysym mapping
	minKeycode := conn.Setup.MinKeycode
	codeNum := int(r.Length) / int(r.KeysymsPerKeycode)
	keyCodeToSym = make([]xgb.Keysym, int(minKeycode)+codeNum)
	keySymToCode = make(map[xgb.Keysym]byte)
	for i := 0; i < codeNum; i++ {
		s := r.Keysyms[i*int(r.KeysymsPerKeycode)]
		keyCodeToSym[byte(i)+minKeycode] = s
		keySymToCode[s] = byte(i) + minKeycode
	}
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
		handleEvent(conn.WaitForEvent())
	}
}
