package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"math"
)

var struts = make(Struts)

func manage(w Window, panel ParentBox, vievableOnly bool) {
	d.Printf("Manage %s in %s", w, panel)
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if root.Children().BoxByWindow(w, true) != nil {
		d.Printf("  %s - alredy managed", w)
		return
	}
	attr, err := w.Attrs()
	if err != nil {
		l.Print("Attrs: ", err)
		return
	}
	// During startup manage obtains each existing window with vievableOnly=true
	if vievableOnly && attr.MapState != xgb.MapStateViewable {
		d.Print("  not vievable")
		return
	}
	// Check strut property
	if struts.Update(w, true) {
		return // For now we don't manage windows with strut property
	}
	// Don't manage if OverrideRedirect flag is set
	if attr.OverrideRedirect {
		d.Print("  OverrideRedirect")
		return
	}
	// Check window type
	p, err := w.Prop(AtomNetWmWindowType, math.MaxUint32)
	if err != nil {
		l.Print("Prop: ", err)
		return
	}
	wm_type := atomList(p)
	if wm_type.Contains(AtomNetWmWindowTypeDock) {
		d.Printf("  window %s is of type dock", w)
		return // For now we don't manage dock windows
	}
	// NewWindowBox(w) changes some property of w so it can't be used before!
	b := NewBoxedWindow(w)
	if wm_type.Contains(AtomNetWmWindowTypeDialog) ||
		wm_type.Contains(AtomNetWmWindowTypeUtility) ||
		wm_type.Contains(AtomNetWmWindowTypeToolbar) ||
		wm_type.Contains(AtomNetWmWindowTypeSplash) {
		b.SetFloat(true)
	}
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
		d.Printf("  Window %s will be floating", w)
		currentDesk.Insert(b)
	} else {
		panel.Insert(b)
	}
}

func unmanage(w Window) {
	if b := root.Children().BoxByWindow(w, true); b != nil {
		b.Parent().Remove(b)
	}
	struts.Update(w, false)
}

type strutGeometry struct {
	left, right, top, bottom int16
}

type Struts map[Window]strutGeometry

func (s Struts) Update(w Window, add bool) bool {
	x, y, width, height := currentDesk.PosSize()
	if add {
		strut, err := w.Prop(AtomNetWmStrutPartial, math.MaxUint32)
		if err != nil {
			l.Print("Prop: ", err)
			return false
		}
		if strut.ValueLen != 12 {
			strut, err = w.Prop(AtomNetWmStrut, math.MaxUint32)
			if err != nil {
				l.Print("Prop: ", err)
				return false
			}
		}
		if strut.ValueLen != 4 && strut.ValueLen != 12 {
			return false
		}
		sa := atomList(strut)
		sg := strutGeometry{
			int16(sa[0]), int16(sa[1]), int16(sa[2]), int16(sa[3]),
		}
		x += sg.left
		width -= sg.left + sg.right
		y += int16(sg.top)
		height -= int16(sg.top + sg.bottom)
		s[w] = sg
	} else {
		sg, ok := s[w]
		if !ok {
			return false
		}
		x -= sg.left
		width += sg.left + sg.right
		y -= int16(sg.top)
		height += int16(sg.top + sg.bottom)
		delete(s, w)
	}
	// Change size and position for all desks
	i := root.Children().FrontIter(false)
	for p := i.Next(); p != nil; p = i.Next() {
		p.ReqPosSize(x, y, width, height)
	}
	return true
}
