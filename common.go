package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"reflect"
	"unsafe"
)

func currentPanel() ParentBox {
	if p, ok := currentBox.(ParentBox); ok {
		return p
	}
	return currentBox.Parent()
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

func atomList(prop *xgb.GetPropertyReply) IdList {
	if prop == nil || prop.ValueLen == 0 {
		return nil
	}
	if uintptr(prop.Format / 8) != reflect.TypeOf(xgb.Id(0)).Size() {
		l.Panic("Property reply has wrong format for atoms: ", prop.Format)
	}
	return (*[1 << 24]xgb.Id)(unsafe.Pointer(&prop.Value[0]))[:prop.ValueLen]
}

func prop32(prop *xgb.GetPropertyReply) []uint32 {
	if prop == nil || prop.ValueLen == 0 {
		return nil
	}
	if prop.Format != 32 {
		l.Panic("Property reply contains %d-bit values (need 32-bit).",
			prop.Format)
	}
	return (*[1 << 24]uint32)(unsafe.Pointer(&prop.Value[0]))[:prop.ValueLen]
}

func statusLog() {
	if cfg.StatusLogger == nil {
		return
	}
	var cur, n int
	for p := root.Children().Front(); p != nil; p = p.Next() {
		if p == currentDesk {
			cur = n
			break
		}
		n++
	}
	cfg.StatusLogger.Log(Status{cur, root.Children().Len(), currentBox.Name()})
}
