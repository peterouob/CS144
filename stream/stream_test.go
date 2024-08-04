package stream

import (
	"fmt"
	"testing"
)

func TestDeque_PushBack(t *testing.T) {
	dq := newDeque(10)
	dq.PushBack('1')
	dq.PushBack('2')
	dq.PushBack('3')
	dq.PushBack('0')
	for k, v := range dq.StringItem() {
		t.Log(k, v)
	}
}

func TestDeque_PushFront(t *testing.T) {
	dq := newDeque(10)
	dq.PushFront('1')
	dq.PushFront('2')
	dq.PushFront('3')
	dq.PushFront('0')
	for k, v := range dq.item {
		t.Log(k, v)
	}
}

func TestDeque_PopBack(t *testing.T) {
	dq := newDeque(10)
	dq.PushBack('1')
	dq.PushBack('2')
	dq.PushBack('3')
	dq.PushBack('0')
	s, _ := dq.PopBack()
	t.Log(s)
	s, _ = dq.PopBack()
	t.Log(s)
	s, _ = dq.PopBack()
	t.Log(s)
}

func TestDeque_PopFront(t *testing.T) {
	dq := newDeque(10)
	dq.PushFront('1')
	dq.PushFront('2')
	dq.PushFront('3')
	dq.PushFront('0')
	s, _ := dq.PopFront()
	t.Log(s)
	s, _ = dq.PopFront()
	t.Log(s)
	s, _ = dq.PopFront()
	t.Log(s)
}

func TestNewStream(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)

	if stream.capacitySize != 10 {
		t.Errorf("Expected capacitySize to be 10, got %d", stream.capacitySize)
	}

	if stream.writtenSize != 0 {
		t.Errorf("Expected writtenSize to be 0, got %d", stream.writtenSize)
	}

	if stream.readSize != 0 {
		t.Errorf("Expected readSize to be 0, got %d", stream.readSize)
	}

	if stream.endInput {
		t.Errorf("Expected endInput to be false, got %v", stream.endInput)
	}

	if stream.error {
		t.Errorf("Expected error to be false, got %v", stream.error)
	}
}

func TestWrite(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	data := "hello"

	n := stream.Write(data)

	if n != len(data) {
		t.Errorf("Expected written length to be %d, got %d", len(data), n)
	}

	if stream.writtenSize != len(data) {
		t.Errorf("Expected writtenSize to be %d, got %d", len(data), stream.writtenSize)
	}

	if stream.BufferSize() != len(data) {
		t.Errorf("Expected buffer size to be %d, got %d", len(data), stream.BufferSize())
	}
}

func TestPeekOutput(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	data := "hello"
	stream.Write(data)

	output := fmt.Sprintf("%s", stream.PeekOutput(3))
	expected := "hel"

	if output != expected {
		t.Errorf("Expected output to be %s, got %s", expected, output)
	}
}

func TestRead(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	data := "hello"
	stream.Write(data)

	readData := fmt.Sprintf("%s", stream.Read(3))
	expected := "hel"

	if readData != expected {
		t.Errorf("Expected read data to be %s, got %s", expected, readData)
	}

	if stream.BytesRead() != 3 {
		t.Errorf("Expected bytes read to be 3, got %d", stream.BytesRead())
	}

	if stream.BufferSize() != 2 {
		t.Errorf("Expected buffer size to be 2, got %d", stream.BufferSize())
	}
}

// Test sequential write and read
func TestSequentialWriteAndRead(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	data := "hello"

	n := stream.Write(data)
	if n != len(data) {
		t.Errorf("Expected written length to be %d, got %d", len(data), n)
	}

	readData := stream.Read(len(data))
	if readData != data {
		t.Errorf("Expected read data to be %s, got %s", data, readData)
	}
}

// Test end input and EOF detection
func TestEndInputAndEOF(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	data := "hello"

	stream.Write(data)
	stream.EndInput()

	if !stream.InputEnded() {
		t.Errorf("Expected input to be ended")
	}

	readData := stream.Read(len(data))
	if readData != data {
		t.Errorf("Expected read data to be %s, got %s", data, readData)
	}

	if !stream.EOF() {
		t.Errorf("Expected EOF to be true")
	}
}

// Test flow control
func TestFlowControl(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 5, 0, 0, false, false)
	data := "hello"

	n := stream.Write(data)
	if n != 5 {
		t.Errorf("Expected written length to be 5, got %d", n)
	}

	m := stream.Write("world")
	if m != 0 {
		t.Errorf("Expected written length to be 0, got %d", m)
	}

	stream.Read(3)
	n = stream.Write("world")
	if n != 3 {
		t.Errorf("Expected written length to be 3, got %d", n)
	}
}

// Test long byte stream
func TestLongByteStream(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 1, 0, 0, false, false)
	data := "hello"

	n := stream.Write(data)
	if n != 1 {
		t.Errorf("Expect write data to be 0,got %d", n)
	}

	readData := stream.Read(1)

	if readData != string(data[0]) {
		t.Errorf("Expected read data to be %s, got %s", data, readData)
	}
}

// Test the SetError and Errors functions
func TestSetError(t *testing.T) {
	q := Deque{}
	stream := NewStream(q, 10, 0, 0, false, false)
	stream.SetError()

	if !stream.Errors() {
		t.Errorf("Expected error to be true, got %v", stream.Errors())
	}
}

//go test -run={func Name} -cover -v ./
// for all go test -v
