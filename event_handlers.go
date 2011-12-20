package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

func handleEvent(event xgb.Event) {
	switch e := event.(type) {
	// *Request events
	case xgb.MapRequestEvent:
		mapRequest(e)
	case xgb.ConfigureRequestEvent:
		configureRequest(e)
	// *Notify events
	case xgb.EnterNotifyEvent:
		enterNotify(e)
	case xgb.UnmapNotifyEvent:
		unmapNotify(e)
	case xgb.DestroyNotifyEvent:
		destroyNotify(e)
	// Keyboard and mouse events
	case xgb.KeyPressEvent:
		keyPress(e)
	case xgb.ButtonPressEvent:
		buttonPress(e)
	case xgb.ButtonReleaseEvent:
		buttonRelease(e)
	case xgb.MotionNotifyEvent:
		motionNotify(e)
	default:
		l.Printf("*** Unhandled event: %T", e)
	}
}

func mapRequest(e xgb.MapRequestEvent) {
	l.Print("MapRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	manage(w, currentPanel(), false)
}

func enterNotify(e xgb.EnterNotifyEvent) {
	l.Print("EnterNotifyEvent: ", Window(e.Event))
	if e.Mode != xgb.NotifyModeNormal {
		return
	}
	changeFocusTo(Window(e.Event))
}

func destroyNotify(e xgb.DestroyNotifyEvent) {
	l.Print("DestroyNotifyEvent: ", Window(e.Event))
	removeWindow(Window(e.Event), false)
}

func unmapNotify(e xgb.UnmapNotifyEvent) {
	l.Print("xgb.UnmapNotifyEvent: ", Window(e.Event), e)
	removeWindow(Window(e.Event), false)
}

func configureRequest(e xgb.ConfigureRequestEvent) {
	l.Print("ConfigureRequestEvent: ", Window(e.Window))
	w := Window(e.Window)
	b := root.Children().BoxByWindow(w, true)
	if b == nil || b.Float() {
		// We accept request from floating windows.
		// Unmanaged window will be configured by manage() function so
		// now we can simply execute its request.
		mask := (xgb.ConfigWindowX | xgb.ConfigWindowY |
			xgb.ConfigWindowWidth | xgb.ConfigWindowHeight |
			xgb.ConfigWindowBorderWidth | xgb.ConfigWindowSibling |
			xgb.ConfigWindowStackMode) & e.ValueMask
		v := make([]interface{}, 0, 7)
		if mask&xgb.ConfigWindowX != 0 {
			v = append(v, e.X)
		}
		if mask&xgb.ConfigWindowY != 0 {
			v = append(v, e.Y)
		}
		if mask&xgb.ConfigWindowWidth != 0 {
			v = append(v, e.Width)
		}
		if mask&xgb.ConfigWindowHeight != 0 {
			v = append(v, e.Height)
		}
		if mask&xgb.ConfigWindowBorderWidth != 0 {
			v = append(v, e.BorderWidth)
		}
		if mask&xgb.ConfigWindowSibling != 0 {
			v = append(v, e.Sibling)
		}
		if mask&xgb.ConfigWindowStackMode != 0 {
			v = append(v, e.StackMode)
		}
		w.Configure(mask, v...)
		return
	}
	// Force box configuration
	g := b.Geometry()
	cne := xgb.ConfigureNotifyEvent{
		Event:        e.Window,
		Window:       e.Window,
		AboveSibling: xgb.WindowNone,
		X:            g.X,
		Y:            g.Y,
		Width:        Pint16(g.W),
		Height:       Pint16(g.H),
		BorderWidth:  Pint16(g.B),
	}
	w.Send(false, xgb.EventMaskStructureNotify, cne)
}

func keyPress(e xgb.KeyPressEvent) {
	l.Print("KeyPressEvent: ", Window(e.Event))
	if e.State == cfg.ModMask {
		cmd, ok := cfg.Keys[e.Detail]
		if !ok {
			l.Print("Unhandled key: ", e.Detail)
		}
		if err := cmd.Run(); err != nil {
			l.Printf("cmd(%s): %s", cmd.Param, err)
		}
	}
}


// Distinguishes following actions (they can be used in specified handler):
// 1. One click and move: Motion, ButtonRelease
// 2. One long click without move: ButtonRelease
// 3. Two clicks and move: Motion, ButtonRelease
// 4. Two clicks without move when second click is long: ButtonRelease
// 4. Three clicks and move: Motion, ButtonRelease
// 5. Three clicks without move: ButtonRelease
type Multiclick struct {
	Box   Box
	X, Y  int16
	Num   int  // number of clicks
	Moved bool // Indicates that cursor has moved during click

	last    xgb.Timestamp // time of last click
	counter int           // Multiclick counter
}

func (c *Multiclick) Inc(t xgb.Timestamp) {
	c.counter++
	if c.counter == 1 || t-c.last > cfg.MultiClickTime * 2 {
		// First click or to long interval betwen clicks
		c.last = t
		c.Num = 0
		c.Moved = false
	} else if c.counter == 3 {
		// Maximum number of clicks
		c.Num = c.counter
		c.counter = 0
	}
}

func (c *Multiclick) First() bool {
	return c.counter == 1
}

func (c *Multiclick) Update(t xgb.Timestamp, moved bool) {
	if !c.Moved {
		c.Moved = moved
	}
	if c.counter == 0 {
		return
	}
	if t-c.last > cfg.MultiClickTime || moved {
		// We obtained all clicks from this multiclick
		c.Num = c.counter
		c.counter = 0
	}
}

var click Multiclick

func buttonPress(e xgb.ButtonPressEvent) {
	l.Print("ButtonPressEvent: ", e.Event)
	click.Inc(e.Time)
	if click.First() {
		// Save first clicked box and coordinates of first click
		click.Box = currentBox
		click.X, click.Y = e.RootX, e.RootY
		// We can't do any action on first ButtonPressEvent
		return
	}
	click.Update(e.Time, false)
}

func buttonRelease(e xgb.ButtonReleaseEvent) {
	l.Print("ButtonReleaseEvent: ", e.Event)
	click.Update(e.Time, false)
	// Actions
	switch click.Num {
	case 1: // One click
		if _, ok := click.Box.(ParentBox); ok {
			return // For now, we don't move panels
		}
		w := click.Box.Window()
		if !click.Moved {
			// Send right click to the box
			var (
				child Window
				err error
			)
			e.EventX, e.EventY, child, _, err = w.TranslateCoordinates(
				root.Window(), e.RootX, e.RootY,
			)
			if err != nil {
				l.Print("buttonRelease: ", err)
				return
			}
			e.Event = w.Id()
			e.Child = child.Id()
			e.Time = xgb.TimeCurrentTime
			e.State = 0
			w.Send(false, xgb.EventMaskNoEvent, xgb.ButtonPressEvent(e))
			e.State = xgb.EventMaskButton3Motion
			w.Send(false, xgb.EventMaskNoEvent, e)
			return
		}
		if currentBox.Window() == w {
			return // Box wasn't moved
		}
		// Border color will be set properly by EnterNotify event.
		w.SetBorderColor(cfg.NormalBorderColor)
		// Move a box
		click.Box.Parent().Remove(click.Box, false)
		currentPanel().Insert(click.Box)
	case 2: // Two clicks
	case 3: // Three clicks
		if !click.Moved {
			if _, ok := click.Box.(ParentBox); ok {
				return // For now we don't destroy panels
			}
			// TODO: Use A_WM_DELETE_WINDOW if box supports it
			click.Box.Window().Destroy()
		}
	}
}

func motionNotify(e xgb.MotionNotifyEvent) {
	//l.Print("xgb.MotionNotifyEvent: ", Window(e.Event))
	dx := int(e.RootX - click.X)
	dy := int(e.RootY - click.Y)
	click.Update(e.Time, dx*dx+dy*dy > cfg.MovedClickRadius)
	// Actions
	switch click.Num {
	case 1: // One click and move
		if _, ok := click.Box.(ParentBox); ok {
			return // For now, we don't move panels
		}
		conn.ChangeActivePointerGrab(cfg.MoveCursor, xgb.TimeCurrentTime,
			rightButtonEventMask)
	case 2: // Two clicks and move
	case 3: // Three clicks
	}
	/*dx, dy := e.RootX-move.x, e.RootY-move.y
	x, y, w, h := move.b.PosSize()
	move.b.SetPosSize(x+dx, y+dy, w, h)
	move.x += dx
	move.y += dy*/
}
