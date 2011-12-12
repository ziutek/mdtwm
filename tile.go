package main

// Update geometry for boxes in panel
func tile(panel *Box) {
	n := int16(panel.Children.Len())
	if n == 0 {
		return // there is no boxes in panel
	}
	borderSpace := 2 * cfg.BorderWidth
	pg := panel.Window.Geometry()
	switch panel.Type {
	case BoxTypePanelV:
		l.Print("tile V in: ", panel.Window)
		h := pg.H / n
		// Set new geometry for all boxes in panel
		g := Geometry{0, 0, pg.W, h}
		i := panel.Children.FrontIter(false)
		for n > 1 {
			b := i.Next()
			if b.Type == BoxTypeWindow {
				b.Window.SetGeometry(g.Resize(-borderSpace))
			} else {
				b.Window.SetGeometry(g)
			}
			g.Y += h
			n--
		}
		// Last window obtain all remaining space
		b := i.Next()
		g.H = pg.H - g.Y
		if b.Type == BoxTypeWindow {
			b.Window.SetGeometry(g.Resize(-borderSpace))
		} else {
			b.Window.SetGeometry(g)
		}
	case BoxTypePanelH:
		l.Print("tile H in:", panel.Window)
		w := pg.W / n
		// Set new geometry for all boxes in panel
		g := Geometry{0, 0, w, pg.H}
		i := panel.Children.FrontIter(false)
		for n > 1 {
			b := i.Next()
			if b.Type == BoxTypeWindow {
				b.Window.SetGeometry(g.Resize(-borderSpace))
			} else {
				b.Window.SetGeometry(g)
			}
			g.X += w
			n--
		}
		// Last window obtain all remaining space
		b := i.Next()
		g.W = pg.W - g.X
		if b.Type == BoxTypeWindow {
			b.Window.SetGeometry(g.Resize(-borderSpace))
		} else {
			b.Window.SetGeometry(g)
		}
	default:
		panic("Tile on unknown panel type")
	}
}
