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
	MultiClickTime   xgb.Timestamp
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
		MovedClickRadius: 5,   // minimal radius for moved click [pixel]

		ModMask: xgb.ModMask4,
		Keys: map[byte]Cmd{
			Key1:     {chDesk, 1},
			Key2:     {chDesk, 2},
			KeyEnter: {spawn, "gnome-terminal"},
			KeyQ:     {exit, 0},
		},

		Ignore: List{"Unity-2d-panel", "Unity-2d-launcher"},
		Float:  List{"Mplayer", "Gimp"},
	}
	// We use square of radius
	cfg.MovedClickRadius *= cfg.MovedClickRadius

	// Layout
	root = NewRootPanel()
	// Setup list of desk
	desk1 := NewPanel(Horizontal, 1.75)
	desk2 := NewPanel(Horizontal, 1)
	root.Insert(desk1)
	root.Insert(desk2)
	// Setup two main panels on first desk
	desk1.Insert(NewPanel(Vertical, 1))
	desk1.Insert(NewPanel(Vertical, 0.3))
	// Setup one main panel on second desk
	desk2.Insert(NewPanel(Horizontal, 1))
	// Set current desk and current box
	currentDesk = desk1
	currentDesk.Raise()
	// In this box all existing windows will be placed
	currentBox = currentDesk.Children().Front()
}
