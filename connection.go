package rudp

import (
	"fmt"
	"net"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/orcaman/concurrent-map"
)

type handleFunc func(conn *RuConn, buf []byte)

type Listener struct {
	*net.UDPConn
	HandleFunc  handleFunc
	connections map[string]*RuConn
}

type RuConn struct {
	Listener   *Listener
	Conn       *net.UDPConn
	RemoteAddr *net.UDPAddr
	lastAck    uint16 // representing the last consecutive sequence number of a packet we have sent that was acknowledged by our peer. For example, if we have sent packets whose sequence numbers are in the range [0, 256], and we have received acknowledgements for packets (0, 1, 2, 3, 4, 8, 9, 10, 11, 12), then (oui) would be 4
	lastSend   uint16
	rq         mapset.Set         // read queue
	wq         cmap.ConcurrentMap /* map[uint16]*WrittenPacket */
}

func newConn(udpaddr *net.UDPAddr, conn *net.UDPConn, listener *Listener) *RuConn {
	return &RuConn{
		RemoteAddr: udpaddr,
		rq:         mapset.NewSet(),
		wq:         cmap.New(),
		Conn:       conn,
		Listener:   listener,
	}
}

func (listener *Listener) getConn(udpaddr *net.UDPAddr) (conn *RuConn) {
	conn, ok := listener.connections[udpaddr.String()]
	if !ok {
		fmt.Println("new Conn")
		conn = newConn(udpaddr, nil, listener)
		listener.connections[udpaddr.String()] = conn
	}
	return
}

func (listener *Listener) RegisterHandleFunc(handler handleFunc) {
	listener.HandleFunc = handler
}

func (listener *Listener) Run() {
	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(buf)
		LogError(err)
		conn := listener.getConn(remoteAddr)
		_, seq, left := UnmarshalPacketHeader(buf[:n])
		fmt.Printf("seq: <%d> \n", seq)
		if seq > conn.lastAck && !conn.rq.Contains(seq) {
			fmt.Printf("conn.lastAck: <%d> \n", conn.lastAck)
			if seq == conn.lastAck+1 {
				conn.lastAck = seq
				fmt.Printf("conn.lastAck after: <%d> \n", conn.lastAck)
				index := seq + 1
				for {
					if conn.rq.Contains(index) {
						conn.rq.Remove(index)
						conn.lastAck = index
						index++
					} else {
						break
					}
				}
				conn.ReplyAck(conn.lastAck)
			} else {
				conn.rq.Add(seq)
			}
			listener.HandleFunc(conn, left)
		}
	}
}

func (conn *RuConn) ReplyAck(ack uint16) {
	data := NewAck(ack)
	conn.Listener.WriteToUDP(data, conn.RemoteAddr)
}

func (conn *RuConn) Write(data []byte) {
	conn.lastSend++
	buf := MarshalPacketHeader(conn.lastSend, data)
	conn.wq.Set(string(conn.lastSend), &WrittenPacket{
		sendTime: time.Now(),
		data:     buf,
	})
	conn.Conn.Write(buf)
	// conn.Conn.WriteToUDP(buf, conn.RemoteAddr) 这里有问题，不知道为什么
}

func ListenUDP(laddr string) *Listener {
	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	LogError(err)
	conn, err := net.ListenUDP("udp", udpaddr)
	LogError(err)
	return &Listener{UDPConn: conn,
		connections: make(map[string]*RuConn, 100)}
}

func (conn *RuConn) ReceiveMsg() {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Conn.Read(buf)
		fmt.Printf("receive from server data:%v\n", buf[:10])
		LogError(err)
		msgType, seq, left := UnmarshalPacketHeader(buf[:n])
		switch msgType {
		case ackMsg:
			conn.DealAck(seq)
		case commonMsg:
			conn.DealCommMsg(seq, left)
		}

	}
}

func (conn *RuConn) DealAck(seq uint16) {
	for ; conn.lastAck <= seq; conn.lastAck++ {
		conn.wq.Remove(string(conn.lastAck))
	}
}

func (conn *RuConn) DealCommMsg(seq uint16, data []byte) {
	// fmt.Printf("seq: <%d> \n", seq)
	// 	if seq > conn.lastAck && !conn.rq.Contains(seq) {
	// 		if seq == conn.lastAck+1 {
	// 			index := seq + 1
	// 			for {
	// 				if conn.rq.Contains(index) {
	// 					conn.rq.Remove(index)
	// 					conn.lastAck = index
	// 					index++
	// 				} else {
	// 					break
	// 				}
	// 			}
	// 			conn.ReplyAck(conn.lastAck)
	// 		} else {
	// 			conn.rq.Add(seq)
	// 		}
	// 		listener.HandleFunc(conn, left)
	// 	}
}

func DialUDP(laddr string) *RuConn {
	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	LogError(err)
	conn, err := net.DialUDP("udp", nil, udpaddr)
	LogError(err)
	ruConn := newConn(udpaddr, conn, nil)
	go ruConn.ReceiveMsg()
	go ruConn.StartResendCheck()
	// ruConn.
	return ruConn
}

func (conn *RuConn) StartResendCheck() {
	ticker := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for tuple := range conn.wq.IterBuffered() {
				packet := tuple.Val.(*WrittenPacket)
				if packet.resent<10 {
					conn.Conn.Write(packet.data)
					packet.resent++
					packet.sendTime = time.Now()
				}
			}
		}
	}
}
