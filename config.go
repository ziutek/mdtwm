package main

import (
	"os"
	"path/filepath"
	"x-go-binding.googlecode.com/hg/xgb"
)

func allocColor(r, g, b uint16) uint32 {
	c, err := conn.AllocColor(screen.DefaultColormap, r, g, b)
	if err != nil {
		l.Fatalf("Cannot allocate a color (%x,%x,%x): %s", r, g, b, err)
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

type Config struct {
	Instance string
	Class    string

	Layout             []Geometry
	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        int16

	ModMask uint16

	Ignore List
	Float  List
}

var cfg *Config

func loadConfig() {
	l.Print("loadConfig")
	cfg = &Config{
		Instance: filepath.Base(os.Args[0]),
		Class:    "Mdtwm",

		NormalBorderColor:  allocColor(0xaaaa, 0xaaaa, 0xaaaa),
		FocusedBorderColor: allocColor(0xf444, 0x0000, 0x000f),
		BorderWidth:        1,

		ModMask: xgb.ModMask1,

		Ignore: List{"Unity-2d-panel", "Unity-2d-launcher"},
		Float:  List{"Mplayer", "Gimp"},
	}
}
