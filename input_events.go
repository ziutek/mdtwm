package main

import (
	"code.google.com/p/x-go-binding/xgb"
)

func keyPress(e xgb.KeyPressEvent) {
	d.Printf("%T: %+v", e, e)
	if e.State == cfg.ModMask {
		cmd, ok := cfg.Keys[e.Detail]
		if !ok {
			l.Panic("Unhandled key: ", e.Detail)
		}
		if err := cmd.Run(); err != nil {
			l.Printf("cmd(%s): %s", cmd.Param, err)
		}
	}
}

// Distinguishes following actions (they can be used in specified handler):
// 1. One click and move, in: MotionNotify, ButtonRelease
// 2. One long click without move, in: ButtonRelease
// 3. Two clicks and move, in: MotionNotify, ButtonRelease
// 4. Two clicks without move when second click is long, in: ButtonRelease
// 4. Three clicks and move, in: MotionNotify, ButtonRelease
// 5. Three clicks without move, in: ButtonRelease
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
		click.X, click.Y = e.RootX, e.RootY
		// We can't do any action on first ButtonPressEvent
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
			// Send right click to the box
			var (
				child Window
				ok    bool
			)
			e.EventX, e.EventY, child, _, ok = w.TranslateCoordinates(
				root.Window(), e.RootX, e.RootY,
			)
			if !ok {
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
		if currentBox == nil || currentBox.Window() == w {
			return // Box moved outside of desk or wasn't moved
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
		// Use left and right borders for change desktop
		_, _, rootWidth, _ := root.PosSize()
		switch e.RootX {
		case 0: // Left border
			// WarpPointer must be first, if not we obtain to many Motion events
			conn.WarpPointer(xgb.WindowNone, root.Window().Id(), 0, 0, 0, 0,
				rootWidth-2, e.RootY)
			prevDesk := currentDesk.Prev()
			if prevDesk == nil {
				prevDesk = root.Children().Back()
			}
			currentDesk = prevDesk.(*Panel)
			currentDesk.Raise()
			skipBorderEvents()
		case rootWidth - 1: // Right border
			// WarpPointer must be first, if not we obtain to many Motion events
			conn.WarpPointer(xgb.WindowNone, root.Window().Id(), 0, 0, 0, 0,
				1, e.RootY)
			nextDesk := currentDesk.Next()
			if nextDesk == nil {
				nextDesk = root.Children().Front()
			}
			currentDesk = nextDesk.(*Panel)
			currentDesk.Raise()
			skipBorderEvents()
		}
	case 2: // Two clicks and move
	case 3: // Three clicks
	}
	/*dx, dy := e.RootX-move.x, e.RootY-move.y
	x, y, w, h := move.b.PosSize()
	move.b.SetPosSize(x+dx, y+dy, w, h)
	move.x += dx
	move.y += dy*/
}

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
