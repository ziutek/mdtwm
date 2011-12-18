package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"math"
	"os/exec"
	"reflect"
	"unsafe"
)

func currentPanel() ParentBox {
	if p, ok := currentBox.(ParentBox); ok {
		return p
	}
	return currentBox.Parent()
}

func changeFocusTo(w Window) {
	currentDesk.SetFocus(currentDesk.Window() == w)
	// Iterate over all boxes in current desk
	bi := currentDesk.Children().FrontIter(true)
	for b := bi.Next(); b != nil; b = bi.Next() {
		b.SetFocus(b.Window() == w)
	}
}

type IdList []xgb.Id

func (l IdList) Contains(id xgb.Id) bool {
	for _, i := range l {
		if i == id {
			return true
		}
	}
	return false
}

func propReplyAtoms(prop *xgb.GetPropertyReply) IdList {
	if prop == nil || prop.ValueLen == 0 {
		return nil
	}
	atom_size := uintptr(prop.Format / 8)
	if atom_size != reflect.TypeOf(xgb.Id(0)).Size() {
		panic("Property reply has wrong format for atoms")
	}
	num_atoms := prop.ValueLen / uint32(atom_size)
	return (*[1 << 24]xgb.Id)(unsafe.Pointer(&prop.Value[0]))[:num_atoms]
}

func Uint16(x int16) uint16 {
	if x < 0 {
		panic("Can't convert negative int16 to uint16")
	}
	return uint16(x)
}

func Pint16(x int16) uint16 {
	r := Uint16(x)
	if r == 0 {
		l.Print("Warn: Pint16(0)")
		return 1
	}
	return r
}

func Int16(x uint16) int16 {
	if x > math.MaxInt16 {
		panic("Can't convert big uint16 to int16")
	}
	return int16(x)
}

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
	// TODO: check what filedescriptors are inherited from WM by cmd when
	// exec.Command is used
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
