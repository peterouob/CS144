package tcp_helper

import (
	"lab/utils"
	"testing"
	"unsafe"
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
	header := NewTcpHeader[uint32]()
	result := header.Serialize()
	expectedLength := int(unsafe.Sizeof(uint16(0)))*2 + int(unsafe.Sizeof(uint32(0)))*2 + 1 + 1 + int(unsafe.Sizeof(uint16(0)))*3
	if len(result) != expectedLength {
		t.Errorf("Serialize failed, expected length %d, got %d", expectedLength, len(result))
	}

	header.sport = 12345
	header.dport = 80
	header.seqno = 12345678
	header.ackno = 87654321
	header.doff = 5
	header.urg = true
	header.ack = true
	header.psh = false
	header.rst = false
	header.syn = true
	header.fin = false
	header.win = 65535
	header.cksum = 1234
	header.uptr = 5678

	result = header.Serialize()
	expectedLength = int(header.doff) * 4
	if len(result) != expectedLength {
		t.Errorf("Serialize failed, expected length %d, got %d", expectedLength, len(result))
	}

	// Check some field values in the serialized output
	if result[0] != byte(header.sport>>8) || result[1] != byte(header.sport&0xff) {
		t.Errorf("Serialize failed, source port does not match")
	}
	if result[2] != byte(header.dport>>8) || result[3] != byte(header.dport&0xff) {
		t.Errorf("Serialize failed, destination port does not match")
	}

	// Check flags
	flagsIndex := 13 // After doff, flags start at byte index 13
	flags := uint8(0) |
		(0b00100000 * BoolToUint8(header.urg)) |
		(0b00010000 * BoolToUint8(header.ack)) |
		(0b00000010 * BoolToUint8(header.syn))
	if result[flagsIndex] != flags {
		t.Errorf("Serialize failed, flags do not match. Expected %08b, got %08b", flags, result[flagsIndex])
	}
}
