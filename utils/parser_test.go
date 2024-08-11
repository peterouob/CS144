package utils

import (
	"testing"
)

func TestNetParser_U32(t *testing.T) {
	// Testing parsing 32-bit integer
	buffer := NewBuffer("\x00\x00\x00\x01\x00\x00\x00\x02")
	parser := NewNetParser[uint32](*buffer)

	val := parser.ParseInt(4)
	if val != 1 {
		t.Errorf("NetParser.ParseInt() = %v, want %v", val, 1)
	}

	buf := parser.Buffer()
	if buf.Size() != 4 {
		t.Errorf("Buffer.Size() = %v, want %v", buf.Size(), 4)
	}
}

func TestNetParser_U16(t *testing.T) {
	// Testing parsing 16-bit integer
	buffer := NewBuffer("\x00\x02")
	parser := NewNetParser[uint16](*buffer)

	val := parser.ParseInt(2)
	if val != 2 {
		t.Errorf("NetParser.ParseInt() = %v, want %v", val, 2)
	}

	buf := parser.Buffer()
	if buf.Size() != 0 {
		t.Errorf("Buffer.Size() = %v, want %v", buf.Size(), 0)
	}
}

func TestNetParser_ErrorHandling(t *testing.T) {
	// Testing error handling for short packet
	buffer := NewBuffer("\x00")
	parser := NewNetParser[uint32](*buffer)

	val := parser.ParseInt(4)
	if val != 0 || parser.GetError() != PacketTooShort {
		t.Errorf("Expected PacketTooShort error, got %v", AsString(parser.GetError()))
	}
}

func TestNetParser_RemovePrefix(t *testing.T) {
	buffer := NewBuffer("\x01\x02\x03\x04")
	parser := NewNetParser[uint8](*buffer)

	parser.RemovePrefix(2)

	buf := parser.Buffer()
	if buf.Size() != 2 {
		t.Errorf("Buffer.Size() after RemovePrefix = %v, want %v", buf.Size(), 2)
	}

	valParser := NewNetParser[uint16](parser.Buffer())
	val := valParser.ParseInt(2)
	if val != 0x0304 {
		t.Errorf("NetParser.ParseInt() = %v, want %v", val, 0x0304)
	}
}
