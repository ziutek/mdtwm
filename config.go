package main

func allocColor(r, g, b uint16) uint32 {
	c, err := conn.AllocColor(screen.DefaultColormap, r, g, b)
	if err != nil {
		l.Fatal("Cannot allocate a color: ", err)
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
	NormalBorderColor  uint32
	FocusedBorderColor uint32
	BorderWidth        uint32
	Ignore             List
}

var cfg *Config

func loadConfig() {
	cfg = &Config{
		NormalBorderColor:  allocColor(0xaaaa, 0xaaaa, 0xaaaa),
		FocusedBorderColor: allocColor(0x4444, 0x0000, 0xffff),
		BorderWidth:        1,
		Ignore:             List{"Unity-2d-panel", "unity-2d-launcher"},
	}
}
