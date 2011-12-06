package main

import (
	"container/list"
	"x-go-binding.googlecode.com/hg/xgb"
)

type WindowListIterator struct {
	current *list.Element
	back bool
}

func (i *WindowListIterator) Done() bool {
	return i.current == nil
}

func (i *WindowListIterator) Next() Window {
	w := i.current.Value.(Window)
	if i.back {
		i.current = i.current.Prev()
	} else {
		i.current = i.current.Next()
	}
	return w
}

type WindowList struct {
	raw *list.List
}

func NewWindowList() WindowList {
	return WindowList{list.New()}
}

func (wl WindowList) PushFront(w Window) {
	wl.raw.PushFront(w)
}

func (wl WindowList) PushBack(w Window) {
	wl.raw.PushBack(w)
}

func (wl WindowList) Len() int {
	return wl.raw.Len()
}

func (wl WindowList) FrontIter() *WindowListIterator {
	return &WindowListIterator{wl.raw.Front(), false}
}

func (wl WindowList) BackIter() *WindowListIterator {
	return &WindowListIterator{wl.raw.Back(), true}
}

func (wl WindowList) Contains(w Window) bool {
	for e := wl.raw.Front(); e != nil; e = e.Next() {
		if e.Value.(Window) == w {
			return true
		}
	}
	return false
}

func (wl WindowList) Remove(w Window) {
	for e := wl.raw.Front(); e != nil; e = e.Next() {
		if e.Value.(Window) == w {
			wl.raw.Remove(e)
			return
		}
	}
	panic("Can't remove non existent window form list")
}

var (
	windows = NewWindowList()
)


func winAdd(w Window) {
	l.Print("manageWindow: ", w)
	if cfg.Ignore.Contains(w.Class()) {
		return
	}
	wa := w.Attrs()
	if wa.OverrideRedirect || wa.MapState != xgb.MapStateViewable {
		return
	}
	w.EventMask(xgb.EventMaskPropertyChange | xgb.EventMaskStructureNotify |
		xgb.EventMaskFocusChange)
	// Nice bechavior if wm will be killed, exited, crashed
	w.ChangeSaveSet(xgb.SetModeInsert)

	w.SetBorderWidth(cfg.BorderWidth)
	w.SetBorderColor(cfg.NormalBorderColor)
	windows.PushFront(w)
}

func winFocus(w Window) {
	w.SetBorderColor(cfg.FocusedBorderColor)
	w.SetInputFocus()
}
