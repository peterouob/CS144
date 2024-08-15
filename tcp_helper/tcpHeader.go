package tcp_helper

import (
	"lab/utils"
	"lab/wrapping"
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
	sport uint16
	dport uint16
	seqno wrapping.WrappingInt32
	ackno wrapping.WrappingInt32
	doff  uint8
	urg   bool
	ack   bool
	psh   bool
	rst   bool
	syn   bool
	fin   bool
	win   uint16
	cksum uint16
	uptr  uint16
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

func (t *TCPHeader[T]) Parse(p utils.NetParser[T]) utils.ParseResult {}
func (t *TCPHeader[T]) Serialize() string                            {}
func (t *TCPHeader[T]) ToString() string                             {}
func (t *TCPHeader[T]) Summary() string                              {}
