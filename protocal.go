package rudp

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Protocol struct{
	
}

const (
    commonMsg uint8 = 1 + iota
    ackMsg
)

type PacketHeader struct{
	packetType uint8
	Sequence  uint16

}

type WrittenPacket struct {
	// seq     uint16  
	data []byte 
	sendTime time.Time  // last time the packet was written
	resent  uint8      // total number of times this packet was resent
}

func NewAck(ack uint16) []byte{
	header:=&PacketHeader{
		packetType: ackMsg,
		Sequence: ack,
	}
	buf := new(bytes.Buffer)
	binary.Write(buf,binary.BigEndian,header)
	return buf.Bytes()
}

func MarshalPacketHeader(seq uint16,data []byte) (buf []byte) {
	buf=make([]byte,10)
	binary.BigEndian.PutUint16(buf,seq)
	AppendHead(buf,commonMsg)
	buf=append(buf,data...)
	return
}

func UnmarshalPacketHeader(buf []byte)(msgType uint8,seq uint16, leftover []byte){
	msgType = buf[0]
	seq =binary.BigEndian.Uint16(buf[1:3])
	leftover=buf[3:]
	return
}