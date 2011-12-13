package main

import (
	"math"
	"x-go-binding.googlecode.com/hg/xgb"
)

var x int16

func manage(w Window, panel ParentBox, vievableOnly bool) {
	l.Printf("manage %s in %s", w, panel)
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if root.Children().BoxByWindow(w, true) != nil {
		l.Printf("  %s - alredy managed", w)
		return
	}
	attr := w.Attrs()
	// Don't manage if OverrideRedirect flag is set
	if attr.OverrideRedirect {
		l.Print(  "OverrideRedirect")
		return
	}
	if vievableOnly && attr.MapState != xgb.MapStateViewable {
		l.Print(  "not vievable")
		return
	}
	// Check window type
	p, err := w.Prop(AtomNetWmWindowType, math.MaxUint32)
	if err == nil {
		wm_type := propReplyAtoms(p)
		if wm_type.Contains(AtomNetWmWindowTypeDock) {
			l.Printf("  window %s is of type dock", w)
		}
		if cfg.Float.Contains(class) ||
			wm_type.Contains(AtomNetWmWindowTypeDialog) ||
			wm_type.Contains(AtomNetWmWindowTypeUtility) ||
			wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
			wm_type.Contains(AtomNetWmWindowTypeSplash) {
			l.Printf(" window %s should be treated as float", w)
		}
	} else {
		l.Printf("  can't get AtomNetWmWindowType from %s: %s", w, err)
	}
	// Insert new box in a panel.
	// NewWindowBox(w) changes some property of w so it can't be used before!
	panel.Insert(NewTiledWindow(w))
}
