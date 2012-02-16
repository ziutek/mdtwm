package mdtwm

import (
	"github.com/ziutek/mdtwm/xgb_patched"
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
	if uintptr(prop.Format/8) != reflect.TypeOf(xgb.Id(0)).Size() {
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
	var title string
	if currentBox != nil {
		title = currentBox.Name()
	}
	cfg.StatusLogger.Log(Status{currentDeskNum, root.Children().Len(), title})
}

func setCurrentDesk(deskNum int) {
	currentDeskNum = deskNum
	for d := root.Children().Front(); d != nil; d = d.Next() {
		if deskNum == 0 {
			currentDesk = d.(*Panel)
			break
		}
		deskNum--
	}
	updateCurrentDesk()
}

func setNextDesk() {
	nextDesk := currentDesk.Next()
	currentDeskNum++
	if nextDesk == nil {
		nextDesk = root.Children().Front()
		currentDeskNum = 0
	}
	currentDesk = nextDesk.(*Panel)
	updateCurrentDesk()
}

func setPrevDesk() {
	prevDesk := currentDesk.Prev()
	currentDeskNum--
	if prevDesk == nil {
		prevDesk = root.Children().Back()
		currentDeskNum = root.Children().Len() - 1
	}
	currentDesk = prevDesk.(*Panel)
	updateCurrentDesk()
}

func updateCurrentDesk() {
	currentDesk.Raise()
	root.Window().ChangeProp(xgb.PropModeReplace, AtomNetCurrentDesktop,
		xgb.AtomCardinal, uint32(currentDeskNum))
}
