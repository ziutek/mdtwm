package mdtwm

import (
	"math"
	"runtime"
)

// Utils
func Uint16(x int16) uint16 {
	if x < 0 {
		l.Panic("Can't convert negative int16 to uint16")
	}
	return uint16(x)
}

func Pint16(x int16) uint16 {
	r := Uint16(x)
	if r == 0 {
		l.Print("Pint16(0)")
		return 1
	}
	return r
}

func Int16(x uint16) int16 {
	if x > math.MaxInt16 {
		l.Panicf("Can't convert %d to int16", x)
	}
	return int16(x)
}

func put16(buf []byte, v uint16) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
}

func put32(buf []byte, v uint32) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

// Logs error prefixed by function name
func logFuncErr(e error) {
	fname := "Unknown"
	if pc, _, _, ok := runtime.Caller(1); ok {
		f := runtime.FuncForPC(pc)
		fname = f.Name()
	}
	l.Printf("%s: %s", fname, e)
}
