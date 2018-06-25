package main

import (
	"github.com/ziutek/mdtwm/xgb_patched"
	"os"
	"path/filepath"
)

type Config struct {
	Instance string
	Class    string

	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        int16
	StatusLogger       StatusLogger

	DefaultCursor     xgb.Id
	MoveCursor        xgb.Id
	MultiClickTime    xgb.Timestamp
	MovedClickRadius  int16
	ResizeBorderWidth int16

	ModMask uint16
	Keys    map[xgb.Keysym]Cmd

	Ignore TextList
	Float  TextList
}

func configure() {
	// Default configuration
	cfg = &Config{
		Instance: filepath.Base(os.Args[0]),
		Class:    "MDtwm",

		NormalBorderColor: rgbColor(0x8888, 0x8888, 0x8888),
		//FocusedBorderColor: rgbColor(0x4444, 0x0000, 0xffff),
		FocusedBorderColor: rgbColor(0xffff, 0x9999, 0x0000),
		BorderWidth:        1,

		StatusLogger: &Dzen2Logger{
			Writer:     os.Stdout,
			FgColor:    "#ddddcc",
			BgColor:    "#555588",
			BatPath:    "/sys/class/power_supply/BAT1",
			TimeFormat: "Mon, Jan _2 15:04:05",
			InfoPos:    -332, // Negatife value means pixels from right border
		},

		DefaultCursor:     stdCursor(68),
		MoveCursor:        stdCursor(52),
		MultiClickTime:    300, // maximum interval for multiclick [ms]
		MovedClickRadius:  5,   // minimal radius for moved click [pixel]
		ResizeBorderWidth: 6,   // width of imaginary resize border [pixel]

		ModMask: xgb.ModMask4,
		Keys: map[xgb.Keysym]Cmd{
			Key1:      {chDesk, 0},
			Key2:      {chDesk, 1},
			Key3:      {chDesk, 2},
			KeyLeft:   {prevDesk, nil},
			KeyRight:  {nextDesk, nil},
			KeyReturn: {spawn, "x-terminal-emulator"},
			KeyEscape: {exit, 0},
			KeyA:      {spawn, "x-text-editor"},
			KeyC:      {closeCurrentWindow, nil},
			KeyE:      {spawn, "x-email-client"},
			KeyW:      {spawn, "x-www-browser"},
			KeyF11:    {spawn, "sudo systemctl -i suspend"},
			KeyF12:    {spawn, "sudo systemctl -i hibernate"},
		},

		Ignore: TextList{},
		Float:  TextList{"MPlayer", "QEMU", "mpv"},
	}
	// Read configuration from file
	//cfg.Load(filepath.Join(os.Getenv("HOME"), ".mdtwm"))

	// Layout
	root = NewRootPanel()
	// Setup all desks
	desk0 := NewPanel(Horizontal, 484, 1)
	desk1 := NewPanel(Horizontal, 484, 1)
	desk2 := NewPanel(Horizontal, 0, 1)
	root.Append(desk0)
	root.Append(desk1)
	root.Append(desk2)
	// Setup two main vertical panels on first desk
	left := NewPanel(Vertical, 0, 1)
	right := NewPanel(Vertical, 0, 0.3)
	desk0.Append(left)
	desk0.Append(right)
	// Divide right panel into two horizontal panels
	right.Append(NewPanel(Horizontal, 484, 1))
	right.Append(NewPanel(Horizontal, 0, 1))
	// Setup two panels on second desk
	desk1.Append(NewPanel(Vertical, 0, 1.03))
	desk1.Append(NewPanel(Vertical, 0, 0.3))
	// Setup one horizontal panel on thrid desk
	desk2.Append(NewPanel(Horizontal, 0, 1))
	// Set current desk and current box
	setCurrentDesk(0)
	// In this box all existing windows will be placed
	currentBox = currentDesk.Children().Front()

	// Some operation on configuration
	cfg.MovedClickRadius *= cfg.MovedClickRadius // We need square of radius
	if cfg.StatusLogger != nil {
		cfg.StatusLogger.Start()
	}
}
