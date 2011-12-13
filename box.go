package main

import (
	"fmt"
	"unicode/utf16"
	"x-go-binding.googlecode.com/hg/xgb"
)

const (
	/*boxEventMask = xgb.EventMaskButtonPress |
	xgb.EventMaskButtonRelease |
	//xgb.EventMaskPointerMotion |
	xgb.EventMaskExposure | // window needs to be redrawn
	xgb.EventMaskStructureNotify | // window gets destroyed
	xgb.EventMaskSubstructureRedirect | // app tries to resize itself
	xgb.EventMaskSubstructureNotify | // subwindows get notifies
	xgb.EventMaskEnterWindow |
	xgb.EventMaskPropertyChange |
	xgb.EventMaskFocusChange*/
	boxEventMask = xgb.EventMaskEnterWindow | xgb.EventMaskStructureNotify
)

type Box interface {
	String() string

	Window() Window
	Parent() ParentBox
	SetParent(p ParentBox)
	Children() BoxList

	SetPosSize(x, y, width, height int16) // Set position and EXTERNAL size
	SetFocus(cur bool)

	// Properties
	Name() string
	NameX() []uint16 // UCS2 encoded name
	SetName(name string)
	Class() (instance, class string)
	SetClass(instance, class string)
}

type ParentBox interface {
	Box

	Insert(b Box)
	Remove(b Box)
}

type commonBox struct {
	w        Window     // window stored in this box
	parent   ParentBox // parent panel
	children BoxList    // child boxes contains childs of this box
}

func (b *commonBox) String() string {
	return fmt.Sprintf("%s (%s)", b.Name(), b.w)
}

func (b *commonBox) Window() Window {
	return b.w
}

func (b *commonBox) init(w Window) {
	b.w = w
	b.children = NewBoxList()
}

func (b *commonBox) grabInput(confineTo Window) {
	// Grab right mouse buttons for WM actions
	b.w.GrabButton(false, xgb.EventMaskButtonPress, xgb.GrabModeSync,
		xgb.GrabModeAsync, confineTo, xgb.CursorNone, 3, xgb.ButtonMaskAny)
	for k, _ := range cfg.Keys {
		b.w.GrabKey(true, cfg.ModMask, k, xgb.GrabModeAsync, xgb.GrabModeAsync)
	}
}

func (b *commonBox) Parent() ParentBox {
	return b.parent
}

func (b *commonBox) Children() BoxList {
	return b.children
}

func (b *commonBox) SetParent(p ParentBox) {
	b.parent = p
	b.w.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmapNotify
	b.w.Reparent(p.Window(), 0, 0)
	b.w.SetEventMask(boxEventMask)
}

// Properties

func (b *commonBox) Name() string {
	// We prefer utf8 version
	if p, err := b.w.Prop(AtomNetWmName, 128); err == nil && len(p.Value) > 0 {
		return string(p.Value)
	}
	if p, err := b.w.Prop(xgb.AtomWmName, 128); err == nil && len(p.Value) > 0 {
		return string(p.Value)
	}
	return ""
}

func (b *commonBox) NameX() []uint16 {
	return utf16.Encode([]rune(b.Name()))
}

func (b *commonBox) SetName(name string) {
	b.w.ChangeProp(xgb.PropModeReplace, xgb.AtomWmName, xgb.AtomString, name)
	b.w.ChangeProp(xgb.PropModeReplace, AtomNetWmName, AtomUtf8String, name)
}

func (b *commonBox) Class() (instance, class string) {
	return b.w.Class()
}

func (b *commonBox) SetClass(instance, class string) {
	v := make([]byte, 0, len(instance)+len(class)+2)
	v = append(v, instance...)
	v = append(v, 0)
	v = append(v, class...)
	b.w.ChangeProp(xgb.PropModeReplace, xgb.AtomWmClass, xgb.AtomString, v)
}
