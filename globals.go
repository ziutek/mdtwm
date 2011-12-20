package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"log"
	"os"
)

var (
	conn   *xgb.Conn
	screen *xgb.ScreenInfo

	// Desk in mdtwm means workspace. Desk contains panels. Panel contains
	// panels or windows.
	root        *RootPanel
	currentDesk *Panel
	currentBox  Box

	cfg *Config

	l = log.New(os.Stderr, "mdtwm: ", log.Lshortfile )
)
