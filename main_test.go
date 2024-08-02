package main

import "testing"

func TestDeque_PushBack(t *testing.T) {
	dq := newDeque(10)
	dq.PushBack("123")
	dq.PushBack("456")
	dq.PushBack("789")
	dq.PushBack("0")
	for k, v := range dq.item {
		t.Log(k, v)
	}
}

func TestDeque_PushFront(t *testing.T) {
	dq := newDeque(10)
	dq.PushFront("123")
	dq.PushFront("456")
	dq.PushFront("789")
	dq.PushFront("0")
	for k, v := range dq.item {
		t.Log(k, v)
	}
}

func TestDeque_PopBack(t *testing.T) {
	dq := newDeque(10)
	dq.PushBack("123")
	dq.PushBack("456")
	dq.PushBack("789")
	dq.PushBack("0")
	s, _ := dq.PopBack()
	t.Log(s)
	s, _ = dq.PopBack()
	t.Log(s)
	s, _ = dq.PopBack()
	t.Log(s)
}

func TestDeque_PopFront(t *testing.T) {
	dq := newDeque(10)
	dq.PushFront("123")
	dq.PushFront("456")
	dq.PushFront("789")
	dq.PushFront("0")
	s, _ := dq.PopFront()
	t.Log(s)
	s, _ = dq.PopFront()
	t.Log(s)
	s, _ = dq.PopFront()
	t.Log(s)
}

//go test -run={func Name} -cover -v ./
