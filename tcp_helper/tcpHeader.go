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
//!  |           CheCksum            |         Urgent Pointer        |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                    Options                    |    Padding    |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//!  |                             data                              |
//!  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type TCPHeader[T uint32 | uint16 | uint8] struct {
	Sport T
	Dport T
	Seqno uint32
	Ackno uint32
	Doff  T
	Urg   bool
	Ack   bool
	Psh   bool
	Rst   bool
	Syn   bool
	Fin   bool
	Win   T
	Cksum T
	Uptr  T
}

var _ TCPHeaderInterface[uint32] = (*TCPHeader[uint32])(nil)

func NewTcpHeader[T uint32 | uint16 | uint8]() *TCPHeader[T] {
	return &TCPHeader[T]{
		Sport: 0,
		Dport: 0,
		Seqno: 0,
		Ackno: 0,
		Doff:  TCPHeaderLENGTH / 4,
		Urg:   false,
		Ack:   false,
		Psh:   false,
		Rst:   false,
		Syn:   false,
		Fin:   false,
		Win:   0,
		Cksum: 0,
		Uptr:  0,
	}
}

func (t *TCPHeader[T]) Parse(p utils.NetParser[T]) utils.ParseResult {
	t.Sport = p.ParseInt(int(unsafe.Sizeof(uint16(0)))) // 2
	t.Dport = p.ParseInt(int(unsafe.Sizeof(uint16(0))))

	Seqno := p.ParseInt(int(unsafe.Sizeof(uint32(0)))) // 4
	Ackno := p.ParseInt(int(unsafe.Sizeof(uint32(0))))

	wrappingInt32 := wrapping.WrappingInt32{}
	wrappingInt32.SetRawValue(uint32(Seqno))
	t.Seqno = wrappingInt32.RawValue()
	wrappingInt32.SetRawValue(uint32(Ackno))

	t.Ackno = wrappingInt32.RawValue()
	t.Doff = p.ParseInt(int(unsafe.Sizeof(uint8(0)))) >> 4
	f := p.ParseInt(int(unsafe.Sizeof(uint8(0))))

	t.Urg = f&0b00100000 != 0
	t.Ack = f&0b00010000 != 0
	t.Psh = f&0b00001000 != 0
	t.Rst = f&0b00000100 != 0
	t.Syn = f&0b00000010 != 0
	t.Fin = f&0b00000001 != 0

	t.Win = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	t.Cksum = p.ParseInt(int(unsafe.Sizeof(uint16(0))))
	t.Uptr = p.ParseInt(int(unsafe.Sizeof(uint16(0))))

	if t.Doff < 5 {
		return utils.PacketTooShort
	}

	p.RemovePrefix(int(t.Doff*4 - TCPHeaderLENGTH))
	if p.Error() {
		return p.GetError()
	}
	return utils.NoError
}
func (t *TCPHeader[T]) Serialize() string {
	if t.Doff < 5 {
		return "TCP header too short"
	}
	ret := make([]byte, 0, 4*t.Doff)
	unparser := utils.NetUnparser[T]{}
	unparser.UnparseInt(&ret, t.Sport, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.Dport, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, T(t.Seqno), int(unsafe.Sizeof(uint32(0))))
	unparser.UnparseInt(&ret, T(t.Ackno), int(unsafe.Sizeof(uint32(0))))
	unparser.UnparseInt(&ret, t.Doff<<4, int(unsafe.Sizeof(uint8(0))))
	flags := (0b00100000 * BoolToUint8(t.Urg)) |
		(0b00010000 * BoolToUint8(t.Ack)) |
		(0b00001000 * BoolToUint8(t.Psh)) |
		(0b00000100 * BoolToUint8(t.Rst)) |
		(0b00000010 * BoolToUint8(t.Syn)) |
		(0b00000001 * BoolToUint8(t.Fin))
	unparser.UnparseInt(&ret, T(flags), int(unsafe.Sizeof(uint8(0))))
	unparser.UnparseInt(&ret, t.Win, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.Cksum, int(unsafe.Sizeof(uint16(0))))
	unparser.UnparseInt(&ret, t.Uptr, int(unsafe.Sizeof(uint16(0))))
	//ret = append(ret, make([]byte, 4*t.Doff-T(len(ret)))...)
	if cap(ret) >= int(t.Doff<<4) {
		ret = ret[:4*t.Doff]
	} else {
		newRet := make([]byte, 4*t.Doff)
		copy(newRet, ret)
		ret = newRet
	}
	log.Printf("ret size = %d", len(ret))
	log.Printf("0x = %+v \n s = %s \n the header struct = %v", ret, fmt.Sprintf("%s", ret), t)
	return string(ret)
}
func (t *TCPHeader[T]) ToString() string {
	return fmt.Sprintf("TCP source port: %d\n"+
		"TCP dest port: %d\n"+
		"TCP Seqno: %d\n"+
		"TCP Ackno: %d\n"+
		"TCP Doff: %d\n"+
		"Flags: Urg: %t Ack: %t Psh: %t Rst: %t Syn: %t Fin: %t\n"+
		"TCP Winsize: %d\n"+
		"TCP Cksum: %d\n"+
		"TCP Uptr: %d\n",
		t.Sport, t.Dport, t.Seqno, t.Ackno, t.Doff,
		t.Urg, t.Ack, t.Psh, t.Rst, t.Syn, t.Fin,
		t.Win, t.Cksum, t.Uptr)
}

func (t *TCPHeader[T]) Summary() string {
	flags := ""
	if t.Syn {
		flags += "S"
	}
	if t.Ack {
		flags += "A"
	}
	if t.Rst {
		flags += "R"
	}
	if t.Fin {
		flags += "F"
	}
	return fmt.Sprintf("Header(flags=%s, Seqno=%d, Ack=%d, Win=%d)",
		flags, t.Seqno, t.Ackno, t.Win)
}

func BoolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
