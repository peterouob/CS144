package tcpReceiver

import (
	"lab/stream"
	"lab/streamReassembler"
	"lab/tcp_helper"
	"lab/wrapping"
)

type ReceiverInterface interface {
	Ackno() wrapping.WrappingInt32
	WindowSize() int
	UnassembledBytes() int
	SegmentReceived(tcp_helper.TCPSegment)
	SegmentOut() stream.Stream
}

type TcpReceiver struct {
	isn         wrapping.WrappingInt32
	setSynFlag  bool
	reassembler streamReassembler.StreamReassembler
	capacity    int
}

var _ ReceiverInterface = (*TcpReceiver)(nil)

func (rcv *TcpReceiver) Ackno() wrapping.WrappingInt32 {
	if !rcv.setSynFlag {
		return wrapping.WrappingInt32{}
	}
	ressmebler := rcv.reassembler.StreamOut()
	absAckNo := ressmebler.BytesWritten() + 1
	if ressmebler.InputEnded() {
		absAckNo += 1
	}
	val := rcv.isn.RawValue()
	return *rcv.isn.SetRawValue(val + uint32(absAckNo))
}

func (rcv *TcpReceiver) WindowSize() int {
	out := rcv.reassembler.StreamOut()
	return rcv.capacity - out.BufferSize()
}

func (rcv *TcpReceiver) UnassembledBytes() int {
	return rcv.reassembler.UnassembledBytes()
}

func (rcv *TcpReceiver) SegmentReceived(seg tcp_helper.TCPSegment) {
	header := seg.GetHeader()
	if !rcv.setSynFlag {
		if !header.Syn {
			return
		}
		rcv.isn.SetRawValue(header.Seqno)
		rcv.setSynFlag = true
	}
	rstream := rcv.reassembler.StreamOut()
	absAckno := rstream.BytesWritten() + 1

	wrap := wrapping.WrappingInt32{}
	currAbsSeqno := wrap.UnWrap(*wrap.SetRawValue(header.Seqno), rcv.isn, uint64(absAckno))

	streamIndex := currAbsSeqno - 1 + checkSyn(header.Syn)
	payload := seg.GetPayload()
	rcv.reassembler.PushSubString(payload.Copy(), int(streamIndex), header.Fin)
}

func (rcv *TcpReceiver) SegmentOut() stream.Stream {
	return rcv.reassembler.StreamOut()
}

func checkSyn(syn bool) uint64 {
	if syn {
		return 1
	} else {
		return 0
	}
}
