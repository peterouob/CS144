package tcp_helper

import (
	"fmt"
	"lab/utils"
	"lab/wrapping"
	"log"
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
	seqno uint32
	ackno uint32
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
		seqno: 0,
		ackno: 0,
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
	wrappingInt32.SetRawValue(uint32(seqno))
	t.seqno = wrappingInt32.RawValue()
	wrappingInt32.SetRawValue(uint32(ackno))

	t.ackno = wrappingInt32.RawValue()
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
func (t *TCPHeader[T]) Serialize() string {
	if t.doff < 5 {
		return "TCP header too short"
	}
	ret := make([]byte, 0, 4*t.doff)
	unparser := utils.NetUnparser[T]{}
	unparser.UnparseInt(&ret, t.sport, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.dport, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, T(t.seqno), int(unsafe.Sizeof(uint32(0))))
	unparser.UnparseInt(&ret, T(t.ackno), int(unsafe.Sizeof(uint32(0))))
	unparser.UnparseInt(&ret, t.doff<<4, int(unsafe.Sizeof(uint8(0))))
	flags := uint8((0b00100000 * boolToUint8(t.urg)) |
		(0b00010000 * boolToUint8(t.ack)) |
		(0b00001000 * boolToUint8(t.psh)) |
		(0b00000100 * boolToUint8(t.rst)) |
		(0b00000010 * boolToUint8(t.syn)) |
		(0b00000001 * boolToUint8(t.fin)))
	unparser.UnparseInt(&ret, T(flags), int(unsafe.Sizeof(uint8(0))))
	unparser.UnparseInt(&ret, t.win, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.cksum, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.uptr, int(unsafe.Sizeof(uint16(0))))
	if T(len(ret)) > 4*t.doff {
		ret = append(ret, make([]byte, 4*t.doff-T(len(ret)))...)
	}
	log.Printf("0x = %+v \n s = %s \n the header struct = %v", ret, fmt.Sprintf("%s", ret), t)
	return string(ret)
}
func (t *TCPHeader[T]) ToString() string { return "" }
func (t *TCPHeader[T]) Summary() string  { return "" }
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
