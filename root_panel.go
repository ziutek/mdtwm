package mdtwm

import (
	"github.com/ziutek/mdtwm/xgb_patched"
)

const rightButtonEventMask = xgb.EventMaskButtonPress |
	xgb.EventMaskButtonRelease | xgb.EventMaskPointerMotion

// Box for root window
type RootPanel struct {
	commonBox

	checkWindow Window
}

func NewRootPanel() *RootPanel {
	var p RootPanel
	p.init(
		Window(screen.Root),
		xgb.EventMaskSubstructureNotify| // For override-redirect windows
			xgb.EventMaskSubstructureRedirect,
	)
	p.width = int16(screen.WidthInPixels)
	p.height = int16(screen.HeightInPixels)
	p.w.ChangeAttrs(xgb.CWCursor, uint32(cfg.DefaultCursor))
	p.SetClass(cfg.Instance, cfg.Class)
	p.SetName("mdtwm root")
	// Create infrastructure for check of existence of active WM
	p.checkWindow = NewWindow(p.Window(), Geometry{0, 0, 1, 1, 0},
		xgb.WindowClassInputOutput, 0)
	p.checkWindow.ChangeProp(xgb.PropModeReplace, AtomNetWmName,
		AtomUtf8String, p.Name())
	p.checkWindow.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, p.checkWindow)
	p.w.ChangeProp(xgb.PropModeReplace, AtomNetSupportingWmCheck,
		xgb.AtomWindow, p.checkWindow)
	// Supported WM properties
	p.w.ChangeProp(
		xgb.PropModeReplace, AtomNetSupported, xgb.AtomAtom,
		[]xgb.Id{
			AtomNetSupportingWmCheck,
			AtomNetWmName,
			AtomNetWmDesktop,
			AtomNetNumberOfDesktops,
			AtomNetCurrentDesktop,
			AtomNetWmStateModal,
			AtomNetWmStateHidden,
			AtomNetActiveWindow,
			AtomNetWmStrut,
		},
	)
	p.w.DeleteProp(AtomNetVirtualRoots) // clear for future append
	// Grab right mouse buttons for WM actions
	p.w.GrabButton(
		true, // Needed for EnterNotify events during grab
		rightButtonEventMask,
		xgb.GrabModeAsync, xgb.GrabModeAsync,
		xgb.WindowNone, cfg.DefaultCursor, 3,
		xgb.ButtonMaskAny,
	)
	// Grab keys for WM actions
	for k, _ := range cfg.Keys {
		p.w.GrabKey(true, cfg.ModMask, keySymToCode[k], xgb.GrabModeAsync,
			xgb.GrabModeAsync)
	}
	return &p
}

func (p *RootPanel) Geometry() Geometry {
	return Geometry{
		X: p.x, Y: p.y,
		W: p.width, H: p.height,
	}
}

func (p *RootPanel) SetPosSize(x, y, width, height int16) {
	l.Panic("Can't request a geometry change of root window")
}

func (p *RootPanel) SetFocus(f bool, t xgb.Timestamp) {
	return
}

// Inserts a box into panel 
func (p *RootPanel) Append(b Box) {
	b.SetParent(p)
	p.children.PushBack(b)
	b.SetPosSize(p.x, p.y, p.width, p.height)
	b.SetName("mdtwm desktop")
	b.Window().Map()
	// Update desktops properties
	p.w.ChangeProp(xgb.PropModeReplace, AtomNetNumberOfDesktops,
		xgb.AtomCardinal, uint32(p.children.Len()))
	p.w.ChangeProp(xgb.PropModeAppend, AtomNetVirtualRoots, xgb.AtomWindow,
		b.Window())
}

func (p *RootPanel) InsertNextTo(b, mark Box, x, y int16) {
	// TODO: Proper implementation needed.
	panic("Unimplemented.")
}

func (p *RootPanel) Remove(b Box) {
	p.children.Remove(b)
}
