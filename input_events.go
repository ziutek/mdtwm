package main

import (
	"github.com/ziutek/mdtwm/xgb_patched"
)

func keyPress(e xgb.KeyPressEvent) {
	d.Printf("%T: %+v", e, e)
	if e.State == cfg.ModMask {
		cmd, ok := cfg.Keys[keyCodeToSym[e.Detail]]
		if !ok {
			l.Print("Unhandled key: ", e.Detail)
			return
		}
		if err := cmd.Run(); err != nil {
			l.Printf("cmd(%s): %s", cmd.Param, err)
		}
	}
}


// TODO: Following code isn't good (it need to be reimplemented!)

const (
	resizeLeft = 1 << iota
	resizeRight
	resizeTop
	resizeBottom
)

// Distinguishes following actions (they can be used in specified handler):
// 1. One click and move, in: MotionNotify, ButtonRelease
// 2. One long click without move, in: ButtonRelease
// 3. Two clicks and move, in: MotionNotify, ButtonRelease
// 4. Two clicks without move when second click is long, in: ButtonRelease
// 4. Three clicks and move, in: MotionNotify, ButtonRelease
// 5. Three clicks without move, in: ButtonRelease
type Multiclick struct {
	Box                Box
	X, Y, RootX, RootY int16
	Num                int  // number of clicks
	Moved              bool // indicates that cursor has moved during click
	Resize             byte
	Child              Window

	last    xgb.Timestamp // time of last click
	counter int           // Multiclick counter
}

