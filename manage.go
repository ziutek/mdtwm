package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"math"
)

func manage(w Window, panel ParentBox, vievableOnly bool) {
	l.Printf("*** Manage %s in %s", w, panel)
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if root.Children().BoxByWindow(w, true) != nil {
		l.Printf("***   %s - alredy managed", w)
		return
	}
	attr, err := w.Attrs()
	if err != nil {
		l.Print("Attrs: ", err)
		return
	}
	// Don't manage if OverrideRedirect flag is set
	if attr.OverrideRedirect {
		l.Print("***   OverrideRedirect")
		return
	}
	if vievableOnly && attr.MapState != xgb.MapStateViewable {
		l.Print("***   not vievable")
		return
	}
	p, err := w.Prop(AtomNetWmWindowType, math.MaxUint32)
	if err != nil {
		l.Print("Prop: ", err)
		return
	}
	wm_type := atomList(p)
	if wm_type.Contains(AtomNetWmWindowTypeDock) {
		l.Printf("***   window %s is of type dock", w)
		strut, err := w.Prop(AtomNetWmStrutPartial, math.MaxUint32)
		if err != nil {
			l.Print("Prop: ", err)
			return
		}
		if err == nil && strut.ValueLen == 12 {
			sa := atomList(strut)
			left, right, top, bottom := sa[0], sa[1], sa[2], sa[3]
			x, y, width, height := currentDesk.PosSize()
			x += int16(left)
			width -= int16(left + right)
			y += int16(top)
			height -= int16(top + bottom)
			// Change size and position for all desks
			i := root.Children().FrontIter(false)
			for d := i.Next(); d != nil; d = i.Next() {
				l.Printf("***   Change size of %s: %d %d %d %d",
					d, x, y, width, height)
				d.ReqPosSize(x, y, width, height)
			}
		}
		return
	}
	// NewWindowBox(w) changes some property of w so it can't be used before!
	b := NewBoxedWindow(w)
	if wm_type.Contains(AtomNetWmWindowTypeDialog) ||
		wm_type.Contains(AtomNetWmWindowTypeUtility) ||
		wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
		wm_type.Contains(AtomNetWmWindowTypeSplash) {
		b.SetFloat(true)
	}
	// Check window type
	p, err = w.Prop(xgb.AtomWmTransientFor, math.MaxUint32)
	if err != nil {
		l.Print("Prop: ", err)
		return
	}
	tr_for := atomList(p)
	if len(tr_for) > 0 && tr_for[0] != xgb.WindowNone {
		b.SetFloat(true)
	}
	if cfg.Float.Contains(class) {
		b.SetFloat(true)
	}
	// Insert new box in a panel.
	if b.Float() {
		l.Printf("***   Window %s will be floating", w)
		currentDesk.Insert(b)
	} else {
		panel.Insert(b)
	}
}
