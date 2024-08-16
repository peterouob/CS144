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

func TestUnparseInt(t *testing.T) {
	unparser := NetUnparser[uint32]{}

	tests := []struct {
		name string
		val  uint32
		n    int
		want []byte
	}{
		{
			name: "Unparse 16-bit integer",
			val:  0x1234,
			n:    2,
			want: []byte{0x12, 0x34},
		},
		{
			name: "Unparse 32-bit integer",
			val:  0x12345678,
			n:    4,
			want: []byte{0x12, 0x34, 0x56, 0x78},
		},
		{
			name: "Unparse 8-bit integer",
			val:  0xFF,
			n:    1,
			want: []byte{0xFF},
		},
		{
			name: "Unparse 24-bit integer",
			val:  0x123456,
			n:    3,
			want: []byte{0x12, 0x34, 0x56},
		},
		{
			name: "Unparse zero value",
			val:  0x00000000,
			n:    4,
			want: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "Unparse maximum 16-bit integer",
			val:  0xFFFF,
			n:    2,
			want: []byte{0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret := make([]byte, 0)
			unparser.UnparseInt(&ret, tt.val, tt.n)

			if len(ret) != len(tt.want) {
				t.Errorf("UnparseInt() length = %v, want %v", len(ret), len(tt.want))
			}

			for i := range ret {
				if ret[i] != tt.want[i] {
					t.Errorf("UnparseInt() byte at index %d = %v, want %v", i, ret[i], tt.want[i])
				}
			}
		})
	}
}
