package main

func tile(b *Box) {
	l.Print("Tile", b)
	// Obtain a number of tiling boxes in parent
	i, n := b.Children.FrontIter(false), int16(0)
	for c := i.Next(); c != nil; c = i.Next() {
		if !c.Float {
			n++
		}
	}
	if n == 0 {
		return // there is no tiling boxes in b
	}
	// Calculate new geometry for boxes in parent
	borderSpace := 2 * cfg.BorderWidth
	bg := b.Window.Geometry()
	h := bg.H / n - 2 * cfg.BorderWidth // new height
	g := Geometry{0, 0, bg.W - borderSpace, h}
	// Set new height for windows
	i, n = b.Children.FrontIter(false), n-1
	for c := i.Next(); c != nil; c = i.Next() {
		c.Window.SetGeometry(g)
		if n--; n > 0 {
			g.Y += h
		} else {
			// Last window obtain all remaining space
			g.Y = g.Y + h
			g.H = bg.Y-g.Y - borderSpace
		}
	}
}
