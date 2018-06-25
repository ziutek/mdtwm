package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ziutek/mdtwm/xgb_patched"
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

func closeCurrentWindow(cmd interface{}) error {
	if currentBox == nil {
		return nil
	}
	if b, ok := currentBox.(*BoxedWindow); ok {
		if b.Protocols().Contains(AtomWmDeleteWindow) {
			b.SendMessage(AtomWmDeleteWindow, b.Window())
		} else {
			b.Window().Destroy()
		}
	}
	return nil
}

func exit(retval interface{}) error {
	os.Exit(retval.(int))
	return nil
}

func chDesk(deskNum interface{}) error {
	setCurrentDesk(deskNum.(int))
	return nil
}

func nextDesk(interface{}) error {
	setNextDesk()
	return nil
}

func prevDesk(interface{}) error {
	setPrevDesk()
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

	KeyMinus        = 0x002d
	KeyEqual        = 0x003d
	KeySemicolon    = 0x003b
	KeyComma        = 0x002c
	KeyPeriod       = 0x002e
	KeySlash        = 0x002f
	KeyApostrophe   = 0x0027
	KeyBackslash    = 0x005c
	KeyBracketRight = 0x005d
	KeyBracketLeft  = 0x005b

	KeySpace      = 0x0020
	KeyBackSpace  = 0xff08
	KeyTab        = 0xff09
	KeyReturn     = 0xff0d
	KeyPause      = 0xff13
	KeyScrollLock = 0xff14
	KeySysReq     = 0xff15
	KeyEscape     = 0xff1b
	KeyInsert     = 0xff63 // ?
	KeyDelete     = 0xffff

	KeyHome     = 0xff50
	KeyPageDown = 0xff56
	KeyPageUp   = 0xff55
	KeyEnd      = 0xff57
	KeyLeft     = 0xff51
	KeyUp       = 0xff52
	KeyRight    = 0xff53
	KeyDown     = 0xff54

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

type sysStat struct {
	act, sum uint64
}

type Dzen2Logger struct {
	io.Writer
	FgColor    string
	BgColor    string
	BatPath    string
	TimeFormat string
	InfoPos    int16

	ch chan *Status

	i    [9]byte
	c    [4]byte
	s    [1]byte
	stat []sysStat
}

func (d *Dzen2Logger) invColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.FgColor, d.BgColor)
}
func (d *Dzen2Logger) nrmColors() {
	fmt.Fprintf(d.Writer, "^bg(%s)^fg(%s)", d.BgColor, d.FgColor)
}

func (d *Dzen2Logger) Start() {
	if d.InfoPos < 0 {
		_, _, width, _ := root.PosSize()
		d.InfoPos += width
	}
	d.ch = make(chan *Status)
	if f, err := os.Open("/proc/stat"); err == nil {
		sc := bufio.NewScanner(f)
		n := 0
		for sc.Scan() {
			line := sc.Text()
			if !strings.HasPrefix(line, "cpu") {
				break
			}
			n++
		}
		f.Close()
		if sc.Err() == nil {
			d.stat = make([]sysStat, n-1)
		}
	}
	go d.thr()
}

func (d *Dzen2Logger) Log(s Status) {
	d.ch <- &s
}

func readFull(fname string, buf []byte) (int, error) {
	f, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.ReadFull(f, buf)
}

func (d *Dzen2Logger) batInfo() string {
	if d.BatPath == "" {
		return ""
	}

	c := d.c[:]
	n, err := readFull(filepath.Join(d.BatPath, "capacity"), c)
	if err != nil && err != io.ErrUnexpectedEOF {
		return ""
	}
	if n > 0 {
		c = c[:n-1]
	} else {
		c = nil
	}

	i := d.i[:]
	i[0] = ' '
	n, err = readFull(filepath.Join(d.BatPath, "current_now"), i[1:])
	if err != nil && err != io.ErrUnexpectedEOF {
		return ""
	}
	if n > 0 {
		if n > 3 {
			i = i[:n-3]
		} else {
			i = i[:2]
			i[1] = '0'
		}
	} else {
		i = i[:1]
	}

	_, err = readFull(filepath.Join(d.BatPath, "status"), d.s[:])
	if err == nil && d.s[0] == 'C' {
		i[0] = '~'
	}
	return fmt.Sprintf("[%s%%%5smA]", c, i)
}

func (d *Dzen2Logger) cpuLoad() string {
	if d.stat == nil {
		return ""
	}
	f, err := os.Open("/proc/stat")
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	load := make([]string, len(d.stat))
	i := 0
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "cpu") {
			break
		}
		if strings.HasPrefix(line, "cpu ") {
			continue
		}
		a := strings.Fields(line)
		user, _ := strconv.ParseUint(a[1], 10, 64)
		nice, _ := strconv.ParseUint(a[2], 10, 64)
		system, _ := strconv.ParseUint(a[3], 10, 64)
		idle, _ := strconv.ParseUint(a[4], 10, 64)
		iowait, _ := strconv.ParseUint(a[5], 10, 64)
		irq, _ := strconv.ParseUint(a[6], 10, 64)
		soft, _ := strconv.ParseUint(a[7], 10, 64)
		act := user + nice + system + iowait + irq + soft
		sum := act + idle
		dact := act - d.stat[i].act + 1
		dsum := sum - d.stat[i].sum + 1
		load[i] = fmt.Sprintf("%3d", (dact*100+dsum/2)/dsum)
		d.stat[i].act = act
		d.stat[i].sum = sum
		i++
	}
	return "cpu:" + strings.Join(load, "")
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
			d.Writer, "   %s^pa(%d)%s   %13s   %s\n",
			s.title, d.InfoPos, d.cpuLoad(), d.batInfo(),
			t.Format(d.TimeFormat),
		)
	}
}

func (c *Config) Load(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		if e, ok := err.(*os.PathError); !ok || e.Err != syscall.ENOENT {
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
