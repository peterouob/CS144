package utils

type ParseResult int

const (
	NoError ParseResult = iota
	BadCheckSum
	PacketTooShort
	WrongIPVersion
	HeaderTooShort
	TruncatedPacket
	Unsupported
)

func AsString(r ParseResult) string {
	names := []string{
		"NoError",
		"BadChecksum",
		"PacketTooShort",
		"WrongIPVersion",
		"HeaderTooShort",
		"TruncatedPacket",
		"Unsupported",
	}

	if int(r) < len(names) {
		return names[r]
	}
	return "Unknown"
}

type NetParserInterface[T uint32 | uint16 | uint8] interface {
	Buffer() Buffer
	GetError() ParseResult
	SetError(ParseResult)
	Error() bool
	CheckSize(int)
	ParseInt(int) T
	RemovePrefix(int)
}

type NetParser[T uint32 | uint16 | uint8] struct {
	buffer Buffer
	error  ParseResult
}

func NewNetParser[T uint32 | uint16 | uint8](buffer Buffer) *NetParser[T] {
	return &NetParser[T]{
		buffer: buffer,
		error:  NoError,
	}
}

var _ NetParserInterface[uint32] = (*NetParser[uint32])(nil)

func (n *NetParser[T]) Buffer() Buffer {
	return n.buffer
}

func (n *NetParser[T]) GetError() ParseResult {
	return n.error
}

func (n *NetParser[T]) SetError(res ParseResult) {
	n.error = res
}

func (n *NetParser[T]) Error() bool {
	return n.GetError() != NoError
}

func (n *NetParser[T]) CheckSize(size int) {
	if n.buffer.Size() < size {
		n.SetError(PacketTooShort)
	}
}

func (n *NetParser[T]) ParseInt(len int) T {
	n.CheckSize(len)
	if n.Error() {
		return 0
	}
	var ret T
	for i := 0; i < len; i++ {
		ret = ret << 8
		pos, at := n.buffer.At(i)
		if !at {
			return 0
		}
		ret += T(pos)
	}
	n.buffer.RemovePrefix(len)
	return ret
}

func (n *NetParser[T]) RemovePrefix(len int) {
	n.CheckSize(len)
	if n.Error() {
		return
	}
	n.buffer.RemovePrefix(len)
}

type NetUnparserInterface[T uint32 | uint16 | uint8] interface {
	UnparseInt(string, T, int) string
}

type NetUnparser[T uint32 | uint16 | uint8] struct{}

var _ NetUnparserInterface[uint32] = (*NetUnparser[uint32])(nil)

func (nu *NetUnparser[T]) UnparseInt(s string, val T, n int) string {
	for i := 0; i < n; i++ {
		theByte := byte((val >> uint((n-i-1)*8)) & 0xff)
		s += string(theByte)
	}
	return s
}