func (c *Multiclick) Inc(t xgb.Timestamp) {
	c.counter++
	if c.counter == 1 || t-c.last > cfg.MultiClickTime*2 {
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
	d.Printf("%T: %+v", e, e)
	click.Inc(e.Time)
	if click.First() {
		// Save first clicked box and coordinates of first click
		click.Box = currentBox
		click.RootX, click.RootY = e.RootX, e.RootY
		click.Resize = 0
		if click.Box.Float() {
			click.Box.Raise()
		}
		// Check for click in resize border
		w := click.Box.Window()
		var ok bool
		click.X, click.Y, click.Child, _, ok = w.TranslateCoordinates(
			root.Window(), e.RootX, e.RootY,
		)
		if !ok {
			return
		}
		g := click.Box.Geometry()
		rbw := cfg.ResizeBorderWidth - cfg.BorderWidth
		if click.X < rbw {
			click.Resize |= resizeLeft
		} else if click.X >= g.W-rbw {
			click.Resize |= resizeRight
		}
		if click.Y < rbw {
			click.Resize |= resizeTop
		} else if click.Y >= g.H-rbw {
			click.Resize |= resizeBottom
		}
		return
	}
	click.Update(e.Time, false)
}

func buttonRelease(e xgb.ButtonReleaseEvent) {
	d.Printf("%T: %+v", e, e)
	click.Update(e.Time, false)
	// Actions
	switch click.Num {
	case 1: // One click
		if _, ok := click.Box.(ParentBox); ok {
			return // For now, we don't move panels
		}
		w := click.Box.Window()
		if !click.Moved {
			e.Event = w.Id()
			e.EventX = click.X
			e.EventY = click.Y
			e.Child = click.Child.Id()
			e.Time = xgb.TimeCurrentTime
			e.State = 0
			w.Send(false, xgb.EventMaskNoEvent, xgb.ButtonPressEvent(e))
			e.State = xgb.EventMaskButton3Motion
			w.Send(false, xgb.EventMaskNoEvent, e)
			return
		}
		if click.Box.Float() || currentBox == nil || currentBox.Window() == w {
			// There isn't any action for floating box or if box have been
			// moved outside of desk or havent moved
			return
		}
		// Move a box
		click.Box.Parent().Remove(click.Box)
		x, y, _, _, ok := currentBox.Window().TranslateCoordinates(
			root.Window(), e.RootX, e.RootY,
		)
		if !ok {
			currentPanel().Append(click.Box)
			return
		}
		currentPanel().InsertNextTo(click.Box, currentBox, x, y)
	case 2: // Two clicks
	case 3: // Three clicks
		if click.Moved {
			return
		}
		switch b := click.Box.(type) {
		case *BoxedWindow:
			if b.Protocols().Contains(AtomWmDeleteWindow) {
				b.SendMessage(AtomWmDeleteWindow, b.Window())
			} else {
				b.Window().Destroy()
			}
		}
	}
}

func motionNotify(e xgb.MotionNotifyEvent) {
	//d.Printf("%T: %+v", e, e)
	dx := e.RootX - click.RootX
	dy := e.RootY - click.RootY
	click.Update(e.Time, dx*dx+dy*dy > cfg.MovedClickRadius)
	// Actions
	switch click.Num {
	case 2: // Two clicks and move: tile/untile box
		click.Num = 1
		click.Resize = 0
		click.Box.SetFloat(!click.Box.Float())
		fallthrough
	case 1: // One click and move: move/resize box
		if _, ok := click.Box.(ParentBox); ok {
			return // For now, we don't move panels
		}
		conn.ChangeActivePointerGrab(cfg.MoveCursor, xgb.TimeCurrentTime,
			rightButtonEventMask)
		if click.Resize != 0 {
			if click.Box.Float() {
				x, y, w, h := click.Box.PosSize()
				bb := cfg.ResizeBorderWidth * 2
				if click.Resize & resizeLeft != 0 && w - dx > bb {
						x += dx
						w -= dx
				} else if click.Resize & resizeRight != 0 && w + dx > bb {
						w += dx
				}
				if click.Resize & resizeTop != 0 && h - dy > bb {
						y += dy
						h -= dy
				} else if click.Resize & resizeBottom != 0 && h + dy > bb {
					h += dy
				}
				click.Box.SetPosSize(x, y, w, h)
				click.RootX += dx
				click.RootY += dy
			} else {
				// TODO: implement resize of tiled window
			}
			return
		}
		// Use left and right borders for change desktop
		_, _, rootWidth, _ := root.PosSize()
		switch e.RootX {
		case 0: // Left border
			// WarpPointer must be first, if not we obtain to many Motion events
			e.RootX = rootWidth - 2
			conn.WarpPointer(xgb.WindowNone, root.Window().Id(), 0, 0, 0, 0,
				e.RootX, e.RootY)
			setPrevDesk()
			skipBorderEvents()
		case rootWidth - 1: // Right border
			// WarpPointer must be first, if not we obtain to many Motion events
			e.RootX = 1
			conn.WarpPointer(xgb.WindowNone, root.Window().Id(), 0, 0, 0, 0,
				e.RootX, e.RootY)
			setNextDesk()
			skipBorderEvents()
		}
		if click.Box.Float() {
			// Move floating box
			x, y, w, h := click.Box.PosSize()
			dx, dy := e.RootX-click.RootX, e.RootY-click.RootY
			click.Box.SetPosSize(x+dx, y+dy, w, h)
			click.RootX += dx
			click.RootY += dy
			if click.Box.Parent().Window() != currentDesk.Window() {
				// Floating window moved to new desk
				click.Box.Parent().Remove(click.Box)
				currentDesk.Append(click.Box)
			}
		}
	case 3: // Three clicks
	}
}

/*func isRootOrDesk(b Box) bool {
	return b.Parent() == nil || b.Parent().Window() == root.Window()
}*/

func skipBorderEvents() {
	_, _, maxX, maxY := root.PosSize()
	maxX--
	maxY--
	var (
		event xgb.Event
		err   error
	)
	for {
		event, err = conn.WaitForEvent()
		if err != nil {
			break
		}
		e, ok := event.(xgb.MotionNotifyEvent)
		if !ok {
			break
		}
		if e.RootX != 0 && e.RootX != maxX && e.RootY != 0 && e.RootY != maxY {
			break
		}
	}
	handleEvent(event, err)
}
