package tcp_helper

import (
	"lab/utils"
	"testing"
)

func TestTCPHeaderParse(t *testing.T) {
	data := []byte{
		0x00, 0x50, // sport (80)
		0x00, 0x50, // dport (80)
		0x00, 0x00, 0x00, 0x01, // seqno (1)
		0x00, 0x00, 0x00, 0x02, // ackno (2)
		0x50,       // doff + reserved (5 << 4)
		0x18,       // flags (PSH + ACK)
		0x00, 0x10, // win (16)
		0x00, 0x3, // cksum (0x003)
		0x00, 0x00, // uptr
	}
	buffer := utils.NewBuffer(string(data))
	parser := utils.NewNetParser[uint8](*buffer)

	// 初始化 TCPHeader
	tcpHeader := TCPHeader[uint8]{}
	result := tcpHeader.Parse(*parser)

	if result != utils.NoError {
		t.Errorf("Parse() result = %v, want %v", result, utils.NoError)
	}
	if tcpHeader.sport != 80 {
		t.Errorf("tcpHeader.sport = %v, want %v", tcpHeader.sport, 80)
	}
	if tcpHeader.dport != 80 {
		t.Errorf("tcpHeader.dport = %v, want %v", tcpHeader.dport, 80)
	}
	if tcpHeader.seqno != 1 {
		t.Errorf("tcpHeader.seqno = %v, want %v", tcpHeader.seqno, 1)
	}
	if tcpHeader.ackno != 2 {
		t.Errorf("tcpHeader.ackno = %v, want %v", tcpHeader.ackno, 2)
	}
	if tcpHeader.doff != 5 {
		t.Errorf("tcpHeader.doff = %v, want %v", tcpHeader.doff, 5)
	}
	if !tcpHeader.psh || !tcpHeader.ack {
		t.Errorf("tcpHeader flags incorrect: psh = %v, ack = %v", tcpHeader.psh, tcpHeader.ack)
	}
	if tcpHeader.win != 16 {
		t.Errorf("tcpHeader.win = %v, want %v", tcpHeader.win, 16)
	}
	if tcpHeader.cksum != 0x003 {
		t.Errorf("tcpHeader.cksum = %v, want %v", tcpHeader.cksum, 0x123)
	}
	if tcpHeader.uptr != 0 {
		t.Errorf("tcpHeader.uptr = %v, want %v", tcpHeader.uptr, 0)
	}
}

func TestTCPHeaderParse_HeaderTooShort(t *testing.T) {
	data := []byte{
		0x00, 0x14, // sport (20)
		0x00, 0x50, // dport (80)
		0x00, 0x00, 0x00, 0x01, // seqno (1)
	}
	buffer := utils.NewBuffer(string(data))
	parser := utils.NewNetParser[uint8](*buffer)

	tcpHeader := TCPHeader[uint8]{}
	result := tcpHeader.Parse(*parser)

	if result != utils.PacketTooShort {
		t.Errorf("Parse() result = %v, want %v", result, utils.PacketTooShort)
	}
}

func TestTCPHeaderParse_InvalidDOFF(t *testing.T) {
	data := []byte{
		0x00, 0x14, // sport (20)
		0x00, 0x50, // dport (80)
		0x00, 0x00, 0x00, 0x01, // seqno (1)
		0x00, 0x00, 0x00, 0x00, // ackno (0)
		0x40,       // doff + reserved (invalid, should be >= 5)
		0x00,       // flags (none)
		0x00, 0x10, // win (16)
		0x12, 0x34, // cksum (0x1234)
		0x00, 0x00, // uptr
	}
	buffer := utils.NewBuffer(string(data))
	parser := utils.NewNetParser[uint8](*buffer)

	tcpHeader := TCPHeader[uint8]{}
	result := tcpHeader.Parse(*parser)

	if result != utils.PacketTooShort {
		t.Errorf("Parse() result = %v, want %v", result, utils.PacketTooShort)
	}
}

func TestTCPHeaderParse_WithOptions(t *testing.T) {
	data := []byte{
		0x00, 0x14, // sport (20)
		0x00, 0x50, // dport (80)
		0x00, 0x00, 0x00, 0x01, // seqno (1)
		0x00, 0x00, 0x00, 0x00, // ackno (0)
		0x60,       // doff + reserved (6 << 4)
		0x00,       // flags (none)
		0x00, 0x10, // win (16)
		0x12, 0x34, // cksum (0x1234)
		0x00, 0x00, // uptr
		0x01, 0x02, 0x03, 0x04, // 选项字段（4字节）
	}
	buffer := utils.NewBuffer(string(data))
	parser := utils.NewNetParser[uint8](*buffer)

	tcpHeader := TCPHeader[uint8]{}
	result := tcpHeader.Parse(*parser)

	if result != utils.NoError {
		t.Errorf("Parse() result = %v, want %v", result, utils.NoError)
	}
	if tcpHeader.doff != 6 {
		t.Errorf("tcpHeader.doff = %v, want %v", tcpHeader.doff, 6)
	}
}

func TestTCPHeaderParse_InvalidFlags(t *testing.T) {
	data := []byte{
		0x00, 0x14, // sport (20)
		0x00, 0x50, // dport (80)
		0x00, 0x00, 0x00, 0x01, // seqno (1)
		0x00, 0x00, 0x00, 0x00, // ackno (0)
		0x50,       // doff + reserved (5 << 4)
		0xFF,       // 无效标志位
		0x00, 0x10, // win (16)
		0x12, 0x34, // cksum (0x1234)
		0x00, 0x00, // uptr
	}
	buffer := utils.NewBuffer(string(data))
	parser := utils.NewNetParser[uint8](*buffer)

	tcpHeader := TCPHeader[uint8]{}
	result := tcpHeader.Parse(*parser)

	if result != utils.NoError {
		t.Errorf("Parse() result = %v, want %v", result, utils.NoError)
	}
	if !tcpHeader.urg || !tcpHeader.ack || !tcpHeader.psh || !tcpHeader.rst || !tcpHeader.syn || !tcpHeader.fin {
		t.Errorf("tcpHeader flags incorrect")
	}
}

func TestTCPHeader_Serialize(t *testing.T) {
	tests := []struct {
		name     string
		header   TCPHeader[uint32]
		expected string
	}{
		{
			name: "Basic Test Case",
			header: TCPHeader[uint32]{
				sport: 12,
				dport: 0x12,
				seqno: 0x9ABCDEF0,
				ackno: 0x12345678,
				doff:  5,
				urg:   true,
				ack:   false,
				psh:   true,
				rst:   false,
				syn:   true,
				fin:   false,
				win:   0xFFFF,
				cksum: 0xABCD,
				uptr:  0x1234,
			},
			expected: string([]byte{
				0x12,                   // sport
				0x12,                   // dport
				0x9A, 0xBC, 0xDE, 0xF0, // seqno
				0x12, 0x34, 0x56, 0x78, // ackno
				0x50,       // doff << 4
				0b00101001, // flags (urg=1, psh=1, syn=1)
				0xFF, 0xFF, // win
				0xAB, 0xCD, // cksum
				0x12, 0x34, // uptr
				0x00, 0x00, 0x00, 0x00, // Padding to 4 * doff bytes
			}),
		},
		{
			name: "Header too short",
			header: TCPHeader[uint32]{
				doff: 4,
			},
			expected: "TCP header too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.header.Serialize()
			if got != tt.expected {
				t.Errorf("Serialize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
