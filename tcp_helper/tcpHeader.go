package tcp_helper

import (
	"lab/utils"
	"lab/wrapping"
	"unsafe"
)

const TCPHeaderLENGTH = 20

type TCPHeaderInterface[T uint32 | uint16 | uint8] interface {
	Parse(utils.NetParser[T]) utils.ParseResult
	Serialize() string
	ToString() string
	Summary() string
}

//!   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |          Source Port          |       Destination Port        |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                        Sequence Number                        |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                    Acknowledgment Number                      |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |  Data |           |U|A|P|R|S|F|                               |
//!  | Offset| Reserved  |R|C|S|S|Y|I|            Window             |
//!  |       |           |G|K|H|T|N|N|                               |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |           Checksum            |         Urgent Pointer        |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                    Options                    |    Padding    |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                             data                              |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type TCPHeader[T uint32 | uint16 | uint8] struct {
	sport T
	dport T
	seqno wrapping.WrappingInt32
	ackno wrapping.WrappingInt32
	doff  T
	urg   bool
	ack   bool
	psh   bool
	rst   bool
	syn   bool
	fin   bool
	win   T
	cksum T
	uptr  T
}

var _ TCPHeaderInterface[uint32] = (*TCPHeader[uint32])(nil)

func NewTcpHeader[T uint32 | uint16 | uint8]() *TCPHeader[T] {
	return &TCPHeader[T]{
		sport: 0,
		dport: 0,
		seqno: *(wrapping.NewWrrappingInt32()),
		ackno: *(wrapping.NewWrrappingInt32()),
		doff:  TCPHeaderLENGTH / 4,
		urg:   false,
		ack:   false,
		psh:   false,
		rst:   false,
		syn:   false,
		fin:   false,
		win:   0,
		cksum: 0,
		uptr:  0,
	}
}

func (t *TCPHeader[T]) Parse(p utils.NetParser[T]) utils.ParseResult {
	t.sport = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	t.dport = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	seqno := p.ParseInt(int(unsafe.Sizeof(uint32(0))))
	ackno := p.ParseInt(int(unsafe.Sizeof(uint32(0))))
	wrappingInt32 := wrapping.WrappingInt32{}
	t.seqno = *(wrappingInt32.SetRawValue(uint32(seqno)))
	t.ackno = *(wrappingInt32.SetRawValue(uint32(ackno)))
	t.doff = p.ParseInt(int(unsafe.Sizeof(uint8(0)))) >> 4
	f := p.ParseInt(int(unsafe.Sizeof(uint8(0))))

	t.urg = f&0b00100000 != 0
	t.ack = f&0b00010000 != 0
	t.psh = f&0b00001000 != 0
	t.rst = f&0b00000100 != 0
	t.syn = f&0b00000010 != 0
	t.fin = f&0b00000001 != 0

	t.win = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	t.cksum = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	t.uptr = p.ParseInt(int(unsafe.Sizeof(uint16(0))))

	if t.doff < 5 {
		return utils.PacketTooShort
	}

	p.RemovePrefix(int(t.doff*4 - TCPHeaderLENGTH))
	if p.Error() {
		return p.GetError()
	}
	return utils.NoError
}
func (t *TCPHeader[T]) Serialize() string { return "" }
func (t *TCPHeader[T]) ToString() string  { return "" }
func (t *TCPHeader[T]) Summary() string   { return "" }
