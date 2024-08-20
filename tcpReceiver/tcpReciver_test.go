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
	stream := stream.NewStream(stream.Deque{}, 1024, 0, 0, false, false)
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

	segment := tcp_helper.TCPSegment{
		Header:  *header,
		Payload: stream.Deque{}, // No payload
	}

	// Call SegmentReceived to simulate receiving the SYN
	receiver.SegmentReceived(segment)

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
	header.Seqno = 101 // Sequence number following the SYN
	payload := stream.Deque{}
	payload.PushBack('h')
	payload.PushBack('i')
	segment.Payload = payload

	// Call SegmentReceived to simulate receiving a data segment
	receiver.SegmentReceived(segment)

	// Check if the payload is correctly reassembled
	outputStream := receiver.SegmentOut()
	if outputStream.BytesWritten() != 2 {
		t.Errorf("Expected 2 bytes written, got %d", outputStream.BytesWritten())
	}

	// Check if the received data matches the payload
	data := outputStream.Read(2)
	if data != "hi" {
		t.Errorf("Expected payload 'hi', got '%s'", data)
	}
}
