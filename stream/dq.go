package stream

type DequeFunc interface {
	Empty() bool
	PushFront(byte)
	PushBack(byte)
	PopFront() (byte, bool)
	PopBack() (byte, bool)
	StringItem() []string
}

type Deque struct {
	item []byte
}

var _ DequeFunc = (*Deque)(nil)

func (d *Deque) Empty() bool {
	return len(d.item) == 0
}

func (d *Deque) PushFront(item byte) {
	d.item = append([]byte{item}, d.item...)
}

func (d *Deque) PushBack(item byte) {
	d.item = append(d.item, item)
}

func (d *Deque) PopFront() (byte, bool) {
	if len(d.item) == 0 {
		return ' ', false
	}
	frontEle := d.item[0]
	d.item = d.item[1:]
	return frontEle, true
}

func (d *Deque) PopBack() (byte, bool) {
	if len(d.item) == 0 {
		return ' ', false
	}
	rearEle := d.item[len(d.item)-1]
	d.item = d.item[:len(d.item)-1]
	return rearEle, true
}

func (d *Deque) StringItem() []string {
	var tmp []string
	for _, v := range d.item {
		tmp = append(tmp, string(v))
	}
	return tmp
}

func NewDeque(length int) *Deque {
	return &Deque{
		item: make([]byte, length),
	}
}
