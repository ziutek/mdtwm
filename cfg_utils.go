package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	Func  func(interface{}) error
	Param interface{}
}

func (c *Cmd) Run() error {
	return c.Func(c.Param)
}

func spawn(cmd interface{}) error {
	// TODO: check what filedescriptors are inherited from WM by cmd when
	// exec.Command is used
	return exec.Command(fmt.Sprint(cmd)).Start()
}

func exit(retval interface{}) error {
	os.Exit(retval.(int))
	return nil
}

func chDesk(deskNum interface{}) error {
	i := root.Children().FrontIter(false)
	n := deskNum.(int)
	for d := i.Next(); d != nil; d = i.Next() {
		if n--; n == 0 {
			currentDesk = d.(*Panel)
			d.Raise()
			break
		}
	}
	return nil
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
	TimePos    int

	ch chan *Status
}

func (d *Dzen2Logger) invColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.FgColor, d.BgColor)
}
func (d *Dzen2Logger) nrmColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.BgColor, d.FgColor)
}

func (d *Dzen2Logger) Start() {
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
			s.curDesk++
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
