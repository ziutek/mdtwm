package main

import (
	"fmt"
)

type Geometry struct {
	X, Y, W, H, B int16
	// int16 for W, H and B (see "Why X Is Not Our Ideal Window System")
}

func (g Geometry) String() string {
	return fmt.Sprintf("(%d,%d,%d,%d,%d)", g.X, g.Y, g.W, g.H, g.B)
}

func (g Geometry) Resize(i int16) Geometry {
	return Geometry{g.X, g.Y, g.W + i, g.H + i, g.B}
}

func (g Geometry) ResizeWidth(i int16) Geometry {
	return Geometry{g.X, g.Y, g.W + i, g.H, g.B}
}

func (g Geometry) ResizeHeight(i int16) Geometry {
	return Geometry{g.X, g.Y, g.W, g.H + i, g.B}
}

func (g Geometry) ResizeBorder(i int16) Geometry {
	bb := i + i
	return Geometry{g.X, g.Y, g.W - bb, g.H - bb, g.B + i}
}

func (g Geometry) External() Geometry {
	bb := g.B + g.B
	return Geometry{g.X, g.Y, g.W + bb, g.H + bb, 0}
}


func (g Geometry) Position() (x, y int16) {
	return g.X, g.Y
}

func (g Geometry) Size() (width, height int16) {
	return g.W, g.H
}

type Orientation bool

const (
	Vertical = Orientation(true)
	Horizontal = Orientation(false)
)

func (o Orientation) String() string {
	if o == Vertical {
		return "vertical"
	}
	return "horizontal"
}

