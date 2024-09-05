package stream

type Queue struct {
	item []byte
}

type QueueInterface interface {
	Pop()
	Push()
}

func NewQueue(l int) *Queue {
	return &Queue{
		item: make([]byte, l),
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.item) == 0
}

func (q *Queue) Push(val byte) {
	q.item = append(q.item, val)
}

func (q *Queue) Pop() byte {
	if q.IsEmpty() {
		return ' '
	}
	v := q.item[0]
	q.item = q.item[1:]
	return v
}
