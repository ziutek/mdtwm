package main

import (
	"container/list"
)

type Box struct {
	Frame  Window  // frame window of this box
	Window Window  // window stored in this box
	Childs BoxList // child boxes contains childs of windows

	Float bool // floating box
}

func NewBox() *Box {
	return &Box{Childs: NewBoxList()}
}

type BoxList struct {
	raw *list.List
}

func NewBoxList() BoxList {
	return BoxList{list.New()}
}

func (bl BoxList) PushFront(f *Box) {
	bl.raw.PushFront(f)
}

func (bl BoxList) PushBack(f *Box) {
	bl.raw.PushBack(f)
}

func (bl BoxList) Len() int {
	return bl.raw.Len()
}

func (bl BoxList) FrontIter(full_tree bool) BoxListIterator {
	return &frontBoxListIterator{
		boxListIterator{bl.raw.Front(), full_tree, nil},
	}
}

func (bl BoxList) BackIter(full_tree bool) BoxListIterator {
	return &backBoxListIterator{
		boxListIterator{bl.raw.Back(), full_tree, nil},
	}
}

func (bl BoxList) BoxByWindow(w Window) *Box {
	for e := bl.raw.Front(); e != nil; e = e.Next() {
		b := e.Value.(*Box)
		if b.Window == w {
			return b
		}
		b = b.Childs.BoxByWindow(w)
		if b != nil {
			return b
		}
	}
	return nil
}

func (bl BoxList) Remove(f *Box) {
	for e := bl.raw.Front(); e != nil; e = e.Next() {
		if e.Value.(*Box) == f {
			bl.raw.Remove(e)
			return
		}
	}
	panic("Can't remove non existent frame form a list")
}

type BoxListIterator interface {
	Next() *Box
	Done() bool
}

type boxListIterator struct {
	current *list.Element
	full_tree bool
	child BoxListIterator
}

func (i *boxListIterator) Done() bool {
	return i.current == nil && i.child == nil
}

type frontBoxListIterator struct {
	boxListIterator
}

// Returns nil if end of list
func (i *frontBoxListIterator) Next() (b *Box) {
	if i.child != nil {
		b = i.child.Next()
		if b != nil {
			return
		}
		i.child = nil // There is no more data in child iterator
	}
	if i.current == nil {
		return	// There is no more data at all
	}
	b = i.current.Value.(*Box)
	i.current = i.current.Next()
	if i.full_tree {
		i.child = b.Childs.FrontIter(i.full_tree)
	}
	return
}

type backBoxListIterator struct {
	boxListIterator
}

// Returns nil if end of list
func (i *backBoxListIterator) Next() (b *Box) {
	if i.child != nil {
		b = i.child.Next()
		if b != nil {
			return
		}
		i.child = nil // There is no more data in child iterator
	}
	if i.current == nil {
		return	// There is no more data at all
	}
	b = i.current.Value.(*Box)
	i.current = i.current.Prev()
	if i.full_tree {
		i.child = b.Childs.BackIter(i.full_tree)
	}
	return
}
