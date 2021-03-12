package main

import (
	"strconv"
	"rudp"
)
func main() {
	conn:= rudp.DialUDP("127.0.0.1:9981")
	writeFunc(conn)
	// go readFunc(conn)
	select{}

}

func writeFunc(conn *rudp.RuConn){
	for i:=1;i<=100;i++{
		conn.Write([]byte("seq:"+strconv.Itoa(i)))
	}
}

// func readFunc(conn *net.UDPConn){
// 	for{
// 		data := make([]byte, 1024)
// 		n, _ := conn.Read(data)
// 		fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
// 	}

// }