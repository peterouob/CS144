package main

type DequeFunc interface {
	PushFront(string)
	PushBack(string)
	PopFront() (string, bool)
	PopBack() (string, bool)
}

type Deque struct {
	item []string
}

var _ DequeFunc = (*Deque)(nil)

func (d *Deque) PushFront(item string) {
	d.item = append([]string{item}, d.item...)
}

func (d *Deque) PushBack(item string) {
	d.item = append(d.item, item)
}

func (d *Deque) PopFront() (string, bool) {
	if len(d.item) == 0 {
		return "", false
	}
	frontEle := d.item[0]
	d.item = d.item[1:]
	return frontEle, true
}

func (d *Deque) PopBack() (string, bool) {
	if len(d.item) == 0 {
		return "", false
	}
	rearEle := d.item[len(d.item)-1]
	d.item = d.item[:len(d.item)-1]
	return rearEle, true
}

func newDeque(length int) *Deque {
	return &Deque{
		item: make([]string, length),
	}
}

type Stream struct {
	dq Deque
}

func NewStream(dq Deque) *Stream {
	return &Stream{
		dq: dq,
	}
}
