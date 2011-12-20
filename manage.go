package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"math"
)

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
	attr, err := w.Attrs()
	if err != nil {
		l.Print("manage: ", err)
	}
	// Don't manage if OverrideRedirect flag is set
	if attr.OverrideRedirect {
		l.Print("OverrideRedirect")
		return
	}
	if vievableOnly && attr.MapState != xgb.MapStateViewable {
		l.Print("not vievable")
		return
	}
	// NewWindowBox(w) changes some property of w so it can't be used before!
	b := NewBoxedWindow(w)
	// Check window type
	p, err := w.Prop(xgb.AtomWmTransientFor, math.MaxUint32)
	if err == nil {
		tr_for := atomList(p)
		l.Printf("  tr_for: %+v %v", p, tr_for)
		if len(tr_for) > 0 && tr_for[0] != xgb.WindowNone {
			b.SetFloat(true)
		}
	}
	p, err = w.Prop(AtomNetWmWindowType, math.MaxUint32)
	if err == nil {
		wm_type := atomList(p)
		if wm_type.Contains(AtomNetWmWindowTypeDock) {
			l.Printf("  window %s is of type dock", w)
		}
		if cfg.Float.Contains(class) ||
			wm_type.Contains(AtomNetWmWindowTypeDialog) ||
			wm_type.Contains(AtomNetWmWindowTypeUtility) ||
			wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
			wm_type.Contains(AtomNetWmWindowTypeSplash) {
			b.SetFloat(true)
		}
	}
	// Insert new box in a panel.
	if b.Float() {
		l.Printf(" window %s should be treated as float", w)
		currentDesk.Insert(b)
	} else {
		panel.Insert(b)
	}
}
