package utils

import (
	"testing"
)

func TestNetParser_U32(t *testing.T) {
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

func TestNetUnparser_UnparseInt(t *testing.T) {
	t.Run("Unparse uint32", func(t *testing.T) {
		nu := NetUnparser[uint32]{}
		val := uint32(0x12345678)
		s := ""

		result := nu.UnparseInt(s, val, 4)

		expected := "\x12\x34\x56\x78"
		if result != expected {
			t.Errorf("UnparseInt() = %x, want %x", result, expected)
		}
	})

	t.Run("Unparse uint16", func(t *testing.T) {
		nu := NetUnparser[uint16]{}
		val := uint16(0x1234)
		s := ""

		result := nu.UnparseInt(s, val, 2)

		expected := "\x12\x34"
		if result != expected {
			t.Errorf("UnparseInt() = %x, want %x", result, expected)
		}
	})

	t.Run("Unparse uint8", func(t *testing.T) {
		nu := NetUnparser[uint8]{}
		val := uint8(0x12)
		s := ""

		result := nu.UnparseInt(s, val, 1)

		expected := "\x12"
		if result != expected {
			t.Errorf("UnparseInt() = %x, want %x", result, expected)
		}
	})

	t.Run("Unparse multiple types", func(t *testing.T) {
		s := ""
		nu32 := NetUnparser[uint32]{}
		nu16 := NetUnparser[uint16]{}
		nu8 := NetUnparser[uint8]{}

		s = nu32.UnparseInt(s, uint32(0x01020304), 4)
		s = nu16.UnparseInt(s, uint16(0x0506), 2)
		s = nu8.UnparseInt(s, uint8(0x07), 1)

		expected := "\x01\x02\x03\x04\x05\x06\x07"
		if s != expected {
			t.Errorf("UnparseInt() = %x, want %x", s, expected)
		}
	})
}
