package main

import (
	"fmt"
	"net"
	"strconv"
	"rudp"
)
func main() {
	sip := net.ParseIP("127.0.0.1.155")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	rudp.LogError(err)
	defer conn.Close()
	go writeFunc(conn)
	go readFunc(conn)
	select{}

}

func writeFunc(conn *net.UDPConn){
	for i:=1;i<=100;i++{
		conn.Write([]byte("seq:"+strconv.Itoa(i)))
	}
}

func readFunc(conn *net.UDPConn){
	for{
		data := make([]byte, 1024)
		n, _ := conn.Read(data)
		fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
	}

}