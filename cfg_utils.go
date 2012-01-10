package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
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

type TextList []string

func (l TextList) Contains(e string) bool {
	for _, v := range l {
		if v == e {
			return true
		}
	}
	return false
}

type Cmd struct {
	Func  func(interface{}) error
	Param interface{}
}

func (c *Cmd) Run() error {
	return c.Func(c.Param)
}

func spawn(cmd interface{}) error {
	// TODO: check what filedescriptors are inherited from WM by cmd when
	// exec.Command is used
	args := strings.Split(cmd.(string), " ")
	return exec.Command(args[0], args[1:]...).Start()
}

func exit(retval interface{}) error {
	os.Exit(retval.(int))
	return nil
}

func chDesk(deskNum interface{}) error {
	setCurrentDesk(deskNum.(int))
	return nil
}

// Keycodes
const (
	KeyA = 0x0061
	KeyB = 0x0062
	KeyC = 0x0063
	KeyD = 0x0064
	KeyE = 0x0065
	KeyF = 0x0066
	KeyG = 0x0067
	KeyH = 0x0068
	KeyI = 0x0069
	KeyJ = 0x006a
	KeyK = 0x006b
	KeyL = 0x006c
	KeyM = 0x006d
	KeyN = 0x006e
	KeyO = 0x006f
	KeyP = 0x0070
	KeyQ = 0x0071
	KeyR = 0x0072
	KeyS = 0x0073
	KeyT = 0x0074
	KeyU = 0x0075
	KeyV = 0x0076
	KeyW = 0x0077
	KeyX = 0x0078
	KeyY = 0x0079
	KeyZ = 0x007a

	Key0 = 0x0030
	Key1 = 0x0031
	Key2 = 0x0032
	Key3 = 0x0033
	Key4 = 0x0034
	Key5 = 0x0035
	Key6 = 0x0036
	Key7 = 0x0037
	Key8 = 0x0038
	Key9 = 0x0039

	KeyBackSpace  = 0xff08
	KeyTab        = 0xff09
	KeyReturn     = 0xff0d
	KeyPause      = 0xff13
	KeyScrollLock = 0xff14
	KeySysReq     = 0xff15
	KeyEscape     = 0xff1b
	KeyDelete     = 0xffff

	KeyF1  = 0xffbe
	KeyF2  = 0xffbf
	KeyF3  = 0xffc0
	KeyF4  = 0xffc1
	KeyF5  = 0xffc2
	KeyF6  = 0xffc3
	KeyF7  = 0xffc4
	KeyF8  = 0xffc5
	KeyF9  = 0xffc6
	KeyF10 = 0xffc7
	KeyF11 = 0xffc8
	KeyF12 = 0xffc9
)

var stdCursorFont xgb.Id

func stdCursor(id uint16) xgb.Id {
	if stdCursorFont == 0 {
		stdCursorFont = conn.NewId()
		conn.OpenFont(stdCursorFont, "cursor")
	}
	cursor := conn.NewId()
	conn.CreateGlyphCursor(cursor, stdCursorFont, stdCursorFont, id, id+1,
		0, 0, 0, 0xffff, 0xffff, 0xffff)
	return cursor
}

type Status struct {
	curDesk, desks int
	title          string
}

type StatusLogger interface {
	Log(s Status)
	Start()
}

type Dzen2Logger struct {
	io.Writer
	FgColor    string
	BgColor    string
	TimeFormat string
	TimePos    int16

	ch chan *Status
}

func (d *Dzen2Logger) invColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.FgColor, d.BgColor)
}
func (d *Dzen2Logger) nrmColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.BgColor, d.FgColor)
}

func (d *Dzen2Logger) Start() {
	if d.TimePos < 0 {
		_, _, width, _ := root.PosSize()
		d.TimePos += width
	}
	d.ch = make(chan *Status)
	go d.thr()
}

func (d *Dzen2Logger) Log(s Status) {
	d.ch <- &s
}

func (d *Dzen2Logger) thr() {
	var s *Status
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
		case s = <-d.ch:
			s.curDesk++ // Printed desk names starts from 1
		}
		if s == nil {
			continue
		}
		d.nrmColors()
		for i := 1; i <= s.desks; i++ {
			if i == s.curDesk {
				d.invColors()
			}
			fmt.Fprintf(d.Writer, " %d ", i)
			if i == s.curDesk {
				d.nrmColors()
			}
		}
		t := time.Now()
		fmt.Fprintf(
			d.Writer, "   %s^pa(%d)%s\n",
			s.title, d.TimePos, t.Format(d.TimeFormat),
		)
	}
}

func (c *Config) Load(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		if e, ok := err.(*os.PathError); !ok || e.Err != os.ENOENT {
			l.Fatalf("Can't open a configuration file: %s", err)
		}
		// Configuration file doesn't exists: create default
		if f, err = os.Create(fname); err != nil {
			l.Fatal("Can't create a configuration file: ", err)
		}
		buf, err := json.MarshalIndent(c, "", "\t")
		if err != nil {
			l.Fatal("Can't encode a configuration: ", err)
		}
		if _, err = f.Write(buf); err != nil {
			l.Fatal("Can't write a configuration file: ", err)
		}
	} else {
		if err = json.NewDecoder(f).Decode(c); err != nil {
			l.Fatal("Can't decode a configuration file: ", err)
		}
	}

}
