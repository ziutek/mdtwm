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

func atomList(prop *xgb.GetPropertyReply) IdList {
	if prop == nil || prop.ValueLen == 0 {
		return nil
	}
	if uintptr(prop.Format / 8) != reflect.TypeOf(xgb.Id(0)).Size() {
		l.Panic("Property reply has wrong format for atoms: ", prop.Format)
	}
	return (*[1 << 24]xgb.Id)(unsafe.Pointer(&prop.Value[0]))[:prop.ValueLen]
}

func removeWindow(w Window, unmap bool) {
	if b := root.Children().BoxByWindow(w, true); b != nil {
		b.Parent().Remove(b, unmap)
	}
}
