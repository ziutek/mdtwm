package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"unicode/utf16"
)

const (
	/*boxEventMask = xgb.EventMaskButtonPress |
	xgb.EventMaskButtonRelease |
	xgb.EventMaskPointerMotion |
	xgb.EventMaskExposure | // window needs to be redrawn
	xgb.EventMaskStructureNotify | // Any change in window configuration.
	xgb.EventMaskSubstructureRedirect | // Redirect reconfiguration of children
	xgb.EventMaskSubstructureNotify | //Notify about reconfiguration of children
	xgb.EventMaskEnterWindow |
	xgb.EventMaskPropertyChange |
	xgb.EventMaskFocusChange*/
)

type Box interface {
	String() string

	Window() Window
	Parent() ParentBox
	SetParent(p ParentBox)
	Children() BoxList

	Geometry() Geometry // Cached geometry
	PosSize() (x, y, width, height int16) // Cached position and EXTERNAL size
	SetPosSize(x, y, width, height int16) // Set position and EXTERNAL size
	SetFocus(f bool)
	Raise()

	Float() bool
	SetFloat(float bool)

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
	Remove(b Box, unmap bool)
}

type commonBox struct {
	w        Window    // window stored in this box
	parent   ParentBox // parent panel
	children BoxList   // child boxes contains childs of this box
	eventMask uint32

	float bool

	// Box configuration
	x, y, width, height int16
}

func (b *commonBox) String() string {
	return fmt.Sprintf("%s (%s)", b.Name(), b.w)
}

func (b *commonBox) Window() Window {
	return b.w
}

func (b *commonBox) init(w Window, eventMask uint32) {
	b.w = w
	b.parent = root
	b.children = NewBoxList()
	b.eventMask = eventMask
	b.w.SetEventMask(eventMask)
}

func (b *commonBox) Parent() ParentBox {
	return b.parent
}

func (b *commonBox) SetParent(p ParentBox) {
	// Translate current coordinates to new parent coordinates (useful when new
	// parent is root and window should stay in place.
	b.x, b.y, _, _ = p.Window().TranslateCoordinates(
		b.parent.Window(), b.x, b.y,
	)
	b.parent = p
	b.w.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmpNotify
	b.w.Reparent(p.Window(), b.x, b.y)
	b.w.SetEventMask(b.eventMask)
}

func (b *commonBox) Children() BoxList {
	return b.children
}

func (b *commonBox) Float() bool {
	return b.float
}

func (b *commonBox) SetFloat(float bool) {
	b.float = float
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

func (b *commonBox) PosSize() (x, y, width, height int16) {
	return b.x, b.y, b.width, b.height
}

func (b *commonBox) Raise() {
	
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
