package mdtwm

import (
	"github.com/ziutek/mdtwm/xgb_patched"
	"math"
)

var struts = make(Struts)

func manage(w Window, panel ParentBox, startup bool) {
	d.Printf("Manage %s in %s", w, panel)
	_, class := w.Class()
	if cfg.Ignore.Contains(class) {
		return
	}
	if root.Children().BoxByWindow(w, true) != nil {
		d.Printf("  %s - alredy managed", w)
		return
	}
	attr := w.Attrs()
	if attr == nil {
		return
	}
	// During startup we obtain each existing window (not vievable too)
	if startup && attr.MapState != xgb.MapStateViewable {
		d.Print("  not vievable")
		return
	}
	wmState := atomList(w.Prop(AtomNetWmState, math.MaxUint32))
	if wmState.Contains(AtomNetWmStateHidden) {
		// TODO: Need better handling of hidden window
		d.Print("  Hidden")
		return
	}
	// Check strut property
	if struts.Update(w, true) {
		d.Printf("  window %s - strut", w)
		w.Map()
		return // For now we don't manage windows with strut property
	}
	wmType := atomList(w.Prop(AtomNetWmWindowType, math.MaxUint32))
	if wmType.Contains(AtomNetWmWindowTypeDock) {
		d.Printf("  window %s is of type dock", w)
		w.Map()
		return // For now we don't manage dock windows
	}
	// Don't manage if OverrideRedirect flag is set
	if attr.OverrideRedirect {
		d.Print("  OverrideRedirect")
		return
	}
	// NewWindowBox(w) changes some property of w so it can't be used before!
	b := NewBoxedWindow(w)
	if wmType.Contains(AtomNetWmWindowTypeDialog) ||
		wmType.Contains(AtomNetWmWindowTypeUtility) ||
		wmType.Contains(AtomNetWmWindowTypeToolbar) ||
		wmType.Contains(AtomNetWmWindowTypeSplash) {
		b.SetFloat(true)
	}
	tr_for := atomList(w.Prop(xgb.AtomWmTransientFor, math.MaxUint32))
	if len(tr_for) > 0 && tr_for[0] != xgb.WindowNone {
		b.SetFloat(true)
	}
	if wmState.Contains(AtomNetWmStateModal) {
		// TODO: Need better handling of modal windows
		b.SetFloat(true)
	}
	if cfg.Float.Contains(class) {
		b.SetFloat(true)
	}
	// Insert new box in a panel.
	if b.Float() {
		currentDesk.Append(b)
	} else {
		panel.Append(b)
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
		strut := w.Prop(AtomNetWmStrutPartial, math.MaxUint32)
		if strut == nil {
			return false
		}
		if strut.ValueLen != 12 {
			strut = w.Prop(AtomNetWmStrut, math.MaxUint32)
			if strut == nil {
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
	for p := root.Children().Front(); p != nil; p = p.Next() {
		p.SetPosSize(x, y, width, height)
	}
	return true
}
