package tcpReceiver

import (
	"lab/stream"
	"lab/streamReassembler"
	"lab/tcp_helper"
	"lab/wrapping"
	"testing"
)

func TestTcpReceiver_SegmentReceived(t *testing.T) {
	// Initialize WrappingInt32 with initial sequence number (ISN)
	isn := wrapping.WrappingInt32{}
	isn.SetRawValue(100) // Starting ISN

	// Create a mock Stream and StreamReassembler
	q := stream.NewDeque(2048)
	stream := stream.NewStream(*q, 1024, 0, 0)
	reassembler := streamReassembler.NewStreamReassembler(1024, stream)

	// Initialize the TcpReceiver
	receiver := &TcpReceiver{
		isn:         isn,
		setSynFlag:  false, // Start with SYN not set
		reassembler: *reassembler,
		capacity:    1024,
	}

	// Prepare a TCP Segment with SYN flag and initial sequence number
	header := tcp_helper.NewTcpHeader[uint32]()
	header.Syn = true
	header.Seqno = 100 // SYN with ISN 100

	segment := tcp_helper.NewTCPSegment(*header)

	// Call SegmentReceived to simulate receiving the SYN
	receiver.SegmentReceived(*segment)

	// Check if the SYN flag is correctly set
	if !receiver.setSynFlag {
		t.Errorf("SYN flag was not set in TcpReceiver after receiving SYN segment")
	}

	// Check if the acknowledgment number is correctly computed
	expectedAckno := isn.SetRawValue(101)
	ackno := receiver.Ackno()
	if ackno != *expectedAckno {
		t.Errorf("Incorrect Ackno: got %v, want %v", ackno.RawValue(), expectedAckno.RawValue())
	}

	// Add a new TCP segment with payload
	header.Syn = false // Not a SYN, regular data
	header.Seqno = 100 // Sequence number following the SYN

	payload := stream
	payload.Write("hi")
	segment.SetPaylaod(*payload)
	reassembler.SetunassembleStrs("hi", "h")
	// Call SegmentReceived to simulate receiving a data segment
	receiver.SegmentReceived(*segment)
	r := receiver.reassembler
	put := r.StreamOut()
	p := put.Read(2)

	if p != "hi" {
		t.Errorf("need =%s,got =%s", "hi", p)
	}

	// Check if the payload is correctly reassembled
	rc := receiver.SegmentOut()
	outputStream := rc.Read(2)

	if outputStream != "hi" {
		t.Errorf("need =%s,got =%s", "hi", outputStream)
	}
}

func TestTcpReceiverWithPayload(t *testing.T) {
	isn := wrapping.WrappingInt32{}
	isn.SetRawValue(100) // Starting ISN

	// Create a mock Stream and StreamReassembler
	q := stream.NewDeque(2048)
	stream := stream.NewStream(*q, 1024, 0, 0)
	reassembler := streamReassembler.NewStreamReassembler(1024, stream)

	// Initialize the TcpReceiver
	receiver := &TcpReceiver{
		isn:         isn,
		setSynFlag:  false, // Start with SYN not set
		reassembler: *reassembler,
		capacity:    1024,
	}

	// Prepare a TCP Segment with SYN flag and initial sequence number
	header := tcp_helper.NewTcpHeader[uint32]()
	header.Syn = false // Not a SYN, regular data
	header.Seqno = 100 // Sequence number following the SYN

	segment := tcp_helper.NewTCPSegment(*header)

	payload := stream
	payload.Write("hi")
	segment.SetPaylaod(*payload)
	reassembler.SetunassembleStrs("hi")
	// Call SegmentReceived to simulate receiving a data segment
	receiver.SegmentReceived(*segment)
	r := receiver.reassembler
	put := r.StreamOut()
	p := put.Read(2)

	if p != "hi" {
		t.Errorf("need =%s,got =%s", "hi", p)
	}

	// Check if the payload is correctly reassembled
	rc := receiver.SegmentOut()
	outputStream := rc.Read(2)

	if outputStream != "hi" {
		t.Errorf("need =%s,got =%s", "hi", outputStream)
	}
}
