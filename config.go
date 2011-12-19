package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"os"
	"path/filepath"
)

type Config struct {
	Instance string
	Class    string

	BackgroundColor    uint32
	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        int16

	DefaultCursor    xgb.Id
	MoveCursor       xgb.Id
	MultiClickTime  xgb.Timestamp
	MovedClickRadius int

	ModMask uint16
	Keys    map[byte]Cmd

	Ignore List
	Float  List
}

func configure() {
	// Configuration variables

	cfg = &Config{
		Instance: filepath.Base(os.Args[0]),
		Class:    "Mdtwm",

		NormalBorderColor:  rgbColor(0x8888, 0x8888, 0x8888),
		FocusedBorderColor: rgbColor(0xeeee, 0x0000, 0x1111),
		BorderWidth:        1,

		DefaultCursor:    stdCursor(68),
		MoveCursor:       stdCursor(52),
		MultiClickTime:   300, // maximum interfal for multiclick [ms]
		MovedClickRadius: 5,  // minimal radius for moved click [pixel]

		ModMask: xgb.ModMask1,
		Keys: map[byte]Cmd{
			KeyEnter:  {spawn, "xterm"},
			KeyBspace: {spawn, "xkill"},
		},

		Ignore: List{"Unity-2d-panel", "Unity-2d-launcher"},
		Float:  List{"Mplayer", "Gimp"},
	}
	// We use square of radius
	cfg.MovedClickRadius *= cfg.MovedClickRadius

	// Initial layout

	root = NewRootPanel()
	// Setup list of desk (for now there is only one desk)
	currentDesk = NewPanel(Horizontal, 1.75)
	root.Insert(currentDesk)
	// Setup two main panels
	currentDesk.Insert(NewPanel(Vertical, 1))
	currentDesk.Insert(NewPanel(Vertical, 0.3))
	// All windows that exists during startup will be placed in currentBox
	currentBox = currentDesk.Children().Front()
}
