package main

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
	Layout             []Geometry
	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        int16
	Ignore             List
	Float              List
}

var cfg *Config

func loadConfig() {
	l.Print("loadConfig")
	cfg = &Config{
		Layout: []Geometry{{440, 0, 1000, 1080}, {0, 0, 440, 1080}},

		NormalBorderColor:  allocColor(0xaaaa, 0xaaaa, 0xaaaa),
		FocusedBorderColor: allocColor(0xf444, 0x0000, 0x000f),
		BorderWidth:        2,

		Ignore: List{"Unity-2d-panel", "Unity-2d-launcher"},
		Float:  List{"Mplayer", "Gimp"},
	}
}
