package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"log"
	"os"
)

var (
	conn *xgb.Conn
	l    = log.New(os.Stderr, "test: ", 0)
)

func main() {
	var err error
	display := os.Getenv("DISPLAY")
	conn, err = xgb.Dial(display)
	if err != nil {
		l.Fatal(err)
	}
	setupAtoms()

	root := Window(conn.DefaultScreen().Root)
	tr, err := conn.QueryTree(root.Id())
	for i, id := range append(tr.Children, root.Id()) {
		w := Window(id)

		inst, class := w.Class()

		tr, err := conn.QueryTree(id)
		if err != nil {
			l.Fatal("QueryTree: ", err)
		}
		wmName, _ := w.Prop(AtomNetWmName, 128)

		info := struct {
			id, root, parent  xgb.Id
			ch_num            uint16
			name, inst, class string
			g                 Geometry
		}{
			id, tr.Root, tr.Parent,
			tr.ChildrenLen,
			string(wmName.Value), inst, class,
			w.Geometry(),
		}

		fmt.Printf("%d: %+v\n", i, info)
	}
}
