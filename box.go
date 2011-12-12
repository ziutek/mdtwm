package main

import (
	"unicode/utf16"
	"x-go-binding.googlecode.com/hg/xgb"
)

type Box interface {
	Window

	NameX() []uint16 // UCS2 encoded name

	Parent() *PanelBox
	SetParent(p *PanelBox)
	Children() BoxList

	SetPosSize(x, y, width, height int16) // Set position and EXTERNAL size
	SetFocus(cur bool)
}

type commonBox struct {
	Window // window stored in this box

	parent   *PanelBox // parent panel
	children BoxList   // child boxes contains childs of this box
}

func (b *commonBox) init(w Window) {
	b.Window = w
	b.children = NewBoxList()
	// Grab right mouse buttons for WM actions
	b.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, root, xgb.CursorNone, 3, xgb.ButtonMaskAny)
	// Chose events that WM is interested in
	b.SetEventMask(EventMask)
}

func (b *commonBox) NameX() []uint16 {
	return utf16.Encode([]rune(b.Name()))
}

func (b *commonBox) Parent() *PanelBox {
	return b.parent
}


func (b *commonBox) Children() BoxList {
	return b.children
}


func (b *commonBox) SetParent(p *PanelBox) {
	b.parent = p
	b.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmapNotify due to reparenting
	b.Reparent(p, 0, 0)
	b.SetEventMask(EventMask)
}
