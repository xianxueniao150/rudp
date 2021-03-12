package main

import (
	"fmt"
	"rudp"
	"strconv"
	"time"
)

var (
	t = time.Now()
	nums = make([]uint32, 0, 100)
)
	

func main() {
	listener := rudp.ListenUDP("0.0.0.0:9981")
	// listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 9981})
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	listener.RegisterHandleFunc(handleFunc)
	time.Sleep(time.Second * 3)
	listener.Run()
}

func handleFunc(conn *rudp.RuConn, data []byte) {
	fmt.Printf("<%s> %s\n", conn.RemoteAddr,data)
	a, _ := strconv.Atoi(string(data[12:]))
	nums = append(nums, uint32(a))
	if t.Add(time.Second * 3).Before(time.Now()) {
		// fmt.Println("begin")
		go testResult()
	}
	t = time.Now()
}

func testResult() {
	t = time.Now()
	for {
		time.Sleep(time.Millisecond)
		// fmt.Println("wait",time.Now())
		if t.Add(time.Second * 3).Before(time.Now()) {
			fmt.Println(len(nums))
			fmt.Println(nums)
			nums = make([]uint32, 0, 100)
			break
		}
	}
}
