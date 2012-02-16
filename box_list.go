package mdtwm

type BoxList struct {
	front, back Box
	length      int
}

func NewBoxList() *BoxList {
	return new(BoxList)
}

func (bl *BoxList) Front() Box {
	return bl.front
}

func (bl *BoxList) Back() Box {
	return bl.back
}

func (bl *BoxList) Len() int {
	return bl.length
}

func (bl *BoxList) PushFront(b Box) {
	if bl.front == nil {
		bl.front, bl.back = b, b
		b.SetPrev(nil)
		b.SetNext(nil)
		b.SetList(bl)
		bl.length = 1
		return
	}
	bl.InsertBefore(b, bl.front)
}

func (bl *BoxList) PushBack(b Box) {
	if bl.back == nil {
		bl.front, bl.back = b, b
		b.SetPrev(nil)
		b.SetNext(nil)
		b.SetList(bl)
		bl.length = 1
		return
	}
	bl.InsertAfter(b, bl.back)
}

// Returns false if b not in bl
func (bl *BoxList) InsertBefore(b, mark Box) bool {
	if mark.List() != bl {
		return false
	}
	if mark.Prev() == nil {
		bl.front = b
	} else {
		mark.Prev().SetNext(b)
	}
	b.SetPrev(mark.Prev())
	mark.SetPrev(b)
	b.SetNext(mark)
	b.SetList(bl)
	bl.length++
	return true
}

// Returns false if mark not in bl
func (bl *BoxList) InsertAfter(b, mark Box) bool {
	if mark.List() != bl {
		return false
	}
	if mark.Next() == nil {
		bl.back = b
	} else {
		mark.Next().SetPrev(b)
	}
	b.SetNext(mark.Next())
	mark.SetNext(b)
	b.SetPrev(mark)
	b.SetList(bl)
	bl.length++
	return true
}

func (bl *BoxList) BoxByWindow(w Window, full_tree bool) Box {
	for b := bl.Front(); b != nil; b = b.Next() {
		if b.Window() == w {
			return b
		}
		if full_tree {
			c := b.Children().BoxByWindow(w, true)
			if c != nil {
				return c
			}
		}
	}
	return nil
}

func (bl *BoxList) Remove(b Box) {
	if b.List() != bl {
		l.Panic("Can't remove a non existent box form a list")
	}
	if b.Prev() == nil {
		bl.front = b.Next()
	} else {
		b.Prev().SetNext(b.Next())
	}
	if b.Next() == nil {
		bl.back = b.Prev()
	} else {
		b.Next().SetPrev(b.Prev())
	}
	b.SetPrev(nil)
	b.SetNext(nil)
	b.SetList(nil)
	bl.length--
}

func (bl *BoxList) FrontIter() BoxListIterator {
	return &frontBoxListIterator{
		boxListIterator{bl.Front(), nil},
	}
}

func (bl *BoxList) BackIter() BoxListIterator {
	return &backBoxListIterator{
		boxListIterator{bl.Back(), nil},
	}
}

// Iterates ovet full tree
type BoxListIterator interface {
	Next() Box
}

type boxListIterator struct {
	current   Box
	child     BoxListIterator
}

type frontBoxListIterator struct {
	boxListIterator
}

// Returns nil if end of list
func (i *frontBoxListIterator) Next() (b Box) {
	if i.child != nil {
		b = i.child.Next()
		if b != nil {
			return
		}
		i.child = nil // There is no more data in child iterator
	}
	if i.current == nil {
		return // There is no more data at all
	}
	b = i.current
	i.current = i.current.Next()
	i.child = b.Children().FrontIter()
	return
}

type backBoxListIterator struct {
	boxListIterator
}

// Returns nil if end of list
func (i *backBoxListIterator) Next() (b Box) {
	if i.child != nil {
		b = i.child.Next()
		if b != nil {
			return
		}
		i.child = nil // There is no more data in child iterator
	}
	if i.current == nil {
		return // There is no more data at all
	}
	b = i.current
	i.current = i.current.Prev()
	i.child = b.Children().BackIter()
	return
}
