package main

func tile(b *Box) {
	l.Print("Tile", b)
	// Obtain a number of tiling boxes in parent
	i, n := b.Children.FrontIter(false), uint16(0)
	for c := i.Next(); c != nil; c = i.Next() {
		if !c.Float {
			n++
		}
	}
	if n == 0 {
		return // there is no tiling boxes in b
	}
	// Calculate new geometry for boxes in parent
	bg := b.Window.Geometry()
	g := Geometry{0, 0, bg.W, bg.H / n}
	h := Int16(int(g.H)) // new height
	// Set new height for windows
	i, n = b.Children.FrontIter(false), n-1
	for c := i.Next(); c != nil; c = i.Next() {
		c.Window.SetGeometry(g)
		if n--; n > 0 {
			g.Y += h
		} else {
			// Last window obtain all remaining space
			g.Y = g.Y + h
			g.H = uint16(bg.Y-g.Y)
		}
	}
}
