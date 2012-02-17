package main

import (
	"github.com/ziutek/mdtwm/xgb_patched"
	"fmt"
	"unicode/utf16"
)

type Box interface {
	String() string

	Window() Window
	Parent() ParentBox
	SetParent(p ParentBox)
	Children() *BoxList

	// Methods for implement an element of BoxList
	Next() Box
	SetNext(next Box)
	Prev() Box
	SetPrev(prev Box)
	List() *BoxList
	SetList(l *BoxList)

	Geometry() Geometry // Get internal geometry

	PosSize() (x, y, width, height int16) // Get externel geometry
	SetPosSize(x, y, width, height int16) // Set external geometry

	SetFocus(f bool, t xgb.Timestamp)
	Raise()

	Float() bool
	SetFloat(float bool)
	Hints() Hints
	SetHints(h Hints)

	// Properties
	Name() string
	NameX() []uint16 // UCS2 encoded name
	SetName(name string)
	Class() (instance, class string)
	SetClass(instance, class string)
}

type ParentBox interface {
	Box

	Append(b Box)
	InsertNextTo(b, mark Box, x, y int16)
	Remove(b Box)
}

type Hints struct {
	W, H, MinW, MinH, MaxW, MaxH, IncW, IncH, BaseW, BaseH int16
	MinAspect, MaxAspect [2]int16
	Gravity byte
}

type commonBox struct {
	w        Window    // window stored in this box
	parent   ParentBox // parent panel
	children *BoxList  // child boxes contains childs of this box

	prev, next Box
	list       *BoxList

	eventMask uint32
	float     bool
	hints     Hints

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
	var ok bool
	b.x, b.y, _, _, ok = p.Window().TranslateCoordinates(
		b.parent.Window(), b.x, b.y,
	)
	if !ok {
		return
	}
	b.parent = p
	b.w.SetEventMask(xgb.EventMaskNoEvent) // avoid UnmpNotify
	b.w.Reparent(p.Window(), b.x, b.y)
	b.w.SetEventMask(b.eventMask)
}

func (b *commonBox) Children() *BoxList {
	return b.children
}

func (b *commonBox) Prev() Box {
	return b.prev
}

func (b *commonBox) SetPrev(prev Box) {
	b.prev = prev
}

func (b *commonBox) Next() Box {
	return b.next
}

func (b *commonBox) SetNext(next Box) {
	b.next = next
}

func (b *commonBox) List() *BoxList {
	return b.list
}

func (b *commonBox) SetList(l *BoxList) {
	b.list = l
}

func (b *commonBox) Float() bool {
	return b.float
}

func (b *commonBox) SetFloat(float bool) {
	b.float = float
}

func (b *commonBox) Hints() Hints {
	return b.hints
}

func (b *commonBox) SetHints(h Hints) {
	b.hints = h
}

// Properties

func (b *commonBox) Name() string {
	// We prefer utf8 version
	if p := b.w.Prop(AtomNetWmName, 128); p != nil && len(p.Value) > 0 {
		return string(p.Value)
	}
	if p := b.w.Prop(xgb.AtomWmName, 128); p != nil && len(p.Value) > 0 {
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
	b.Window().Configure(xgb.ConfigWindowStackMode, uint32(xgb.StackModeAbove))
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
