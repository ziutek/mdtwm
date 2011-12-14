package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"x-go-binding.googlecode.com/hg/xgb"
)

func rgbColor(r, g, b uint16) uint32 {
	c, err := conn.AllocColor(screen.DefaultColormap, r, g, b)
	if err != nil {
		l.Fatalf("Cannot allocate a color (%x,%x,%x): %s", r, g, b, err)
	}
	return c.Pixel
}

func namedColor(name string) uint32 {
	c, err := conn.AllocNamedColor(screen.DefaultColormap, name)
	if err != nil {
		l.Fatalf("Cannot allocate a color by name '%s': %s", name, err)
	}
	return c.Pixel
}

type List []interface{}

func (l List) Contains(e interface{}) bool {
	for _, v := range l {
		if v == e {
			return true
		}
	}
	return false
}

type Cmd struct {
	Func  func(string) error
	Param string
}

func (c *Cmd) Run() error {
	return c.Func(c.Param)
}

func spawn(cmd string) error {
	return exec.Command(cmd).Start()
}

// Keycodes
const (
	KeyA = 38
	KeyB = 56
	KeyC = 54
	KeyD = 40
	KeyE = 26
	KeyF = 41
	KeyG = 42
	KeyH = 43
	KeyI = 31
	KeyJ = 44
	KeyK = 45
	KeyL = 46
	KeyM = 58
	KeyN = 57
	KeyO = 32
	KeyP = 33
	KeyQ = 24
	KeyR = 27
	KeyS = 39
	KeyT = 28
	KeyU = 30
	KeyV = 55
	KeyW = 25
	KeyX = 53
	KeyY = 29
	KeyZ = 52

	Key1 = 10
	Key2 = 11
	Key3 = 12
	Key4 = 13
	Key5 = 14
	Key6 = 15
	Key7 = 16
	Key8 = 17
	Key9 = 18
	Key0 = 19

	KeyComma  = 59
	KeyDot    = 60
	KeySpace  = 65
	KeyEnter  = 36
	KeyBspace = 22

	KeyUp    = 111
	KeyLeft  = 113
	KeyRight = 114
	KeyDown  = 116
)

type Config struct {
	Instance string
	Class    string

	BackgroundColor    uint32
	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        int16

	ModMask uint16
	Keys    map[byte]Cmd

	Ignore List
	Float  List
}

var cfg *Config

func loadConfig() {
	cfg = &Config{
		Instance: filepath.Base(os.Args[0]),
		Class:    "Mdtwm",

		BackgroundColor:    namedColor("gray"),
		NormalBorderColor:  rgbColor(0x8888, 0x8888, 0x8888),
		FocusedBorderColor: rgbColor(0xeeee, 0x0000, 0x1111),
		BorderWidth:        1,

		ModMask: xgb.ModMask1,
		Keys: map[byte]Cmd{
			KeyEnter:  {spawn, "xterm"},
			KeyBspace: {spawn, "xkill"},
		},

		Ignore: List{"Unity-2d-panel", "Unity-2d-launcher"},
		Float:  List{"Mplayer", "Gimp"},
	}
}
