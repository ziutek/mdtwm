package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"log"
	"os"
)

var (
	conn   *xgb.Conn
	screen *xgb.ScreenInfo
	keyCodeToSym []xgb.Keysym
	keySymToCode map[xgb.Keysym]byte

	// Desk in mdtwm means workspace. Desk contains panels. Panel contains
	// panels or windows.
	root        *RootPanel
	currentDesk *Panel
	currentBox  Box

	currentDeskNum int

	cfg *Config

	l = log.New(os.Stderr, "mdtwm: ", log.Lshortfile)
	d = log.New(os.Stderr, "mdtwm debug: ", log.Lshortfile)
)
