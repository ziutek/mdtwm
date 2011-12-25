package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	conn *xgb.Conn
	root Window
	l    = log.New(os.Stderr, "test: ", 0)
)

func createWindow() {
	w := NewWindow(
		root, Geometry{0, 0, 100, 30, 0},
		xgb.WindowClassInputOutput,
		xgb.CWOverrideRedirect, 1,
	)
	w.Map()
}

func main() {
	var err error
	display := os.Getenv("DISPLAY")
	conn, err = xgb.Dial(display)
	if err != nil {
		l.Fatal(err)
	}
	setupAtoms()
	root = Window(conn.DefaultScreen().Root)

	createWindow()

	tr, err := conn.QueryTree(root.Id())
	for i, id := range append(tr.Children, root.Id()) {
		w := Window(id)

		inst, class := w.Class()

		tr, err := conn.QueryTree(id)
		if err != nil {
			l.Fatal("QueryTree: ", err)
		}
		wmName := w.Prop(AtomNetWmName, 128)

		info := struct {
			id, root, parent  xgb.Id
			ch_num            uint16
			name, inst, class string
			g                 Geometry
		}{
			id, tr.Root, tr.Parent,
			tr.ChildrenLen,
			string(wmName.Value), inst, class,
			Geometry{},
		}
		var ok bool
		info.g, ok = w.Geometry()
		if !ok {
			return
		}

		fmt.Printf("%d: %+v\n", i, info)
	}
	time.Sleep(100e9)
}
